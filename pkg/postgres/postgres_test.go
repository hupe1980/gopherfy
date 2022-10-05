package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	my := NewPostgres(func(o *Options) {
		o.DB = "demo"
		o.Query = "Select 1;"
	})

	assert.Equal(t, "gopher://127.0.0.1:5432/_%00%00%00%25%00%03%00%00%75%73%65%72%00%70%6f%73%74%67%72%65%73%00%64%61%74%61%62%61%73%65%00%64%65%6d%6f%00%00%51%00%00%00%0e%53%65%6c%65%63%74%20%31%3b%00%58%00%00%00%04", my.Payload())
}
