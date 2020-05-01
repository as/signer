package signer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/as/signer"
)

// run with 'go test -race .'
func TestConcurrentSignVerify(t *testing.T) {
	const (
		N = 64
		T = 3* time.Second 
	)
	s, err := signer.New([]byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	done, errc := make(chan bool), make(chan error)
	defer close(done)

	for i := 0; i < N; i++ {
		go signVerify(s, done, errc)
	}

	select {
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(T):
	}
}

func signVerify(s *signer.Signer, done chan bool, errc chan error) {
	const input = "hello world"
	for {
		select {
		case <-done:
			return
		default:
		}
		tok, err := s.Sign([]byte(input), nil)
		if err != nil {
			errc <- err
			return
		}
		p, err := s.Verify(tok)
		if err != nil {
			errc <- err
			return
		}
		if string(p) != input {
			if err != nil {
				errc <- fmt.Errorf("have %q, want %q", string(p), input)
				return
			}
		}
	}
}
