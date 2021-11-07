package send

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
	emmail "github.com/emersion/go-message/mail"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"golang.org/x/text/transform"
	"xorkevin.dev/mailcat/transformer"
)

type (
	Opts struct {
		Addr         string
		Username     string
		Password     string
		From         string
		To           string
		DKIMSelector string
	}

	Sender interface {
		ReadMsg(r io.Reader) error
		Send(addr string, username, password string, from, to string, dkimSelector string) error
	}

	sender struct {
		m              *message.Entity
		fromAddr       string
		fromAddrDomain string
		headers        []string
	}
)

func Send(r io.Reader, opts Opts) error {
	s := New()
	if err := s.ReadMsg(r); err != nil {
		return err
	}
	if err := s.Send(opts.Addr, opts.Username, opts.Password, opts.From, opts.To, opts.DKIMSelector); err != nil {
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
	headerReferences  = "References"
	headerContentType = "Content-Type"
	headerMsgID       = "Message-ID"
	headerDate        = "Date"
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
	s.headers = make([]string, 0, 10)
	s.headers = append(s.headers, headerMsgID)
	if headers.Has(headerDate) {
		if _, err := headers.Date(); err != nil {
			return fmt.Errorf("Invalid Date: %w", err)
		}
		s.headers = append(s.headers, headerDate)
	}
	if addrs, err := headers.AddressList(headerFrom); err != nil {
		return fmt.Errorf("Invalid From: %w", err)
	} else if len(addrs) > 1 {
		return fmt.Errorf("%w: multiple From", ErrInvalidHeader)
	} else if len(addrs) == 0 {
		return fmt.Errorf("%w: no From", ErrInvalidHeader)
	} else {
		fromAddr := addrs[0].Address
		fromAddrParts := strings.Split(fromAddr, "@")
		fromAddrDomain := fromAddrParts[1]
		if len(fromAddrParts) != 2 || fromAddrParts[0] == "" || fromAddrDomain == "" {
			return fmt.Errorf("%w: Invalid from header", ErrInvalidHeader)
		}
		s.fromAddr = fromAddr
		s.fromAddrDomain = fromAddrDomain
	}
	s.headers = append(s.headers, headerFrom)
	if addrs, err := headers.AddressList(headerTo); err != nil {
		return fmt.Errorf("Invalid To: %w", err)
	} else if len(addrs) == 0 {
		return fmt.Errorf("%w: no To", ErrInvalidHeader)
	}
	s.headers = append(s.headers, headerTo)
	if headers.Has(headerCc) {
		if addrs, err := headers.AddressList(headerCc); err != nil {
			return fmt.Errorf("Invalid Cc: %w", err)
		} else if len(addrs) == 0 {
			return fmt.Errorf("%w: empty Cc", ErrInvalidHeader)
		}
		s.headers = append(s.headers, headerCc)
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
	s.headers = append(s.headers, headerSubject)
	if headers.Has(headerReplyTo) {
		if addrs, err := headers.AddressList(headerReplyTo); err != nil {
			return fmt.Errorf("Invalid Reply-To: %w", err)
		} else if len(addrs) > 1 {
			return fmt.Errorf("%w: multiple Reply-To", ErrInvalidHeader)
		} else if len(addrs) == 0 {
			return fmt.Errorf("%w: empty Reply-To", ErrInvalidHeader)
		}
		s.headers = append(s.headers, headerReplyTo)
	}
	if headers.Has(headerInReplyTo) {
		if replies, err := headers.MsgIDList(headerInReplyTo); err != nil {
			return fmt.Errorf("Invalid In-Reply-To: %w", err)
		} else if len(replies) > 1 {
			return fmt.Errorf("%w: multiple In-Reply-To", ErrInvalidHeader)
		} else if len(replies) == 0 {
			return fmt.Errorf("%w: empty In-Reply-To", ErrInvalidHeader)
		}
		s.headers = append(s.headers, headerInReplyTo)
	}
	if headers.Has(headerReferences) {
		if replies, err := headers.MsgIDList(headerReferences); err != nil {
			return fmt.Errorf("Invalid References: %w", err)
		} else if len(replies) == 0 {
			return fmt.Errorf("%w: empty References", ErrInvalidHeader)
		}
		s.headers = append(s.headers, headerReferences)
	}
	if headers.Has(headerContentType) {
		if _, _, err := headers.ContentType(); err != nil {
			return fmt.Errorf("Invalid Content-Type: %w", err)
		}
		s.headers = append(s.headers, headerContentType)
	}
	s.m = m
	return nil
}

const (
	durationMonth = 30 * 24 * time.Hour
)

func (s *sender) sign(w io.Writer, r io.Reader) error {
	if err := dkim.Sign(w, r, &dkim.SignOptions{
		Domain:                 s.fromAddrDomain,
		Selector:               "tests",
		Identifier:             s.fromAddr,
		Signer:                 nil,
		Hash:                   crypto.SHA256,
		HeaderCanonicalization: dkim.CanonicalizationRelaxed,
		BodyCanonicalization:   dkim.CanonicalizationRelaxed,
		HeaderKeys:             s.headers,
		Expiration:             time.Now().Round(0).Add(durationMonth),
		QueryMethods:           []dkim.QueryMethod{dkim.QueryMethodDNSTXT},
	}); err != nil {
		return fmt.Errorf("Failed to dkim sign message: %w", err)
	}
	return nil
}

func (s *sender) Send(addr string, username, password string, from, to string, dkimSelector string) error {
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
	if dkimSelector != "" {
		t := &bytes.Buffer{}
		if err := s.sign(t, b); err != nil {
			return err
		}
		b = t
	}
	var auth sasl.Client
	if username != "" {
		auth = sasl.NewPlainClient("", username, password)
	}
	if err := smtp.SendMail(addr, auth, from, []string{to}, b); err != nil {
		return fmt.Errorf("Failed to send mail: %w", err)
	}
	return nil
}
