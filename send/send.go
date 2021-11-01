package send

import (
	"errors"
	"fmt"
	"io"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"golang.org/x/text/transform"
	"xorkevin.dev/mailcat/transformer"
)

type (
	Opts struct {
	}

	Sender interface {
		ReadMsg(r io.Reader) error
	}

	sender struct {
		m *message.Entity
	}
)

func Send(r io.Reader, opts Opts) error {
	fmt.Println("hello world")
	return nil
}

func New() Sender {
	return &sender{}
}

var (
	ErrInvalidHeader = errors.New("Invalid header")
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
