package formatter

import (
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strings"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"golang.org/x/text/transform"
	"xorkevin.dev/mailcat/transformer"
	"xorkevin.dev/mailcat/uid"
)

type (
	Opts struct {
		CRLF        bool
		Body        bool
		Headers     []string
		AddHeaders  []string
		MsgIDDomain string
	}

	Formatter interface {
		SetHeaders(opts Opts) error
		ReadBody(r io.Reader, opts Opts) error
		ReadMsg(r io.Reader, opts Opts) error
		WriteMsg(r io.Writer, opts Opts) error
	}

	formatter struct {
		m *message.Entity
	}
)

func Format(r io.Reader, w io.Writer, opts Opts) error {
	f := New()
	if opts.Body {
		if err := f.ReadBody(r, opts); err != nil {
			return err
		}
	} else {
		if err := f.ReadMsg(r, opts); err != nil {
			return err
		}
	}
	if err := f.SetHeaders(opts); err != nil {
		return err
	}
	if err := f.WriteMsg(w, opts); err != nil {
		return err
	}
	return nil
}

func New() Formatter {
	return &formatter{}
}

var (
	ErrNoMsg         = errors.New("No mail message read")
	ErrInvalidHeader = errors.New("Invalid header")
)

const (
	msgidRandBytes = 16
)

const (
	headerFrom      = "From"
	headerTo        = "To"
	headerCc        = "Cc"
	headerBcc       = "Bcc"
	headerSubject   = "Subject"
	headerReplyTo   = "Reply-To"
	headerInReplyTo = "In-Reply-To"
)

const (
	contentTypeTextPlain = "text/plain"
)

func (f *formatter) genMsgID(opts Opts) (string, error) {
	u, err := uid.NewSnowflake(msgidRandBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to generate msgid: %w", err)
	}
	return fmt.Sprintf("%s@%s", u.Base32(), opts.MsgIDDomain), nil
}

func (f *formatter) SetHeaders(opts Opts) error {
	if f.m == nil {
		return ErrNoMsg
	}
	headers := emmail.Header{
		Header: f.m.Header,
	}
	// headers are in reverse order of appearance since headers are prepended
	for _, i := range opts.Headers {
		parts := strings.SplitN(i, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%w: %s", ErrInvalidHeader, i)
		}
		k := textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(parts[0]))
		if !headers.Has(k) {
			v := strings.TrimSpace(parts[1])
			headers.Set(k, v)
		}
	}
	for _, i := range opts.AddHeaders {
		parts := strings.SplitN(i, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%w: %s", ErrInvalidHeader, i)
		}
		k := textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(parts[0]))
		v := strings.TrimSpace(parts[1])
		headers.Add(k, v)
	}
	if t, params, err := headers.ContentType(); err != nil {
		return fmt.Errorf("Invalid Content-Type: %w", err)
	} else {
		headers.SetContentType(t, params)
	}
	if replies, err := headers.MsgIDList(headerInReplyTo); err != nil {
		return fmt.Errorf("Invalid In-Reply-To: %w", err)
	} else if len(replies) > 1 {
		return fmt.Errorf("%w: multiple In-Reply-To", ErrInvalidHeader)
	} else {
		headers.SetMsgIDList(headerInReplyTo, replies)
	}
	if msgid, err := headers.MessageID(); err != nil {
		return fmt.Errorf("Invalid Message-ID: %w", err)
	} else {
		headers.SetMessageID(msgid)
	}
	if addrs, err := headers.AddressList(headerReplyTo); err != nil {
		return fmt.Errorf("Invalid Reply-To: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple Reply-To", ErrInvalidHeader)
	} else {
		headers.SetAddressList(headerReplyTo, addrs)
	}
	if subj, err := headers.Subject(); err != nil {
		return fmt.Errorf("Invalid Subject: %w", err)
	} else {
		headers.SetSubject(subj)
	}
	for _, i := range []string{headerBcc, headerCc} {
		if addrs, err := headers.AddressList(i); err != nil {
			return fmt.Errorf("Invalid %s: %w", i, err)
		} else {
			headers.SetAddressList(i, addrs)
		}
	}
	if addrs, err := headers.AddressList(headerTo); err != nil {
		return fmt.Errorf("Invalid To: %w", err)
	} else if len(addrs) == 0 {
		headers.SetAddressList(headerTo, []*emmail.Address{
			{Name: "Name", Address: "mail@example.com"},
		})
	} else {
		headers.SetAddressList(headerTo, addrs)
	}
	if addrs, err := headers.AddressList(headerFrom); err != nil {
		return fmt.Errorf("Invalid From: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple From", ErrInvalidHeader)
	} else if len(addrs) == 0 {
		headers.SetAddressList(headerFrom, []*emmail.Address{
			{Name: "Name", Address: "mail@example.com"},
		})
	} else {
		headers.SetAddressList(headerFrom, addrs)
	}
	f.m.Header = headers.Header
	return nil
}

func (f *formatter) ReadBody(r io.Reader, opts Opts) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.New(message.Header{}, r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	f.m = m
	return nil
}

func (f *formatter) ReadMsg(r io.Reader, opts Opts) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.Read(r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	f.m = m
	return nil
}

func (f *formatter) WriteMsg(w io.Writer, opts Opts) error {
	if f.m == nil {
		return ErrNoMsg
	}
	if !opts.CRLF {
		w = transform.NewWriter(w, transformer.LF{})
	}
	if err := f.m.WriteTo(w); err != nil {
		return fmt.Errorf("Failed writing mail message: %w", err)
	}
	return nil
}
