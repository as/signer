package signer

import "encoding/base64"

var (
	codec = base64.RawURLEncoding
)

// Token is a byte slice that knows how to marshal and unmarshal itself in base64
type Token []byte

// String returns a url-safe base64-encoded token
func (t Token) String() string {
	s, _ := t.MarshalText()
	return string(s)
}

// MarshalText returns a url-safe base64-encoded token as a byte slice
func (t Token) MarshalText() ([]byte, error) {
	dst := make([]byte, codec.EncodedLen(len(t)))
	codec.Encode(dst, t)
	return dst, nil
}

// UnmarshalText decodes a url-safe base64-encoded token
func (t *Token) UnmarshalText(p []byte) error {
	n := codec.DecodedLen(len(p))
	if len(*t) < n {
		*t = append(*t, make([]byte, n-len(*t))...)
	}
	_, err := codec.Decode(*t, p)
	return err
}
