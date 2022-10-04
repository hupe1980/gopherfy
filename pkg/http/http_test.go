package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	http := NewHTTP()

	assert.Equal(t, "gopher://127.0.0.1:80/_GET%20%2F%20HTTP/1.0%0D%0AHost:%20127.0.0.1:80%0D%0AUser-Agent:%20gopherfy%0D%0AAccept:%20*/*%0D%0A", http.Payload())
}
