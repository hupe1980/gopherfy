package gopher

import (
	"fmt"
	"io"
	"net"
	urlpkg "net/url"
)

const (
	// CRLF is the delimiter used per line of response item
	CRLF = "\r\n"
)

type Response struct {
	Body io.ReadCloser
}

// Get fetches a Gopher resource by URI
func Get(url string) (*Response, error) {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "gopher" {
		return nil, fmt.Errorf("invalid scheme for url: %s", u.Scheme)
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %s", u.Host)
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(u.Path[2:] + CRLF))
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Response{Body: conn}, nil
}
