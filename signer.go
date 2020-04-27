package signer

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"time"
)

const (
	maxHDRLen = 64
)

var (
	ErrKeyLen  = errors.New("bad key length")
	ErrShort   = errors.New("message too short")
	ErrExpired = errors.New("message expired")
)

func New(conf Config, key []byte, ttl time.Duration) (*Signer, error) {
	if len(key) != conf.KeyLen {
		return nil, ErrKeyLen
	}
	aead, err := conf.AEADFunc(key)
	if err != nil {
		return nil, err
	}
	return &Signer{Config: &conf, aead: aead, ttl: ttl}, nil
}

type Signer struct {
	*Config
	aead cipher.AEAD
	ttl  time.Duration

	// temporaries
	n int
	p [maxHDRLen]byte
}

func (s *Signer) TTL() (time.Duration) {
	return s.ttl
}

func (s *Signer) Verify(c Token) (m []byte, err error) {
	return s.VerifyAt(time.Now(), c)
}

func (s *Signer) VerifyAt(t time.Time, c Token) (m []byte, err error) {
	if len(c) < s.hdrSize() {
		return nil, ErrShort
	}
	n := s.hdrSize()
	ae, ad := c[n:], c[:n]
	nonce := ad[s.nonceAt():]

	if m, err = s.aead.Open(nil, nonce, ae, ad); err != nil {
		return m, err
	}
	if s.ttl != 0 && s.getTime(c).After(t.Add(s.ttl)) {
		err = ErrExpired
	}
	return m, err
}

func (s *Signer) Sign(msg []byte) (Token, error) {
	nonce, err := s.Nonce()
	if err != nil {
		return nil, err
	}
	return s.signAt(time.Now(), nonce, msg), nil
}

func (s *Signer) SignAt(t time.Time, nonce []byte, msg []byte) Token {
	return s.signAt(t, nonce, msg)
}

func (s Signer) signAt(t time.Time, nonce []byte, msg []byte) []byte {
	s.put([]byte{s.Version})
	s.put32(t.Unix())
	s.put(nonce)
	return append(s.p[:s.n], s.aead.Seal(nil, nonce, msg, s.p[:s.n])...)
}

func (s *Signer) put(p []byte) {
	s.n += copy(s.p[s.n:], p)
}

func (s *Signer) put32(v int64) {
	binary.BigEndian.PutUint32(s.p[s.n:], uint32(v))
	s.n += 4
}
func (s *Signer) getTime(p []byte) time.Time {
	return time.Unix(int64(binary.BigEndian.Uint32(p[1:1+s.TimeLen])), 0)
}
