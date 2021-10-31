package formatter

import (
	"io"

	_ "github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	_ "github.com/emersion/go-message/mail"
	"golang.org/x/text/transform"
)

type (
	Formatter interface {
	}

	formatter struct {
	}
)

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

func (f *formatter) ReadMessage(inp io.Reader) error {
	return nil
}
