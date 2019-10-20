package gpm

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func generateWalletWithEntries() Wallet {
	var wallet Wallet

	for i := 0; i < 10; i++ {
		entry := Entry{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Entry %d", i), Group: "Good Group"}
		wallet.AddEntry(entry)
	}

	return wallet
}

func TestAddBadEntry(t *testing.T) {
	var entry Entry
	var wallet Wallet

	err := wallet.AddEntry(entry)
	if err == nil {
		t.Error("a bad entry must return an error")
	}
}

func TestAddEntries(t *testing.T) {
	var wallet Wallet

	for i := 0; i < 10; i++ {
		entry := Entry{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Entry %d", i)}
		err := wallet.AddEntry(entry)
		if err != nil {
			t.Errorf("a good entry mustn't return an error: %s", err)
		}
	}

	if len(wallet.Entries) != 10 {
		t.Errorf("must have 10 entries: %d", len(wallet.Entries))
	}
}

func TestSearchEntryWithBadID(t *testing.T) {
	wallet := generateWalletWithEntries()
	entry := wallet.SearchEntryByID("BAD-ID")
	if entry.ID != "" {
		t.Errorf("if the entry doesn't exist must return an empty Entry: %s", entry.ID)
	}
}

func TestSearchEntryWithGoodID(t *testing.T) {
	wallet := generateWalletWithEntries()
	entry := wallet.SearchEntryByID("5")
	if entry.ID != "5" {
		t.Errorf("the ID entry must be 5: %s", entry.ID)
	}
}

func TestSearchEntriesByGroup(t *testing.T) {
	wallet := generateWalletWithEntries()
	entries := len(wallet.SearchEntry("", "BAD-GROUP", false))
	if entries != 0 {
		t.Errorf("a search with bad group must return 0 entry: %d", entries)
	}

	entries = len(wallet.SearchEntry("", "good group", false))
	if entries != 10 {
		t.Errorf("a search with good group must return 10 entries: %d", entries)
	}
}

func TestSearchEntriesByPattern(t *testing.T) {
	wallet := generateWalletWithEntries()
	entries := len(wallet.SearchEntry("BAD-PATTERN", "", false))
	if entries != 0 {
		t.Errorf("a search with bad pattern must return 0 entry: %d", entries)
	}

	entries = len(wallet.SearchEntry("entry", "", false))
	if entries != 10 {
		t.Errorf("a search with good pattern must return 10 entries: %d", entries)
	}

	entries = len(wallet.SearchEntry("^entry 5$", "", false))
	if entries != 1 {
		t.Errorf("a search with specific pattern must return 1 entry: %d", entries)
	}
}

func TestSearchEntriesByPatternAndGroup(t *testing.T) {
	wallet := generateWalletWithEntries()
	entries := len(wallet.SearchEntry("entry", "good group", false))
	if entries != 10 {
		t.Errorf("a search with good pattern and godd group must return 10 entries: %d", entries)
	}
}

func TestDeleteNotExistingEntry(t *testing.T) {
	wallet := generateWalletWithEntries()
	err := wallet.DeleteEntry("BAD-ID")
	if err == nil {
		t.Error("if the entry doesn't exist must return an error")
	}

	if len(wallet.Entries) != 10 {
		t.Errorf("must have 10 entries: %d", len(wallet.Entries))
	}
}

func TestDeleteEntry(t *testing.T) {
	wallet := generateWalletWithEntries()
	err := wallet.DeleteEntry("5")
	if err != nil {
		t.Errorf("a good entry mustn't return an error: %s", err)
	}

	if len(wallet.Entries) != 9 {
		t.Errorf("must have 9 entries: %d", len(wallet.Entries))
	}

	if wallet.SearchEntryByID("5").ID != "" {
		t.Error("must return an empty entry for the ID 5")
	}
}

func TestUpdateNotExistingEntry(t *testing.T) {
	wallet := generateWalletWithEntries()
	err := wallet.UpdateEntry(Entry{ID: "BAD-ID"})
	if err == nil {
		t.Error("if the entry doesn't exist must return an error")
	}
}

func TestUpdateEntry(t *testing.T) {
	wallet := generateWalletWithEntries()
	err := wallet.UpdateEntry(Entry{ID: "5"})
	if err == nil {
		t.Error("if the entry is bad must return an error")
	}

	err = wallet.UpdateEntry(Entry{ID: "5", Name: "Name 5"})
	if err != nil {
		t.Errorf("a good entry mustn't return an error: %s", err)
	}

	entry := wallet.SearchEntryByID("5")
	if entry.Name != "Name 5" {
		t.Errorf("the entry name for the ID 5 must be 'Name 5': %s", entry.Name)
	}
}

func TestExportAndImport(t *testing.T) {
	wallet := generateWalletWithEntries()
	export, err := wallet.Export()
	if err != nil {
		t.Errorf("an export mustn't return an error: %s", err)
	}

	wallet = Wallet{}
	err = wallet.Import(export)
	if err != nil {
		t.Errorf("a good import mustn't return an error: %s", err)
	}

	entries := len(wallet.Entries)
	if entries != 10 {
		t.Errorf("must have 10 entries: %d", entries)
	}
}

func TestSaveWallet(t *testing.T) {
	tmpFile, _ := ioutil.TempFile(os.TempDir(), "gpm_test-")
	defer os.Remove(tmpFile.Name())

	wallet := generateWalletWithEntries()
	wallet.Path = tmpFile.Name()
	wallet.Passphrase = "secret"

	err := wallet.Save()
	if err != nil {
		t.Errorf("save wallet mustn't return an error: %s", err)
	}
}

func TestLoadWalletWithGoodPassword(t *testing.T) {
	var loadWallet Wallet

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "gpm_test-")
	defer os.Remove(tmpFile.Name())

	wallet := generateWalletWithEntries()
	wallet.Path = tmpFile.Name()
	wallet.Passphrase = "secret"
	wallet.Save()
	loadWallet.Path = wallet.Path
	loadWallet.Passphrase = wallet.Passphrase

	err := loadWallet.Load()
	if err != nil {
		t.Errorf("load wallet mustn't return an error: %s", err)
	}

	entries := len(loadWallet.Entries)
	if entries != 10 {
		t.Errorf("must have 10 entries: %d", entries)
	}
}

func TestLoadWalletWithBadPassword(t *testing.T) {
	var loadWallet Wallet

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "gpm_test-")
	defer os.Remove(tmpFile.Name())

	wallet := generateWalletWithEntries()
	wallet.Path = tmpFile.Name()
	wallet.Passphrase = "secret"
	wallet.Save()
	loadWallet.Path = wallet.Path
	loadWallet.Passphrase = "bad secret"

	err := loadWallet.Load()
	if err == nil {
		t.Error("load wallet with bad password must return an error")
	}

	entries := len(loadWallet.Entries)
	if entries != 0 {
		t.Errorf("must have 0 entries: %d", entries)
	}
}

func TestGetGroup(t *testing.T) {
	wallet := generateWalletWithEntries()
	groups := wallet.Groups()
	if len(groups) != 1 {
		t.Errorf("there must have 1 group: %d", len(groups))
	}
	if groups[0] != "Good Group" {
		t.Errorf("the group name isn't 'Good Group': %s", groups[0])
	}
}
