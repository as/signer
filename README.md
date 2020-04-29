# Signer

Signer is a simple token scheme based on xchacha20poly1305. It has a binary format
similar to "branca", except it has no 32-bit binary time field, and uses base64
url-safe encoding (instead of base62).

Token is just a byte slice that implements MarshalText and UnmarshalText in base64.

# Wire Format
```
	version[1] nonce[24] ciphertext[...] tag[16] | base64

	The first 24 bytes are the header, authenticated by the AEAD, but not encrypted.
	The version is fixed to 0x41 (A)
	The nonce is a randomly-generate 24-byte string
	The rest is the output of the AEAD, the ciphertext and 16 byte tag. 
```

# Interface (callee defined)
```
type Signer interface{
	// Sign creates a token using the msg and nonce, if nonce is nil
	// one is generated automatically.
	Sign(msg []byte, nonce []byte) (Token, error)
	
	// Verify authenticates the token and returns the decrypted msg
	Verify(t Token) (msg []byte, err error)
}
```

# Usage Snippet
```
	// Configure
	key := [32]byte{ /* random data */ }
	s, _ := signer.New(key[:])

	// Sign
	tok, _ := s.Sign([]byte("hello world"), nil)
	fmt.Println(tok)
	// ul6mbjrzW_Y82_a8sQQRqlzFTPAcA65tn4xlWN3z3bpwIYZiW47JlyF34UwaUzize4yFfrN8Vzs

	// Verify
	p, err := s.Verify(tok)
	if err != nil{
		log.Fatalf("verify: %v", err)
	}
```

