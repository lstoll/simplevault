package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mitchellh/goamz/s3"
)

const (
	AccessPasswordSize = 40
	DatbaseKey         = "internal.databass"
)

type Vault interface {
	GetItem(string, string) ([]byte, error)
	PutItem(string, string, []byte) (string, error)
	GetAccessPassword(string, string) (string, error)
}

type vault struct {
	Encryptor Enc
	S3Client  S3
}

func NewVault(s3client S3, encryptor Enc) Vault {
	return &vault{
		Encryptor: encryptor,
		S3Client:  s3client,
	}
}

func (v *vault) GetAccessPassword(key, pass string) (string, error) {
	if err := v.validateKey(key); err != nil {
		return "", err
	}

	if !v.isMasterPassword(pass) {
		return "", errors.New("A non-master password was provided")
	}

	pass, err := v.getAccessPasswordForItem(key, pass)
	if err != nil {
		return "", err
	}
	return pass, nil
}

func (v *vault) GetItem(key, pass string) ([]byte, error) {
	accessPass := pass

	if err := v.validateKey(key); err != nil {
		return nil, err
	}

	if v.isMasterPassword(pass) {
		pass, err := v.getAccessPasswordForItem(key, pass)
		accessPass = pass
		if err != nil {
			return nil, err
		}
	}

	return v.getEncryptedData(key, accessPass)
}

func (v *vault) PutItem(key, pass string, item []byte) (string, error) {
	if err := v.validateKey(key); err != nil {
		return "", err
	}

	if !v.isMasterPassword(pass) {
		return "", errors.New("Only the master pass can be used to put an item")
	}

	newPass, err := v.generateAccessPassword()
	if err != nil {
		return "", err
	}

	if err := v.putEncryptedData(key, newPass, item); err != nil {
		return "", err
	}

	if err := v.appendAccessPassword(key, pass, newPass); err != nil {
		return "", err
	}

	return newPass, err

}

func (v *vault) putEncryptedData(key, pass string, data []byte) error {
	putData, err := v.Encryptor.Encrypt([]byte(pass), data)
	if err != nil {
		return err
	}

	return v.S3Client.Put(key, putData)
}

func (v *vault) getEncryptedData(key, pass string) ([]byte, error) {
	fmt.Printf("getting data for %s\n", key)

	item, err := v.S3Client.Get(key)
	if err != nil {
		return nil, err
	}

	decrypted, err := v.Encryptor.Decrypt([]byte(pass), item)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func (v *vault) isMasterPassword(pass string) bool {
	return !strings.HasPrefix(pass, "vvvvv")
}

func (v *vault) getAccessPasswordForItem(key, pass string) (string, error) {
	data, err := v.fetchDatabase(key, pass)
	if err != nil {
		return "", err
	}

	accessPass, ok := data[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("Access password for item %s not found in database", key))
	}
	return accessPass, nil
}

func (v *vault) fetchDatabase(key, pass string) (map[string]string, error) {
	databass, err := v.getEncryptedData(DatbaseKey, pass)
	if err != nil {
		if v.isErrorS3NotFound(err) {
			// No existing database, Start an empty one
			return map[string]string{}, nil
		}
		return nil, err
	}

	data, err := databassDecode(databass)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *vault) putDatabase(key, pass string, data map[string]string) error {
	store, err := databassEncode(data)
	if err != nil {
		return err
	}

	return v.putEncryptedData(DatbaseKey, pass, store)
}

func (v *vault) appendAccessPassword(key, pass, newPass string) error {
	data, err := v.fetchDatabase(key, pass)
	if err != nil {
		return err
	}

	data[key] = newPass

	if err := v.putDatabase(key, pass, data); err != nil {
		return err
	}
	return nil
}

func (v *vault) validateKey(key string) error {
	match, err := regexp.MatchString("^[a-zA-Z0-9\\-\\_\\/]*$", key)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("Key is not valid. Keys can only contain alphanumeric, _, - and /")
	}
	return nil
}

func (v *vault) generateAccessPassword() (string, error) {
	rb := make([]byte, AccessPasswordSize-5)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}
	return "vvvvv" + base64.URLEncoding.EncodeToString(rb), nil
}

func (v *vault) isErrorS3NotFound(e error) bool {
	s3err, ok := e.(*s3.Error)
	if ok && s3err.StatusCode == 404 {
		return true
	} else if e.Error() == "ItemNotFound" {
		return true
	}
	return false
}
