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

const (
	KeySize   = 16
	SaltSize  = 32
	NonceSize = 12
)
const Overhead = SaltSize + NonceSize

func Encrypt(pass, message []byte) ([]byte, error) {
	salt, err := randBytes(SaltSize)
	if err != nil {
		return nil, err
	}

	key, err := deriveKey(pass, salt)
	if err != nil {
		return nil, err
	}

	out, err := aesEncrypt(key, message)
	if err != nil {
		return nil, err
	}

	out = append(salt, out...)
	return out, nil
}

func Decrypt(pass, message []byte) ([]byte, error) {
	if len(message) < Overhead {
		return nil, errors.New("Message is less than overhead")
	}

	key, err := deriveKey(pass, message[:SaltSize])
	if err != nil {
		return nil, err
	}

	out, err := aesDecrypt(key, message[SaltSize:])
	if err != nil {
		return nil, err
	}

	return out, nil
}

// deriveKey generates a new NaCl key from a passphrase and salt.
func deriveKey(pass, salt []byte) ([]byte, error) {
	var aesKey = make([]byte, KeySize)
	key, err := scrypt.Key(pass, salt, 1048576, 8, 1, KeySize)
	if err != nil {
		return nil, err
	}

	copy(aesKey[:], key)
	return aesKey, nil
}

// GenerateNonce creates a new random nonce.
func generateNonce() ([]byte, error) {
	nonce := make([]byte, NonceSize)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// Encrypt secures a message using AES-GCM.
func aesEncrypt(key, message []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce, err := generateNonce()
	if err != nil {
		return nil, err
	}

	out := gcm.Seal(nonce, nonce, message, nil)
	return out, nil
}

func aesDecrypt(key, message []byte) ([]byte, error) {
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

func randBytes(n int) ([]byte, error) {
	r := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
