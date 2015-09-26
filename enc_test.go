package main

import (
	"bytes"
	"testing"
)

var (
	testMessage   = []byte("Chicago house generally refers to house music produced during the mid-1980s and late-1980s within Chicago")
	testPassword1 = []byte("Frankie Knuckles")
)

func TestEncryptCycle(t *testing.T) {
	out, err := Encrypt(testPassword1, testMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}

	out, err = Decrypt(testPassword1, out)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !bytes.Equal(testMessage, out) {
		t.Fatal("recovered plaintext doesn't match original")
	}
}
