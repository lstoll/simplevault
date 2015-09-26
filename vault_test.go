package main

import (
	"errors"
	"testing"
)

var (
	item1name  = "item/1"
	item1value = []byte("Chicago house generally refers to house music produced during the mid-1980s and late-1980s within Chicago")
)

func setupMockVault() Vault {
	mocks3 := &mocks3{data: map[string][]byte{}}
	return NewVault(mocks3, &mockenc{})
}

func TestMasterPassword(t *testing.T) {
	vault := &vault{}

	if vault.isMasterPassword("vvvvvSOME") {
		t.Error("Thought access password was master")
	}

	if vault.isMasterPassword("vvSOME") == false {
		t.Error("Thought master password was access")
	}
}

func TestValidateKey(t *testing.T) {
	vault := &vault{}

	if err := vault.validateKey("good-key/for_me"); err != nil {
		t.Errorf("Good key reporting error %v", err)
	}

	if err := vault.validateKey("bA*d\\Key"); err == nil {
		t.Errorf("Bad key not reporting error")
	}
}

func TestVaultCycle(t *testing.T) {
	vault := setupMockVault()

	pass, err := vault.PutItem(item1name, "masterdummy", item1value)
	if err != nil {
		t.Fatalf("Error putting item %v", err)
	}

	t.Logf("Returned password: %s", pass)

	gotItem, err := vault.GetItem(item1name, pass)
	if err != nil {
		t.Errorf("Error getting item from vault with access pass: %v", err)
	} else if string(gotItem) != string(item1value) {
		t.Errorf("With access password retrieved '%s' from vault, expected '%s'", gotItem, item1value)
	}

	gotItem, err = vault.GetItem(item1name, "masterdummy")
	if err != nil {
		t.Errorf("Error getting item from vault with master pass: %v", err)
	} else if string(gotItem) != string(item1value) {
		t.Errorf("With master password retrieved '%s' from vault, expected '%s'", gotItem, item1value)
	}

	accessPass, err := vault.GetAccessPassword(item1name, "masterpass")
	if err != nil {
		t.Errorf("Error getting access password %v", err)
	}
	if accessPass != pass {
		t.Errorf("Retrieved access pass %s, expected %s", accessPass, pass)
	}
}

// Testable Targets
type mocks3 struct {
	data map[string][]byte
}

func (m *mocks3) Get(key string) ([]byte, error) {
	val, ok := m.data[key]
	if !ok {
		//fmt.Printf("Get: %s Return: ErrorItemNotFound", key)
		return nil, errors.New("ItemNotFound")
	}
	//fmt.Printf("Get: %s Return: %s\n", key, val)
	return val, nil
}

func (m *mocks3) Put(key string, data []byte) error {
	//fmt.Printf("Put at: %s data: %s\n", key, data)
	m.data[key] = data
	return nil
}

type mockenc struct{}

func (e *mockenc) Encrypt(pass, message []byte) ([]byte, error) {
	return message, nil
}
func (e *mockenc) Decrypt(pass, message []byte) ([]byte, error) {
	return message, nil
}
