package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

type Enc interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Decrypt([]byte, []byte) ([]byte, error)
}

type enc struct{}

const (
	KeySize   = 16
	SaltSize  = 32
	NonceSize = 12
)
const Overhead = SaltSize + NonceSize

func NewEnc() Enc {
	return &enc{}
}

func (e *enc) Encrypt(pass, message []byte) ([]byte, error) {
	salt, err := e.randBytes(SaltSize)
	if err != nil {
		return nil, err
	}

	key, err := e.deriveKey(pass, salt)
	if err != nil {
		return nil, err
	}

	out, err := e.aesEncrypt(key, message)
	if err != nil {
		return nil, err
	}

	out = append(salt, out...)
	return out, nil
}

func (e *enc) Decrypt(pass, message []byte) ([]byte, error) {
	if len(message) < Overhead {
		return nil, errors.New("Message is less than overhead")
	}

	key, err := e.deriveKey(pass, message[:SaltSize])
	if err != nil {
		return nil, err
	}

	out, err := e.aesDecrypt(key, message[SaltSize:])
	if err != nil {
		return nil, err
	}

	return out, nil
}

// deriveKey generates a new NaCl key from a passphrase and salt.
func (e *enc) deriveKey(pass, salt []byte) ([]byte, error) {
	var aesKey = make([]byte, KeySize)
	key, err := scrypt.Key(pass, salt, 1048576, 8, 1, KeySize)
	if err != nil {
		return nil, err
	}

	copy(aesKey[:], key)
	return aesKey, nil
}

// GenerateNonce creates a new random nonce.
func (e *enc) generateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// Encrypt secures a message using AES-GCM.
func (e *enc) aesEncrypt(key, message []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce, err := e.generateNonce()
	if err != nil {
		return nil, err
	}

	out := gcm.Seal(nonce, nonce, message, nil)
	return out, nil
}

func (e *enc) aesDecrypt(key, message []byte) ([]byte, error) {
	if len(message) <= NonceSize {
		return nil, errors.New("Message is too short")
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, message)

	out, err := gcm.Open(nil, nonce, message[NonceSize:], nil)
	if err != nil {
		fmt.Println("this err")
		return nil, err
	}
	return out, nil
}

func (e *enc) randBytes(n int) ([]byte, error) {
	r := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
