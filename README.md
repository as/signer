# Signer

Signer is a simple token generation scheme using `xchacha20poly1305`. It generates an authenticated
and encrypted token which can be decrypted and verified by the same key. Signer is basically "Branca"
without the crufty base32 and 32-bit binary time field. It uses url-safe base64 encoding.

The Token type is a byte slice implementing base64 MarshalText and UnmarshalText.

# Use Case

You have a server that wants to give clients a token. The server need to be able to verify
that the token it issued came from the server (via the same key) and that this token was
not modified by the client or some other party. The server also wants to keep the information
inside the token private, and only accessible by the server or other parties in possession of
the key. The authentication is symmetric. Only servers in possesion of the key can verify the
token's authenticity.

You should not use Signer if the client (not just the server) needs to read and/or authenticate the token
E.G., asymmetric authentication with RSA or an elliptic curve. This token is for servers that issue
tokens, perhaps for sessions or other data.

For example: If you are using JWT with HMAC (and/or encryption), you can replace it with signer. 
If you are using JWT with RSA/EC, you can not.

# Wire Format

Below is the Token's wire format:
```
	version[1] nonce[24] ciphertext[...] tag[16] | base64

	The first 1+24 bytes are the header, authenticated by the AEAD, but not encrypted.
	The version is fixed to 0x41 (A)
	The nonce is a randomly-generated 24-byte string

	The rest is the output of the AEAD, the ciphertext and 16 byte tag.
	The ciphertext is the encrypted msg.
	The tag is a message authentication code (MAC) used to verify the integrity of the header and ciphertext

	Finally, the vertical bar denotes the Token's intended string encoding is base64 (url-safe)
```

# Interface (caller defined)
```go
type Signer interface{
	// Sign creates a token using the msg and nonce, if nonce is nil
	// one is generated automatically using a CSPRNG (crypto/rand.Read)
	Sign(msg []byte, nonce []byte) (Token, error)
	
	// Verify authenticates the token and returns the decrypted msg
	Verify(t Token) (msg []byte, err error)
}
```

# Usage Snippet
```go
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
Branca's uint32 time field isn't future-proof and ignorant of time predating 1970 (time is a signed value). Time in a token specification is scope creep. Add your own time in the msg to excersize full control, since its guaranteed to be authenitc.  

Branca implementations seen in the wild don't return msg if it is authentic but expired, which is useless for practical deployments. With Signer, you have the ability to log authentic but expired tokens to debug misbehaving clients or bugs in clients software.

Branca uses a binary time field, one advantage of having a time field would be visibility in your server logs. For example, if logs contained a unix timestamp as part of the branca header, you could easily see in Splunk whether a certain timestamp was expired. However, Branca uses a binary timestamp, so the benefit of having a greppable/searchable value is lost. You must decode the timestamp manually and in the correct big endian byte order before making assumptions about it even if you assume its authentic.

Branca uses base62, a poorly-defined standard (branca test vectors could not be decoded by online base62 decoders). A standard base64 encoding is easily accessible across languages and easier to implement and test.

Branca does not offer easily-available test vectors. Java implementations in the wild incorrectly implement AEADs that only authenticate the header and not the cipher text. We provide an interface that allows the user to pass in the nonce because in practice this is CRITICAL to reproducibly verifying test vectors in the implementation.

## Why not use JWT?

JWT is vulnerable to downgrade attacks, because it supports "none" as an encryption algorithm. Signer supports only one algorithm, so a downgrade attack is impossible by design. If chacha20poly1305 is broken in the distant future, you can use another type of token. It is not wise to rely on dynamic implementations based on token versions. Just tell the server what to expect.

JWT spec is complex and bloated. Signer is a bare-bones token that provides authentication and encryption only. It assumes the user can implement their own claims, authorization (not to be confused with authentication), and timeouts using the data inside the payload itself.

## Why not use Signer?

You want the client to be able to validate the contents of the token, using the servers public key. Signer does not support this usecase at the time of writing, and would require asymmetric public key encryption (which is an order of magnitude slower).

You don't want the contents of the token to be encrypted. It is typical to conflate encryption and authentication, but the two are not always necessary together. Signer both authenticates and encrypts the input. 

You completely distrust the programmers. This implementation allows the user to pass in a nonce, which should be random if you want identical inputs to generate distinct outputs. Many "idiot-proof" crypto packages hide the nonce-generation away from the API. However, this makes testing a pain, and sometimes the use-case calls for a deterministic nonce.

# Base64 Encoding

The Token type uses URL-safe base64-encoding without padding characters. Below is the well-known character-set that comprises base64 url-safe encoding:

`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`

# Test Vectors

This is the output of the zero vector, composed of a 32-byte key and nonce of all zero bits in binary and base64 url-safe format for your implementation in other languages:

		name:  "zero",
		input: "",
		key:   [32]byte{},
		nonce: [24]byte{},
		binary: "A\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\b\xac\x81ƕ\xb5;\xefw\n\xde5PU\xde",
		text:   "QQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAisgcaVtc2-73cK3jVQVd4",

