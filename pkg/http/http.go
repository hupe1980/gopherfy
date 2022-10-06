package http

import (
	"bytes"
	"fmt"
	gohttp "net/http"
	"strings"

	"github.com/hupe1980/gopherfy/internal"
)

const (
	DefaultAddr      = "127.0.0.1:80"
	DefaultMethod    = "GET"
	DefaultVersion   = "HTTP/1.0"
	DefaultPath      = "/"
	DefaultUserAgent = "gopherfy"
)

type Options struct {
	Addr         string
	Method       string
	Version      string
	Path         string
	UserAgent    string
	ExtraHeaders map[string]string
	NewLine      string
}

type HTTP struct {
	addr         string
	method       string
	version      string
	path         string
	userAgent    string
	extraHeaders map[string]string
	newLine      string
}

func NewHTTP(optFns ...func(o *Options)) *HTTP {
	options := Options{
		Addr:      DefaultAddr,
		Method:    DefaultMethod,
		Version:   DefaultVersion,
		Path:      DefaultPath,
		UserAgent: DefaultUserAgent,
		NewLine:   "%0D%0A",
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &HTTP{
		addr:         strings.TrimSpace(options.Addr),
		method:       strings.ToUpper(strings.TrimSpace(options.Method)),
		version:      strings.TrimSpace(options.Version),
		path:         internal.URLEncode(strings.TrimSpace(options.Path)),
		userAgent:    strings.TrimSpace(options.UserAgent),
		extraHeaders: options.ExtraHeaders,
		newLine:      options.NewLine,
	}
}

func (http *HTTP) Payload() string {
	payload := fmt.Sprintf("%s%%20%s%%20%s%s", http.method, http.path, http.version, http.newLine)

	headers := http.generateHeaders()
	if headers != "" {
		payload = fmt.Sprintf("%s%s", payload, headers)
	}

	return fmt.Sprintf("gopher://%s/_%s", http.addr, payload)
}

func (http *HTTP) generateHeaders() string {
	headers := internal.NewInsertionOrderMap()
	headers.Set("Host", http.addr)
	headers.Set("User-Agent", http.userAgent)
	headers.Set("Accept", "*/*")

	for k, v := range http.extraHeaders {
		headers.Set(k, v)
	}

	b := new(bytes.Buffer)

	for _, k := range headers.Keys() {
		v, found := headers.Get(k)
		if !found || v == "" {
			continue
		}

		fmt.Fprintf(b, "%s:%%20%s%s", gohttp.CanonicalHeaderKey(k), v, http.newLine)
	}

	return b.String()
}
