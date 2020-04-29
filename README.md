# Signer

Signer is a simple token scheme based on xchacha20poly1305. It is used to generate
authenticated *and encrypted* tokens which can be issued to clients. It has a binary
format similar to "branca", except it has no 32-bit binary time field, and uses base64
url-safe encoding (instead of base62).

Token is just a byte slice that implements MarshalText and UnmarshalText in base64.

# Use Case

You have a server that wants to give clients a token. The server need to be able to verify
that the token it issued came from the server (via the same key) and that this token was
not modified by the client or some other party. The server also wants to keep the information
inside the token private, and only accessible by the server or other parties in possession of
the key. Only the server can verify the authenticity of the key (the authentication is symmetric)

You should not use this if you want the client to be able to read and authenticate the data stored
in the token. E.G., asymmetric authentication with RSA or an elliptic curve. This token is for
servers that issue tokens, perhaps for sessions or other data.

# Wire Format
```
	version[1] nonce[24] ciphertext[...] tag[16] | base64

	The first 1+24 bytes are the header, authenticated by the AEAD, but not encrypted.
	The version is fixed to 0x41 (A)
	The nonce is a randomly-generated 24-byte string
	The rest is the output of the AEAD, the ciphertext and 16 byte tag. 
```

# Interface (caller defined)
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

# Notes

## Why not use branca?
Branca has a 32-bit unsigned binary time field and isn't future-proof. Having time in the specification is a scope creep. You can add your own time field in the message and have full control over how to use it. Most branca implementations don't return the message if the message is authentic, but expired. This is rather useless in practical deployments! We want the ability to log authentic but expired tokens to debug misbehaving clients or bugs in clients software.

Branca uses base62. There are many opinions of what base62 actually is (branca test vectors could not be decoded by online base62 decoders.). We prefer a standard encoding in a binary power of 2 that is easily accessible across languages).

## Why not use JWT?

JWT is vulnerable to downgrade attacks, because it supports "none" as an encryption algorithm. Signer supports only one algorithm, so a downgrade attack is impossible by design. If chacha20poly1305 is broken in the distant future, you can use another type of token. It is not wise to rely on dynamic implementations based on token versions. Just tell the server what to expect.

JWT spec is complex and bloated. Signer is a bare-bones token that provides authentication and encryption only. It assumes the user can implement their own claims, authorization (not to be confused with authentication), and timeouts using the data inside the payload itself.

## Why not use Signer?

You want the client to be able to validate the contents of the token, using the servers public key. Signer does not support this usecase at the time of writing. However, it may be feasible to support this in another version if there is pressing need.

