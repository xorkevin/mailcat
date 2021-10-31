package formatter

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/transform"
)

func Test_lfTransformer(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name string
		Inp  string
		Exp  string
	}{
		{
			Name: "Multiple CR-LF -> LF",
			Inp:  "This is\r\n\r\nsome text\n\rfor\ra\ntest\r",
			Exp:  "This is\n\nsome text\n\rfor\ra\ntest\r",
		},
	} {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			b, err := io.ReadAll(transform.NewReader(strings.NewReader(tc.Inp), lfTransformer{}))
			assert.NoError(err)
			assert.Equal([]byte(tc.Exp), b)
			b2, err := io.ReadAll(transform.NewReader(bytes.NewReader(b), lfTransformer{}))
			assert.NoError(err)
			assert.Equal([]byte(tc.Exp), b2)
		})
	}
}

func Test_crlfTransformer(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name string
		Inp  string
		Exp  string
	}{
		{
			Name: "Multiple LF -> CR-LF",
			Inp:  "This is\r\n\nsome text\n\r\nfor\ra\ntest\r",
			Exp:  "This is\r\n\r\nsome text\r\n\r\nfor\ra\r\ntest\r",
		},
	} {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			b, err := io.ReadAll(transform.NewReader(strings.NewReader(tc.Inp), crlfTransformer{}))
			assert.NoError(err)
			assert.Equal([]byte(tc.Exp), b)
			b2, err := io.ReadAll(transform.NewReader(bytes.NewReader(b), crlfTransformer{}))
			assert.NoError(err)
			assert.Equal([]byte(tc.Exp), b2)
		})
	}
}
