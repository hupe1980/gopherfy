package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	http := NewHTTP()

	assert.Equal(t, "gopher://127.0.0.1:80/_GET%20/%20HTTP/1.0%0A%0A", http.Payload())
}
