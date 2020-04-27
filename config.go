package signer

import (
	"crypto/cipher"
	"crypto/rand"
)

type Config struct {
	Name     string
	Version  byte
	TimeLen  int
	NonceLen int
	KeyLen   int
	AEADFunc func(key []byte) (cipher.AEAD, error)
}

func (c *Config) Nonce() ([]byte, error) {
	p := make([]byte, c.NonceLen)
	_, err := rand.Read(p)
	return p, err
}

func (c *Config) hdrSize() int              { return 1 + c.TimeLen + c.NonceLen }
func (c *Config) nonceAt() int              { return 1 + c.TimeLen }
func (c *Config) timeAt() int               { return 1 }
func (c *Config) nonceOK(nonce []byte) bool { return len(nonce) == c.NonceLen }
