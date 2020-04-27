package branca

import (

	"github.com/as/signer"
	"golang.org/x/crypto/chacha20poly1305"
)

var Config = signer.Config{
	Name:     "branca.xchacha20poly1305",
	Version:  0xBA,
	TimeLen: 4,	// 32-bit timestamp? why?
	NonceLen: 24,
	KeyLen: 32,
	AEADFunc: chacha20poly1305.NewX,
}

/*
	// better idea; for later
type Header struct {
	V byte
	T uint32
	N [24]byte
}

func (h Header) New(nonce []byte, t time.Time) Header {
	h.T = uint32(t.Unix())
	copy(h.N[:], nonce)
	return h
}
func (h *Header) Version() int { return int(h.V) }
func (h *Header) Nonce() []byte     { return h.N[:] }
func (h *Header) Len() int     { return 1 + 4 + 24 }
func (h *Header) Time() time.Time {
	return time.Unix(int64(h.T), 0)
}
*/