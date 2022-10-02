package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	smtp := NewSMTP(func(o *Options) {
		o.From = "from@mail.org"
		o.To = "to@mail.org"
		o.Subject = "subject"
		o.Msg = "msg"
	})

	assert.Equal(t, "gopher://127.0.0.1:25/_HELO%20localhost%0AMAIL%20FROM:%3Cfrom%40mail.org%3E%0ARCPT%20TO:%3Cto%40mail.org%3E%0ADATA%0ASubject:subject%0Amsg%0A.", smtp.Payload())
}
