package ulid

import (
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"time"
	"unsafe"
)

// ULID is a 128-bit Universally Unique Lexicographically Sortable Identifier.
type ULID [16]byte

const (
	// Encodingis Crockford's Base32
	// https://www.crockford.com/base32.html
	Encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// Size is the length of a ULID in bytes.
	Size = len(ULID{})

	// EncodedSize is the length of a ULID in chars.
	EncodedSize = (Size*8 + 4) / 5
)

var (
	// ErrInvalidLength is returned when the length of the ULID is invalid.
	ErrInvalidLength = errors.New("ulid: invalid length")

	// ErrInvalidChar is returned when the char is invalid.
	ErrInvalidChar = errors.New("ulid: invalid character")

	// ErrInvalidULID is returned when the ULID is invalid.
	ErrInvalidULID = errors.New("ulid: invalid ULID")

	// ErrDataSize is returned when the data size is invalid.
	ErrDataSize = errors.New("ulid: invalid data size")

	// base32enc is the base32 encoding.
	base32enc = base32.NewEncoding(Encoding).WithPadding(base32.NoPadding)
)

// MonotonicReader is a reader of monotonic time.
type MonotonicReader interface {
	MonotonicRead(ms int64, p []byte) error
	io.Reader
}

// Max returns the max ULID.
func Max() ULID {
	return ULID{
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
	}
}

// New returns a new ULID.
func New() (ULID, error) {
	var u ULID

	return u, nil
}

// MustNew returns a new ULID.
func MustNew() ULID {
	u, err := New()
	if err != nil {
		panic(err)
	}

	return u
}

// Parse parses a ULID from a byte slice.
func Parse(data []byte) (ULID, error) {
	var u ULID
	err := u.UnmarshalBinary(data)

	return u, err
}

// ParseString parses a ULID from a string.
func ParseString(s string) (ULID, error) {
	var u ULID
	err := u.UnmarshalText([]byte(s))

	return u, err
}

// Time returns the time component of the ULID.
func (u ULID) Time() uint64 {
	return uint64(u[5]) | uint64(u[4])<<8 | uint64(u[3])<<16 |
		uint64(u[2])<<24 | uint64(u[1])<<32 | uint64(u[0])<<40
}

// maxTime returns the max time component of the ULID.
func maxTime() uint64 {
	return ULID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0, 0}.Time()
}

// MaxTime returns the max time component of the ULID.
func MaxTime() uint64 {
	return maxTime()
}

// Now returns a ULID with the given time.
func Now() int64 {
	return time.Now().UnixMilli()
}

// AppendFormat append text encoded ULID to 'dst'.
func (u ULID) AppendFormat(dst []byte) []byte {
	length := len(dst)

	if cap(dst)-length < Size {
		// Allocate larger slice for ULID
		capac := 2*cap(dst) + Size
		dst2 := make([]byte, length, capac)
		copy(dst2, dst)
		dst = dst2
	}

	dst = dst[:length+Size]

	into := dst[length : length+Size]
	base32enc.Encode(into, u[:])

	return dst
}

// Bytes returns text encoded bytes of receiving ULID.
func (u ULID) Bytes() []byte {
	dst := make([]byte, 0, Size)

	return u.AppendFormat(dst)
}

// String returns the string representation of the ULID.
func (u ULID) String() string {
	b := u.Bytes()
	p := unsafe.Pointer(&b)

	return *(*string)(p)
}

// UnmarshalBinary decodes a ULID from binary form.
func (u *ULID) UnmarshalBinary(data []byte) error {
	if len(data) != len(*u) {
		return ErrDataSize
	}

	copy(data, u[:])

	return nil
}

// UnmarshalText decodes a ULID from text form.
func (u *ULID) UnmarshalText(b []byte) error {
	if len(b) != Size {
		return ErrInvalidLength
	}

	switch c := b[Size-1]; c {
	case '0', '4', '8', 'C', 'G', 'M', 'R', 'W':
	default:
		return fmt.Errorf("%w: '%c' outside encoding range", ErrInvalidChar, c)
	}

	// Base32 decode into receiving ULID.
	_, err := base32enc.Decode(u[:], b)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidULID, err)
	}

	return nil
}

// MarshalText encodes a ULID into text form by implementing encoding.TextMarshaler.
func (u ULID) MarshalText() ([]byte, error) {
	return u.Bytes(), nil
}

// MarshalBinary encodes a ULID into binary form.
func (u ULID) MarshalBinary() ([]byte, error) {
	return u[:], nil
}

// MarshalBinaryTo encodes a ULID into binary form.
func (u ULID) MarshalBinaryTo(data []byte) error {
	if len(data) != len(u) {
		return ErrDataSize
	}

	copy(data, u[:])

	return nil
}
