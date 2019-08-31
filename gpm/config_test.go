package gpm

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	var config Config

	err := config.Init()
	if err != nil {
		t.Error("the config init mustn't return an error")
	}

	if config.WalletDefault != "default" {
		t.Errorf("the WalletDefaut must be 'default': %s", config.WalletDefault)
	}

	if config.PasswordLength != 16 {
		t.Errorf("the PasswordLength must be 16: %d", config.PasswordLength)
	}

	if config.PasswordLetter != true {
		t.Error("the PasswordLetter must be true")
	}

	if config.PasswordDigit != true {
		t.Error("the PasswordDigit must be true")
	}

	if config.PasswordSpecial != false {
		t.Error("the PasswordSpecial must be false")
	}
}

func TestSave(t *testing.T) {
	var config Config

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "gpm_test-")
	defer os.Remove(tmpFile.Name())

	config.Init()
	err := config.Save(tmpFile.Name())
	if err != nil {
		t.Errorf("save config mustn't return an error: %s", err)
	}
}

func TestLoadWithFile(t *testing.T) {
	var config Config

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "gpm_test-")
	defer os.Remove(tmpFile.Name())

	config.Init()
	config.Save(tmpFile.Name())
	err := config.Load(tmpFile.Name())
	if err != nil {
		t.Errorf("load config with file mustn't return an error: %s", err)
	}
}

func TestLoadWithoutFile(t *testing.T) {
	var config Config

	err := config.Load("")
	if err != nil {
		t.Errorf("load config without file mustn't return an error: %s", err)
	}
}
