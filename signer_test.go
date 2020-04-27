package signer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/as/signer"
	"github.com/as/signer/branca"
)

func TestBasicSignVerify(t *testing.T) {
	s, err := signer.New(branca.Config, []byte("0123456789abcdef0123456789abcdef"), time.Second*5)
	if err != nil {
		panic(err)
	}
	const input = "hello world"
	tok, err := s.Sign([]byte(input))
	ck(t, "sign", err)

	fmt.Println(tok)

	p, err := s.Verify(tok)
	ck(t, "verify", err)

	if string(p) != input {
		t.Fatalf("have %q, want %q", string(p), input)
	}
	fmt.Printf("%s (err=%v)\n", p, err)

	p, err = s.VerifyAt((time.Time{}), tok)
	if err == nil {
		t.Fatal("verify didnt fail at zero time")
	}
	if len(p) == 0 {
		t.Fatal("verify didn't return the message, expired messages are still usefull to the caller")
	}
	fmt.Printf("%s (err=%v)\n", p, err)
}

func TestWellKnownVector(t *testing.T) {
	z := struct {
		t       time.Time
		n       [24]byte
		k       [32]byte
		input   [128]byte
		want    string
		wantBin string
	}{}
	want := `uohuCQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB4npaJ5SCNf9nh88W1NB9I7xihPkGJmK3a3ZejaTqYf46C7NXBQzv-0a9JdQwPH_KcQXSgWxGao6noMzgS4MD-pJ4e4BNKcKnUnCTgy9j8O6J-l8MyKtSH93j43GoSL6WcvjPneOouULtZCcmXHE_sL5NSP3eJLRfKpYFn3sTWx6qOWWXHu_0PFvjVMBcljBA`
	wantBin := "\xba\x88n\t\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00x\x9e\x96\x89\xe5 \x8d\u007f\xd9\xe1\xf3ŵ4\x1fH\xef\x18\xa1>A\x89\x98\xad\xdaݗ\xa3i:\x98\u007f\x8e\x82\xec\xd5\xc1C;\xfeѯIu\f\x0f\x1f\xf2\x9cAt\xa0[\x11\x9a\xa3\xa9\xe838\x12\xe0\xc0\xfe\xa4\x9e\x1e\xe0\x13Jp\xa9Ԝ$\xe0\xcb\xd8\xfc;\xa2~\x97\xc32*ԇ\xf7x\xf8\xdcj\x12/\xa5\x9c\xbe3\xe7x\xea.P\xbbY\tɗ\x1cO\xec/\x93R?w\x89-\x17ʥ\x81g\xde\xc4\xd6Ǫ\x8eYeǻ\xfd\x0f\x16\xf8\xd50\x17%\x8c\x10"
	s, err := signer.New(branca.Config, z.k[:], 0)
	if err != nil {
		panic(err)
	}
	tok := s.SignAt(z.t, z.n[:], z.input[:])
	have := string([]byte(tok))
	if have != wantBin {
		t.Fatalf("binary encoding: have %q, want %q", have, wantBin)
	}
	have = tok.String()
	if have != want {
		t.Fatalf("string encoding: have %q, want %q", have, want)
	}

	p, err := s.VerifyAt(z.t, tok)
	ck(t, "verify", err)

	if string(p) != string(z.input[:]) {
		t.Fatalf("have %q, want %q", string(p), z.input)
	}
}

func TestWellKnownBrancaVector(t *testing.T) {
	tm := time.Time{}
	// input := "Hello world!"
	s, err := signer.New(branca.Config, []byte("supersecretkeyyoushouldnotcommit"), 0)
	if err != nil {
		panic(err)
	}
	p, err := s.VerifyAt(tm, []byte(branca_hello_world))
	ck(t, "verify", err)
	t.Log(p)
}

func ck(t *testing.T, ctx string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf(ctx, err)
	}
}

// https://github.com/tuupola/branca-js/blob/9d4eee0d73d621deb763f55189ad18544870cd64/test.js#L8-L15
// supersecretkeyyoushouldnotcommit
// Hello world!
var branca_hello_world = "\xba\x00\x00\x00\x00[*\xddB_\xb6&(\x1cIZo\xa8\x83\x1f\xc9\xf0\xcf@2\x87@u\x1a\x009\x8f\x83\xef\xb3q\xd6Q@\xfa\x04\x90u\x9a\xe4pA\xf4\x11:Sd\xc7i\x1a\xf9\xb2"
