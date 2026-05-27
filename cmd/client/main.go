package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nambuitechx/nam-tcp/internal/proxy"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "connect":
		runConnect(os.Args[2:])
	case "send":
		runSend(os.Args[2:])
	case "forward":
		runForward(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `nam-tcp client

Usage:
  nam-tcp-client forward -local <host:port> -proxy <host:port> -token <pat>
  nam-tcp-client connect -proxy <host:port> -token <pat>
  nam-tcp-client send    -proxy <host:port> -token <pat> [-data <text>]

Environment:
  NAM_TCP_PROXY  default proxy address
  NAM_TCP_TOKEN  default pat token
  NAM_TCP_LOCAL  default local bind address for forward

Examples:
  nam-tcp-client forward -local 127.0.0.1:15432 -proxy ec2:8888 -token nam_tcp_...
  psql -h 127.0.0.1 -p 15432 -U admin mydb
  nam-tcp-client connect -proxy localhost:8888 -token nam_tcp_...
  nam-tcp-client send -proxy localhost:8888 -token nam_tcp_... -data "ping"
`)
}

func runForward(args []string) {
	fs := flag.NewFlagSet("forward", flag.ExitOnError)
	localAddr := fs.String("local", envOr("NAM_TCP_LOCAL", "127.0.0.1:15432"), "local address to listen on")
	proxyAddr := fs.String("proxy", envOr("NAM_TCP_PROXY", "localhost:8888"), "proxy listen address")
	token := fs.String("token", envOr("NAM_TCP_TOKEN", ""), "pat token")
	tokenFile := fs.String("token-file", "", "file containing pat token")
	_ = fs.Parse(args)

	tok, err := resolveToken(*token, *tokenFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fl, err := proxy.NewForwardListener(*localAddr, *proxyAddr, tok)
	if err != nil {
		fmt.Fprintf(os.Stderr, "forward failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "forwarding %s -> %s (Ctrl+C to stop)\n", fl.Addr(), *proxyAddr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- fl.Run()
	}()

	select {
	case <-ctx.Done():
		_ = fl.Close()
		fl.Wait()
		fmt.Fprintln(os.Stderr, "forward stopped")
	case err := <-errCh:
		if err != nil {
			fmt.Fprintf(os.Stderr, "forward failed: %v\n", err)
			os.Exit(1)
		}
	}
}

func runConnect(args []string) {
	fs := flag.NewFlagSet("connect", flag.ExitOnError)
	proxyAddr := fs.String("proxy", envOr("NAM_TCP_PROXY", "localhost:8888"), "proxy listen address")
	token := fs.String("token", envOr("NAM_TCP_TOKEN", ""), "pat token")
	tokenFile := fs.String("token-file", "", "file containing pat token")
	_ = fs.Parse(args)

	tok, err := resolveToken(*token, *tokenFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conn, err := proxy.Dial(*proxyAddr, tok)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	done := make(chan struct{}, 2)
	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		done <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		done <- struct{}{}
	}()

	<-done
}

func runSend(args []string) {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	proxyAddr := fs.String("proxy", envOr("NAM_TCP_PROXY", "localhost:8888"), "proxy listen address")
	token := fs.String("token", envOr("NAM_TCP_TOKEN", ""), "pat token")
	tokenFile := fs.String("token-file", "", "file containing pat token")
	data := fs.String("data", "", "payload to send (default: read stdin)")
	_ = fs.Parse(args)

	tok, err := resolveToken(*token, *tokenFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conn, err := proxy.Dial(*proxyAddr, tok)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	var payload []byte
	if *data != "" {
		payload = []byte(*data)
	} else {
		payload, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read stdin: %v\n", err)
			os.Exit(1)
		}
	}

	if len(payload) > 0 {
		if _, err := conn.Write(payload); err != nil {
			fmt.Fprintf(os.Stderr, "write: %v\n", err)
			os.Exit(1)
		}
	}

	buf := make([]byte, 32*1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "read response: %v\n", err)
		os.Exit(1)
	}
	if n > 0 {
		os.Stdout.Write(buf[:n])
	}
}

func resolveToken(token, tokenFile string) (string, error) {
	if tokenFile != "" {
		b, err := os.ReadFile(tokenFile)
		if err != nil {
			return "", fmt.Errorf("read token file: %w", err)
		}
		token = string(b)
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("token is required (-token, -token-file, or NAM_TCP_TOKEN)")
	}

	return token, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
