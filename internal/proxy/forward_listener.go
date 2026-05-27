package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// ForwardListener runs a local TCP listener that can be shut down explicitly.
type ForwardListener struct {
	localAddr string
	proxyAddr string
	token     string

	listener net.Listener
	done     chan struct{}
	once     sync.Once
}

func NewForwardListener(localAddr, proxyAddr, token string) (*ForwardListener, error) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %w", localAddr, err)
	}

	return &ForwardListener{
		localAddr: localAddr,
		proxyAddr: proxyAddr,
		token:     token,
		listener:  listener,
		done:      make(chan struct{}),
	}, nil
}

func (f *ForwardListener) Addr() string {
	return f.listener.Addr().String()
}

func (f *ForwardListener) Run() error {
	defer close(f.done)
	log.Printf("local forward listening on %s -> proxy %s", f.localAddr, f.proxyAddr)

	for {
		localConn, err := f.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("accept: %w", err)
		}

		go bridgeThroughProxy(localConn, f.proxyAddr, f.token)
	}
}

func (f *ForwardListener) Close() error {
	var err error
	f.once.Do(func() {
		err = f.listener.Close()
	})
	return err
}

func (f *ForwardListener) Wait() {
	<-f.done
}

func bridgeThroughProxy(localConn net.Conn, proxyAddr, token string) {
	defer localConn.Close()

	remoteConn, err := Dial(proxyAddr, token)
	if err != nil {
		log.Printf("proxy dial failed from %s: %v", localConn.RemoteAddr(), err)
		return
	}
	defer remoteConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = io.Copy(remoteConn, localConn)
		if tc, ok := remoteConn.Conn.(*net.TCPConn); ok {
			_ = tc.CloseWrite()
		}
	}()

	go func() {
		defer wg.Done()
		_, _ = io.Copy(localConn, remoteConn)
		if tc, ok := localConn.(*net.TCPConn); ok {
			_ = tc.CloseWrite()
		}
	}()

	wg.Wait()
}
