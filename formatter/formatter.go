package formatter

import (
	"fmt"
	"io"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"golang.org/x/text/transform"
	"xorkevin.dev/mailcat/transformer"
	"xorkevin.dev/mailcat/uid"
)

type (
	Opts struct {
		CRLF bool
	}

	Formatter interface {
		ReadMsg(r io.Reader, opts Opts) error
		WriteMsg(r io.Writer, opts Opts) error
	}

	formatter struct {
		m *message.Entity
	}
)

func Format(r io.Reader, w io.Writer, opts Opts) error {
	f := New()
	if err := f.ReadMsg(r, opts); err != nil {
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

func (f *formatter) ReadMsg(r io.Reader, opts Opts) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.Read(r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	headers := emmail.Header{
		Header: m.Header,
	}
	// headers are in reverse order of appearance since headers are prepended
	if t, params, err := headers.ContentType(); err != nil || t == "" {
		headers.SetContentType(contentTypeTextPlain, nil)
	} else {
		headers.SetContentType(t, params)
	}
	if replies, err := headers.MsgIDList(headerInReplyTo); err != nil || len(replies) != 1 {
		headers.SetMsgIDList(headerInReplyTo, nil)
	} else {
		headers.SetMsgIDList(headerInReplyTo, replies)
	}
	if msgid, err := headers.MessageID(); err != nil || msgid == "" {
		u, err := uid.NewSnowflake(msgidRandBytes)
		if err != nil {
			return fmt.Errorf("Failed to generate msgid: %w", err)
		}
		headers.SetMessageID(fmt.Sprintf("%s@%s", u.Base32(), "mail.example.com"))
	} else {
		headers.SetMessageID(msgid)
	}
	if addrs, err := headers.AddressList(headerReplyTo); err != nil || len(addrs) != 1 {
		headers.SetAddressList(headerReplyTo, nil)
	} else {
		headers.SetAddressList(headerReplyTo, addrs)
	}
	if subj, err := headers.Subject(); err != nil {
		headers.SetSubject("")
	} else {
		headers.SetSubject(subj)
	}
	for _, i := range []string{headerBcc, headerCc} {
		if addrs, err := headers.AddressList(i); err != nil || len(addrs) == 0 {
			headers.SetAddressList(i, nil)
		} else {
			headers.SetAddressList(i, addrs)
		}
	}
	if addrs, err := headers.AddressList(headerTo); err != nil || len(addrs) == 0 {
		headers.SetAddressList(headerTo, []*emmail.Address{
			{Name: "Name", Address: "mail@example.com"},
		})
	} else {
		headers.SetAddressList(headerTo, addrs)
	}
	if addrs, err := headers.AddressList(headerFrom); err != nil || len(addrs) != 1 {
		headers.SetAddressList(headerFrom, []*emmail.Address{
			{Name: "Name", Address: "mail@example.com"},
		})
	} else {
		headers.SetAddressList(headerFrom, addrs)
	}
	m.Header = headers.Header
	f.m = m
	return nil
}

func (f *formatter) WriteMsg(w io.Writer, opts Opts) error {
	if f.m == nil {
		return fmt.Errorf("No mail message read")
	}
	if !opts.CRLF {
		w = transform.NewWriter(w, transformer.LF{})
	}
	if err := f.m.WriteTo(w); err != nil {
		return fmt.Errorf("Failed writing mail message: %w", err)
	}
	return nil
}
