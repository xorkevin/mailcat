package send

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"golang.org/x/text/transform"
	"xorkevin.dev/mailcat/transformer"
)

type (
	Opts struct {
		Addr     string
		Username string
		Password string
		From     string
		To       string
	}

	Sender interface {
		ReadMsg(r io.Reader) error
		Send(addr string, username, password string, from, to string) error
	}

	sender struct {
		m *message.Entity
	}
)

func Send(r io.Reader, opts Opts) error {
	s := New()
	if err := s.ReadMsg(r); err != nil {
		return err
	}
	if err := s.Send(opts.Addr, opts.Username, opts.Password, opts.From, opts.To); err != nil {
		return err
	}
	return nil
}

func New() Sender {
	return &sender{}
}

var (
	ErrNoMsg         = errors.New("No mail message read")
	ErrInvalidHeader = errors.New("Invalid header")
	ErrInvalidArgs   = errors.New("Invalid args")
)

const (
	headerFrom        = "From"
	headerTo          = "To"
	headerCc          = "Cc"
	headerBcc         = "Bcc"
	headerSubject     = "Subject"
	headerReplyTo     = "Reply-To"
	headerInReplyTo   = "In-Reply-To"
	headerContentType = "Content-Type"
	headerMsgID       = "Message-ID"
)

func (s *sender) ReadMsg(r io.Reader) error {
	r = transform.NewReader(r, transformer.CRLF{})
	m, err := message.Read(r)
	if err != nil {
		return fmt.Errorf("Failed reading mail message: %w", err)
	}
	headers := emmail.Header{
		Header: m.Header,
	}
	if msgid, err := headers.MessageID(); err != nil {
		return fmt.Errorf("Invalid Message-ID: %w", err)
	} else if msgid == "" {
		return fmt.Errorf("%w: no Message-ID", ErrInvalidHeader)
	}
	if addrs, err := headers.AddressList(headerFrom); err != nil {
		return fmt.Errorf("Invalid From: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple From", ErrInvalidHeader)
	} else if len(addrs) == 0 {
		return fmt.Errorf("%w: no From", ErrInvalidHeader)
	}
	if addrs, err := headers.AddressList(headerTo); err != nil {
		return fmt.Errorf("Invalid To: %w", err)
	} else if len(addrs) == 0 {
		return fmt.Errorf("%w: no To", ErrInvalidHeader)
	}
	if headers.Has(headerCc) {
		if addrs, err := headers.AddressList(headerCc); err != nil {
			return fmt.Errorf("Invalid Cc: %w", err)
		} else if len(addrs) == 0 {
			return fmt.Errorf("%w: empty Cc", ErrInvalidHeader)
		}
	}
	if headers.Has(headerBcc) {
		if _, err := headers.AddressList(headerBcc); err != nil {
			return fmt.Errorf("Invalid Bcc: %w", err)
		}
		headers.Del(headerBcc)
	}
	if subj, err := headers.Subject(); err != nil {
		return fmt.Errorf("Invalid Subject: %w", err)
	} else if subj == "" {
		return fmt.Errorf("%w: no Subject", ErrInvalidHeader)
	}
	if headers.Has(headerReplyTo) {
		if addrs, err := headers.AddressList(headerReplyTo); err != nil {
			return fmt.Errorf("Invalid Reply-To: %w", err)
		} else if len(addrs) > 1 {
			return fmt.Errorf("%w: multiple Reply-To", ErrInvalidHeader)
		} else if len(addrs) == 0 {
			return fmt.Errorf("%w: empty Reply-To", ErrInvalidHeader)
		}
	}
	if headers.Has(headerInReplyTo) {
		if replies, err := headers.MsgIDList(headerInReplyTo); err != nil {
			return fmt.Errorf("Invalid In-Reply-To: %w", err)
		} else if len(replies) > 1 {
			return fmt.Errorf("%w: multiple In-Reply-To", ErrInvalidHeader)
		} else if len(replies) == 0 {
			return fmt.Errorf("%w: empty In-Reply-To", ErrInvalidHeader)
		}
	}
	if headers.Has(headerContentType) {
		if _, _, err := headers.ContentType(); err != nil {
			return fmt.Errorf("Invalid Content-Type: %w", err)
		}
	}
	s.m = m
	return nil
}

func (s *sender) Send(addr string, username, password string, from, to string) error {
	if s.m == nil {
		return ErrNoMsg
	}
	if addr == "" {
		return fmt.Errorf("%w: no address", ErrInvalidArgs)
	}
	if from == "" {
		return fmt.Errorf("%w: no smtp from", ErrInvalidArgs)
	}
	if to == "" {
		return fmt.Errorf("%w: no smtp to", ErrInvalidArgs)
	}
	b := &bytes.Buffer{}
	if err := s.m.WriteTo(b); err != nil {
		return fmt.Errorf("Failed to write mail message: %w", err)
	}
	auth := sasl.NewPlainClient("", username, password)
	if err := smtp.SendMail(addr, auth, from, []string{to}, b); err != nil {
		return fmt.Errorf("Failed to send mail: %w", err)
	}
	return nil
}
