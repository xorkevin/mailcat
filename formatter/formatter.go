package formatter

import (
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strings"
	"time"

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
		Edit        bool
	}

	Formatter interface {
		SetHeaders(setHeaders, addHeaders []string) error
		SetHeadersFinal(msgidDomain string) error
		ReadBody(r io.Reader) error
		ReadMsg(r io.Reader) error
		WriteMsg(w io.Writer, crlf bool) error
	}

	formatter struct {
		m *message.Entity
	}
)

func Format(r io.Reader, w io.Writer, opts Opts) error {
	f := New()
	if opts.Body {
		if err := f.ReadBody(r); err != nil {
			return err
		}
	} else {
		if err := f.ReadMsg(r); err != nil {
			return err
		}
	}
	if err := f.SetHeaders(opts.Headers, opts.AddHeaders); err != nil {
		return err
	}
	if opts.Edit {
		if err := f.WriteMsg(w, false); err != nil {
			return err
		}
		return nil
	}
	if err := f.SetHeadersFinal(opts.MsgIDDomain); err != nil {
		return err
	}
	if err := f.WriteMsg(w, opts.CRLF); err != nil {
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
	ErrNoEditor      = errors.New("No editor found")
)

const (
	msgidRandBytes = 16
)

const (
	headerFrom        = "From"
	headerTo          = "To"
	headerCc          = "Cc"
	headerBcc         = "Bcc"
	headerSubject     = "Subject"
	headerReplyTo     = "Reply-To"
	headerInReplyTo   = "In-Reply-To"
	headerReferences  = "References"
	headerContentType = "Content-Type"
	headerMsgID       = "Message-ID"
	headerDate        = "Date"
)

const (
	contentTypeTextPlain = "text/plain"
)

func (f *formatter) genMsgID(msgidDomain string) (string, error) {
	u, err := uid.NewSnowflake(msgidRandBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to generate msgid: %w", err)
	}
	return fmt.Sprintf("%s@%s", u.Base32(), msgidDomain), nil
}

func (f *formatter) SetHeaders(setHeaders, addHeaders []string) error {
	if f.m == nil {
		return ErrNoMsg
	}
	headers := emmail.Header{
		Header: f.m.Header,
	}
	// headers are in reverse order of appearance since headers are prepended
	for _, i := range setHeaders {
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
	for _, i := range addHeaders {
		parts := strings.SplitN(i, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%w: %s", ErrInvalidHeader, i)
		}
		k := textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(parts[0]))
		v := strings.TrimSpace(parts[1])
		headers.Add(k, v)
	}
	if headers.Has(headerContentType) {
		if t, params, err := headers.ContentType(); err != nil {
			return fmt.Errorf("Invalid Content-Type: %w", err)
		} else {
			headers.SetContentType(t, params)
		}
	}
	if replies, err := headers.MsgIDList(headerReferences); err != nil {
		return fmt.Errorf("Invalid References: %w", err)
	} else {
		headers.SetMsgIDList(headerReferences, replies)
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
	} else if msgid == "" {
		headers.Del(headerMsgID)
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
	for _, i := range []string{headerBcc, headerCc, headerTo} {
		if addrs, err := headers.AddressList(i); err != nil {
			return fmt.Errorf("Invalid %s: %w", i, err)
		} else {
			headers.SetAddressList(i, addrs)
		}
	}
	if addrs, err := headers.AddressList(headerFrom); err != nil {
		return fmt.Errorf("Invalid From: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple From", ErrInvalidHeader)
	} else {
		headers.SetAddressList(headerFrom, addrs)
	}
	if headers.Has(headerDate) {
		t, err := headers.Date()
		if err != nil {
			return fmt.Errorf("Invalid Date: %w", err)
		}
		headers.SetDate(t)
	}
	f.m.Header = headers.Header
	return nil
}

func (f *formatter) SetHeadersFinal(msgidDomain string) error {
	if f.m == nil {
		return ErrNoMsg
	}
	headers := emmail.Header{
		Header: f.m.Header,
	}
	// headers are in reverse order of appearance since headers are prepended
	if headers.Has(headerContentType) {
		if t, params, err := headers.ContentType(); err != nil {
			return fmt.Errorf("Invalid Content-Type: %w", err)
		} else {
			headers.SetContentType(t, params)
		}
	}
	if refs, err := headers.MsgIDList(headerReferences); err != nil {
		return fmt.Errorf("Invalid References: %w", err)
	} else if len(refs) == 0 {
		headers.Del(headerReferences)
	} else {
		headers.SetMsgIDList(headerReferences, refs)
	}
	if replies, err := headers.MsgIDList(headerInReplyTo); err != nil {
		return fmt.Errorf("Invalid In-Reply-To: %w", err)
	} else if len(replies) > 1 {
		return fmt.Errorf("%w: multiple In-Reply-To", ErrInvalidHeader)
	} else if len(replies) == 0 {
		headers.Del(headerInReplyTo)
	} else {
		headers.SetMsgIDList(headerInReplyTo, replies)
	}
	if msgid, err := headers.MessageID(); err != nil {
		return fmt.Errorf("Invalid Message-ID: %w", err)
	} else if msgid == "" {
		id, err := f.genMsgID(msgidDomain)
		if err != nil {
			return err
		}
		headers.SetMessageID(id)
	} else {
		headers.SetMessageID(msgid)
	}
	if addrs, err := headers.AddressList(headerReplyTo); err != nil {
		return fmt.Errorf("Invalid Reply-To: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple Reply-To", ErrInvalidHeader)
	} else if len(addrs) == 0 {
		headers.Del(headerReplyTo)
	} else {
		headers.SetAddressList(headerReplyTo, addrs)
	}
	if subj, err := headers.Subject(); err != nil {
		return fmt.Errorf("Invalid Subject: %w", err)
	} else if subj == "" {
		return fmt.Errorf("%w: no Subject", ErrInvalidHeader)
	} else {
		headers.SetSubject(subj)
	}
	for _, i := range []string{headerBcc, headerCc} {
		if addrs, err := headers.AddressList(i); err != nil {
			return fmt.Errorf("Invalid %s: %w", i, err)
		} else if len(addrs) == 0 {
			headers.Del(i)
		} else {
			headers.SetAddressList(i, addrs)
		}
	}
	if addrs, err := headers.AddressList(headerTo); err != nil {
		return fmt.Errorf("Invalid To: %w", err)
	} else if len(addrs) == 0 {
		return fmt.Errorf("%w: no To", ErrInvalidHeader)
	} else {
		headers.SetAddressList(headerTo, addrs)
	}
	if addrs, err := headers.AddressList(headerFrom); err != nil {
		return fmt.Errorf("Invalid From: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple From", ErrInvalidHeader)
	} else if len(addrs) == 0 {
		return fmt.Errorf("%w: no From", ErrInvalidHeader)
	} else {
		headers.SetAddressList(headerFrom, addrs)
	}
	if headers.Has(headerDate) {
		t, err := headers.Date()
		if err != nil {
			return fmt.Errorf("Invalid Date: %w", err)
		}
		headers.SetDate(t)
	} else {
		headers.SetDate(time.Now().Round(0))
	}
	f.m.Header = headers.Header
	return nil
}

func (f *formatter) ReadBody(r io.Reader) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.New(message.Header{}, r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	f.m = m
	return nil
}

func (f *formatter) ReadMsg(r io.Reader) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.Read(r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	f.m = m
	return nil
}

func (f *formatter) WriteMsg(w io.Writer, crlf bool) error {
	if f.m == nil {
		return ErrNoMsg
	}
	if !crlf {
		w = transform.NewWriter(w, transformer.LF{})
	}
	if err := f.m.WriteTo(w); err != nil {
		return fmt.Errorf("Failed writing mail message: %w", err)
	}
	return nil
}
