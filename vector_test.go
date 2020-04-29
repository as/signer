package signer_test

const testPrint = false // for generating test vectors, set to true and run go test -v -run TestWellKnownVector

var vectorTab = [...]struct {
	name  string
	input string

	key   [32]byte
	nonce [24]byte

	binary string
	text   string
}{
	{
		name:  "zero",
		input: "",
		key:   [32]byte{},
		nonce: [24]byte{},

		binary: "A\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\b\xac\x81ƕ\xb5;\xefw\n\xde5PU\xde",
		text:   "QQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAisgcaVtc2-73cK3jVQVd4",
	},
}
