package signer

import "encoding/base64"

var (
	codec = base64.RawURLEncoding
)

type Token []byte

func (t Token) String() string {
	s, _ := t.MarshalText()
	return string(s)
}
func (t Token) MarshalText() ([]byte, error) {
	dst := make([]byte, codec.EncodedLen(len(t)))
	codec.Encode(dst, t)
	return dst, nil
}
func (t *Token) UnmarshalText(p []byte) error {
	n := codec.DecodedLen(len(p))
	if len(*t) < n {
		*t = append(*t, make([]byte, n-len(*t))...)
	}
	_, err := codec.Decode(*t, p)
	return err
}
