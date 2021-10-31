package formatter

import (
	"fmt"
	"io"
	"os"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"golang.org/x/text/transform"
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

type (
	lfTransformer struct{}

	crlfTransformer struct{}
)

func (t lfTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, tErr error) {
	for nDst < len(dst) && nSrc < len(src) {
		c := src[nSrc]
		if c == '\r' {
			if nSrc+1 < len(src) {
				if src[nSrc+1] == '\n' {
					nSrc++
					continue
				}
			} else if !atEOF {
				tErr = transform.ErrShortSrc
				return
			}
		}
		dst[nDst] = c
		nDst++
		nSrc++
	}
	if nSrc < len(src) {
		tErr = transform.ErrShortDst
	}
	return
}

func (t lfTransformer) Reset() {}

func (t crlfTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, tErr error) {
	for nDst < len(dst) && nSrc < len(src) {
		c := src[nSrc]
		if c == '\r' {
			if nSrc+1 < len(src) {
				if src[nSrc+1] == '\n' {
					if nDst+1 < len(dst) {
						dst[nDst] = '\r'
						dst[nDst+1] = '\n'
						nDst += 2
						nSrc += 2
						continue
					}
					break
				}
			} else if !atEOF {
				tErr = transform.ErrShortSrc
				return
			}
		} else if c == '\n' {
			if nDst+1 < len(dst) {
				dst[nDst] = '\r'
				dst[nDst+1] = '\n'
				nDst += 2
				nSrc++
				continue
			}
			break
		}
		dst[nDst] = c
		nDst++
		nSrc++
	}
	if nSrc < len(src) {
		tErr = transform.ErrShortDst
	}
	return
}

func (t crlfTransformer) Reset() {}

const (
	headerFrom = "From"
)

func (f *formatter) ReadMessage(r io.Reader) error {
	r = transform.NewReader(r, crlfTransformer{})
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
	if err := m.WriteTo(transform.NewWriter(os.Stdout, lfTransformer{})); err != nil {
		return fmt.Errorf("Failed writing mail message: %w", err)
	}
	return nil
}
