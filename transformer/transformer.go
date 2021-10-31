package transformer

import (
	"golang.org/x/text/transform"
)

type (
	LF struct{}

	CRLF struct{}
)

func (t LF) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, tErr error) {
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

func (t LF) Reset() {}

func (t CRLF) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, tErr error) {
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

func (t CRLF) Reset() {}
