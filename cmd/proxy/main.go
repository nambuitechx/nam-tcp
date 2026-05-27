package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	internal_http "github.com/nambuitechx/nam-tcp/internal/http"
	internal_proxy "github.com/nambuitechx/nam-tcp/internal/proxy"
)

func main() {
	httpAddr := envOr("HTTP_ADDR", ":8000")
	proxyAddr := envOr("PROXY_ADDR", ":8888")
	targetAddr := envOr("PROXY_TARGET", "localhost:8001")

	srv := internal_http.NewHttpServer()

	s := &http.Server{
		Addr:              httpAddr,
		Handler:           srv.Mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	go func() {
		log.Printf("http server listening on %s", httpAddr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	ps, err := internal_proxy.NewProxyServer(proxyAddr, targetAddr)
	if err != nil {
		log.Fatalf("proxy server: %v", err)
	}

	go func() {
		log.Printf("proxy server listening on %s -> %s", proxyAddr, targetAddr)
		ps.Run()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}

	ps.Shutdown()
	log.Println("shutdown complete")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
