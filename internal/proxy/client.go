package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// AuthConn is a TCP connection to the proxy after PAT authentication.
// Reads use an internal buffer so bytes sent with the auth line are not lost.
type AuthConn struct {
	net.Conn
	reader *bufio.Reader
}

func (c *AuthConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

// Dial authenticates with the proxy and returns a connection ready for app traffic.
func Dial(proxyAddr, token string) (*AuthConn, error) {
	conn, err := net.DialTimeout("tcp", proxyAddr, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}

	reader := bufio.NewReader(conn)
	if err := conn.SetWriteDeadline(time.Now().Add(AuthReadTimeout)); err != nil {
		conn.Close()
		return nil, err
	}

	if _, err := io.WriteString(conn, strings.TrimSpace(token)+"\n"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("write token: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Time{}); err != nil {
		conn.Close()
		return nil, err
	}

	if err := conn.SetReadDeadline(time.Now().Add(AuthReadTimeout)); err != nil {
		conn.Close()
		return nil, err
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("read auth response: %w", err)
	}

	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		conn.Close()
		return nil, err
	}

	switch strings.TrimSpace(line) {
	case "OK":
		return &AuthConn{Conn: conn, reader: reader}, nil
	case "ERR":
		conn.Close()
		return nil, errors.New("proxy rejected pat")
	default:
		conn.Close()
		return nil, fmt.Errorf("unexpected auth response: %q", strings.TrimSpace(line))
	}
}
