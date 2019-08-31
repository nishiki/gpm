package gpm

import (
	"regexp"
	"testing"
)

func TestEncrypt(t *testing.T) {
	secret := []byte("secret data")

	data, err := Encrypt(secret, "passphrase", "salt")
	if err != nil {
		t.Errorf("Encrypt mustn't return an error: %s", err)

	}
	if data == "" {
		t.Error("Encrypt must generate a string not empty")
	}
}

func TestDecrypt(t *testing.T) {
	secret := "secret data"

	dataEncrypted, _ := Encrypt([]byte(secret), "passphrase", "salt")
	data, err := Decrypt(dataEncrypted, "passphrase", "salt")
	if err != nil {
		t.Errorf("Decrypt mustn't return an error: %s", err)
	}
	if string(data) != secret {
		t.Errorf("the encrypted secret is different of decrypted secret: %s", data)
	}
}

func TestDecryptWithBadPassphrase(t *testing.T) {
	secret := []byte("secret data")

	dataEncrypted, _ := Encrypt(secret, "passphrase", "salt")
	_, err := Decrypt(dataEncrypted, "bad", "salt")
	if err == nil {
		t.Error("Decrypt must return an error with bad passphrase")
	}
}

func TestDecryptWithBadSalt(t *testing.T) {
	secret := []byte("secret data")

	dataEncrypted, _ := Encrypt(secret, "passphrase", "salt")
	_, err := Decrypt(dataEncrypted, "passphrase", "bad")
	if err == nil {
		t.Error("Decrypt must return an error with bad salt")
	}
}

func TestRandomStringLength(t *testing.T) {
	password := RandomString(64, false, false, false)
	if len(password) != 64 {
		t.Errorf("the string must have 64 chars: %d", len(password))
	}
	r := regexp.MustCompile(`^[a-zA-Z0-9]{64}$`)
	match := r.FindSubmatch([]byte(password))
	if len(match) == 0 {
		t.Errorf("the string must contain only digit and alphabetic characters: %s", password)
	}
}

func TestRandomStringOnlyDigit(t *testing.T) {
	password := RandomString(64, false, true, false)
	r := regexp.MustCompile(`^[0-9]{64}$`)
	match := r.FindSubmatch([]byte(password))
	if len(match) == 0 {
		t.Errorf("the string must contain only digit characters: %s", password)
	}
}

func TestRandomStringOnlyAlphabetic(t *testing.T) {
	password := RandomString(64, true, false, false)
	r := regexp.MustCompile(`^[a-zA-Z]{64}$`)
	match := r.FindSubmatch([]byte(password))
	if len(match) == 0 {
		t.Errorf("the string must contain only alphabetic characters: %s", password)
	}
}

func TestRandomStringOnlySpecial(t *testing.T) {
	password := RandomString(64, false, false, true)
	r := regexp.MustCompile(`^[\~\=\+\%\^\*\/\(\)\[\]\{\}\!\@\#\$\?\|]{64}$`)
	match := r.FindSubmatch([]byte(password))
	if len(match) == 0 {
		t.Errorf("the string must contain only alphabetic characters: %s", password)
	}
}
