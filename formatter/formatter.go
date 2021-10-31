package formatter

import (
	"fmt"
	"io"
	"os"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"golang.org/x/text/transform"
	"xorkevin.dev/mailcat/transformer"
)

type (
	Formatter interface {
		ReadMessage(r io.Reader) error
	}

	formatter struct {
	}
)

func Format(r io.Reader) error {
	f := New()
	return f.ReadMessage(r)
}

func New() Formatter {
	return &formatter{}
}

const (
	headerFrom = "From"
)

func (f *formatter) ReadMessage(r io.Reader) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.Read(r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	headers := emmail.Header{
		Header: m.Header,
	}
	if fromAddrs, err := headers.AddressList(headerFrom); err != nil || len(fromAddrs) == 0 {
		headers.SetAddressList(headerFrom, []*emmail.Address{
			{Name: "Name", Address: "mail@example.com"},
		})
	}
	if err := m.WriteTo(transform.NewWriter(os.Stdout, transformer.LF{})); err != nil {
		return fmt.Errorf("Failed writing mail message: %w", err)
	}
	return nil
}
