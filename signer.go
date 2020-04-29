package signer

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	Version   = 'A'
	NonceSize = chacha20poly1305.NonceSizeX // 24

	hdrSize = 1 + NonceSize
)

var (
	ErrKeyLen = errors.New("bad key length")
	ErrShort  = errors.New("message too short")
)

func New(key []byte) (*Signer, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	return &Signer{aead: aead}, nil
}

type Signer struct {
	aead cipher.AEAD

	// temporaries
	n int
	p [hdrSize]byte
}

func (s *Signer) Verify(c Token) (msg []byte, err error) {
	if len(c) < hdrSize {
		return nil, ErrShort
	}
	n := hdrSize
	ae, ad := c[n:], c[:n]
	nonce := ad[1:]
	return s.aead.Open(nil, nonce, ae, ad)
}

func (s *Signer) Sign(msg []byte, nonce []byte) (t Token, err error) {
	if nonce == nil {
		if nonce, err = mknonce(); err != nil {
			return nil, err
		}
	}
	return s.sign(msg, nonce), nil
}

func mknonce() ([]byte, error) {
	p := make([]byte, NonceSize)
	_, err := rand.Read(p)
	return p, err
}

func (s Signer) sign(msg []byte, nonce []byte) []byte {
	s.put([]byte{Version})
	s.put(nonce)
	return append(s.p[:s.n], s.aead.Seal(nil, nonce, msg, s.p[:s.n])...)
}

func (s *Signer) put(p []byte) {
	s.n += copy(s.p[s.n:], p)
}
