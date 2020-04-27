# Signer

This is just a simple token scheme based on chacha20poly1305. It has a binary format
similar to "baranca". 

Only difference is it uses base64 encoding instead of base62* for the Token type. The
Token is just a byte slice that implements MarshalText and UnmarshalText in base64.

# Branca Token 
```
	version[1] time[4] nonce[24] | ciphertext[...] tag[16]

	The first 29 bytes are the header, authenticated by the AEAD, but not encrypted
	The rest is the output of the AEAD, the ciphertext and 16 byte tag. Simple.

	WARNING: "Branca" packages on github for some languages do not authenticate
	the ciphertext, and are incompatible with this format.
```

# Interface (callee defined)
```
type Signer interface{
	VerifyAt(t time.Time, c Token) (m []byte, err error)
	SignAt(t time.Time, nonce []byte, msg []byte) Token
	TTL() time.Duration
}
```

# Usage Snippet
```
	// Configure
	key := [32]byte{ /* random data */ }
	ttl := 5*time.Second
	s, _ := signer.New(branca.Config, key[:], ttl)

	// Sign
	tok, _ := s.Sign([]byte("hello world"))
	fmt.Println(tok)
	// ul6mbjrzW_Y82_a8sQQRqlzFTPAcA65tn4xlWN3z3bpwIYZiW47JlyF34UwaUzize4yFfrN8Vzs

	// Verify
	p, err := s.Verify(tok)
	if err != nil{
		if err == signer.ErrExpired{
			// expired messages are still useful to the caller
			log.Printf("verify: timed out: token: %q", p)
		} else {
			log.Printf("verify: %v", err)
		}
	}
```

# Notes

The choice to include a time field was the decision of the branca people. I consider this
scope creep (like JWT). You can configure a TTL of 0 in the constructor to disable this
check. If you know the message is valid, you can store your own time value (even multiple
create, modified, updated time). Time field assumes too much about your requirements.

The time field is also 32-bits wide, and unsigned (did you know there is a year before 1970?).
We should probably just remove the time stuff from this package. The more I write this out
the more I hate the idea, especially since it uses a fixed binary header.

Here's another reason I hate the time field. In most implementations the message is expired
and all you get is an error. The message's contents might still have important information. This
package returns the message on ErrExpired

Base62 encoding is slow and has no standard support in the standard libraries of many
languages.

# The current token
version[1] time[4] nonce[24] | ciphertext[...] tag[16]

# The ideal token
version[1] nonce[24] ciphertext[...] tag[16] | base64
