package signer_test

import (
	"testing"

	"github.com/as/signer"
)

func TestBasicSignVerify(t *testing.T) {
	s, err := signer.New([]byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		panic(err)
	}
	const input = "hello world"
	tok, err := s.Sign([]byte(input), nil)
	ck(t, "sign", err)
	p, err := s.Verify(tok)
	ck(t, "verify", err)
	if string(p) != input {
		t.Fatalf("have %q, want %q", string(p), input)
	}
}

// TestAdversary ensures that modifying any part of the token is detectable
func TestAdversary(t *testing.T) {
	s, err := signer.New(vectorTab[0].key[:])
	if err != nil {
		panic(err)
	}
	for _, tc := range [...]struct {
		name string
		at   int
	}{
		{"version", 0},
		{"nonce", 1},
		{"ciphertext", 1 + 24},
		{"tag", 65536},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tok := []byte(vectorTab[0].binary)
			i := tc.at
			if i > len(tok) {
				i = len(tok) - 1
			}
			tok[i] = '\t'
			_, err := s.Verify(tok)
			if err == nil {
				t.Fatalf("adversary modified %s with no error", tc.name)
			}
		})
	}
}

func TestWellKnownVector(t *testing.T) {
	// see vector_test.go for test table

	for _, z := range vectorTab {
		s, err := signer.New(z.key[:])
		if err != nil {
			panic(err)
		}
		tok, err := s.Sign([]byte(z.input), z.nonce[:])
		if err != nil {
			t.Fatalf("sign: %v", err)
		}
		if testPrint {
			t.Logf("token: binary: %q", []byte(tok))
			t.Logf("token: string: %q", tok.String())
		}
		have := string([]byte(tok))
		if have != z.binary {
			t.Fatalf("binary encoding: have %q, want %q", have, z.binary)
		}
		have = tok.String()
		if have != z.text {
			t.Fatalf("string encoding: have %q, want %q", have, z.text)
		}

		p, err := s.Verify(tok)
		ck(t, "verify", err)

		if string(p) != string(z.input[:]) {
			t.Fatalf("have %q, want %q", string(p), z.input)
		}
	}
}

func ck(t *testing.T, ctx string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf(ctx, err)
	}
}
