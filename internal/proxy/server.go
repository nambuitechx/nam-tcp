package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type ProxyServer struct {
	listener     net.Listener
	target       string
	clients      map[net.Conn]struct{}
	mu           sync.Mutex
	done         chan struct{}
	closing      atomic.Bool
	shutdownOnce sync.Once
}

func NewProxyServer(address, target string) (*ProxyServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %w", address, err)
	}

	return &ProxyServer{
		listener: listener,
		target:   target,
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

	target, err := net.Dial("tcp", ps.target)
	if err != nil {
		log.Println("failed to dial target:", err)
		return
	}
	defer target.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, _ = io.Copy(target, client)
		done <- struct{}{}
	}()

	go func() {
		_, _ = io.Copy(client, target)
		done <- struct{}{}
	}()

	<-done
	<-done
}
