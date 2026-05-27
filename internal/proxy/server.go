package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"sync/atomic"
)

type TargetResolver interface {
	ResolveTarget(plaintextToken string) (string, error)
}

type ProxyServer struct {
	listener     net.Listener
	resolver     TargetResolver
	clients      map[net.Conn]struct{}
	mu           sync.Mutex
	done         chan struct{}
	closing      atomic.Bool
	shutdownOnce sync.Once
}

func NewProxyServer(address string, resolver TargetResolver) (*ProxyServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %w", address, err)
	}

	return &ProxyServer{
		listener: listener,
		resolver: resolver,
		clients:  make(map[net.Conn]struct{}),
		done:     make(chan struct{}),
	}, nil
}

func (ps *ProxyServer) Run() {
	defer close(ps.done)

	for {
		client, err := ps.listener.Accept()
		if err != nil {
			if ps.closing.Load() || errors.Is(err, net.ErrClosed) {
				return
			}
			log.Println("accept error:", err)
			continue
		}

		if !ps.trackClient(client) {
			client.Close()
			continue
		}

		go ps.handle(client)
	}
}

func (ps *ProxyServer) trackClient(c net.Conn) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closing.Load() {
		return false
	}

	ps.clients[c] = struct{}{}
	return true
}

func (ps *ProxyServer) untrackClient(c net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.clients, c)
}

func (ps *ProxyServer) Shutdown() {
	ps.shutdownOnce.Do(func() {
		ps.closing.Store(true)

		ps.mu.Lock()
		for c := range ps.clients {
			c.Close()
		}
		ps.clients = make(map[net.Conn]struct{})
		ps.listener.Close()
		ps.mu.Unlock()

		<-ps.done
		log.Println("proxy server is shut down")
	})
}

func (ps *ProxyServer) handle(client net.Conn) {
	defer client.Close()
	defer ps.untrackClient(client)

	if err := client.SetReadDeadline(time.Now().Add(AuthReadTimeout)); err != nil {
		log.Println("auth deadline:", err)
		return
	}

	reader := bufio.NewReader(client)
	token, err := readAuthToken(reader)
	if err != nil {
		log.Println("auth read:", err)
		return
	}

	targetAddr, err := ps.resolver.ResolveTarget(token)
	if err != nil {
		_, _ = io.WriteString(client, ResponseERR)
		log.Println("auth failed:", err)
		return
	}

	if err := client.SetReadDeadline(time.Time{}); err != nil {
		log.Println("clear deadline:", err)
		return
	}

	if _, err := io.WriteString(client, ResponseOK); err != nil {
		log.Println("auth response:", err)
		return
	}

	backend, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("failed to dial target:", err)
		return
	}
	defer backend.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, _ = io.Copy(backend, reader)
		done <- struct{}{}
	}()

	go func() {
		_, _ = io.Copy(client, backend)
		done <- struct{}{}
	}()

	<-done
	<-done
}

func readAuthToken(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(line)
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}
