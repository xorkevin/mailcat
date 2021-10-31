package uid

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type (
	// Snowflake is a uid approximately sortable by time
	Snowflake struct {
		u []byte
	}
)

const (
	timeSize = 8
)

// NewSnowflake creates a new snowflake uid
func NewSnowflake(randsize int) (*Snowflake, error) {
	u := make([]byte, timeSize+randsize)
	now := uint64(time.Now().Round(0).UnixMilli())
	binary.BigEndian.PutUint64(u[:timeSize], now)
	_, err := rand.Read(u[timeSize:])
	if err != nil {
		return nil, fmt.Errorf("Failed reading crypto/rand: %w", err)
	}
	return &Snowflake{
		u: u,
	}, nil
}

// Bytes returns the full raw bytes of a snowflake
func (s *Snowflake) Bytes() []byte {
	return s.u
}

var (
	base32RawHexEncoding = base32.HexEncoding.WithPadding(base32.NoPadding)
)

// Base32 returns the full raw bytes of a snowflake in unpadded base32hex
func (s *Snowflake) Base32() string {
	return strings.ToLower(base32RawHexEncoding.EncodeToString(s.u))
}
