// Copyright 2019 Adrien Waksberg
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "math/rand"
  "os"
  "regexp"
  "sort"
  "strings"
)


// WalletFile contains the data in file
type WalletFile struct {
  Salt string
  Data string
}

// Wallet struct have wallet informations
type Wallet struct {
  Name       string
  Path       string
  Salt       string
  Passphrase string
  Entries    []Entry
}

// Load all wallet's Entrys from the disk
func (w *Wallet) Load() error {
  var walletFile WalletFile

  _, err := os.Stat(w.Path)
  if err != nil {
    return nil
  }

  content, err := ioutil.ReadFile(w.Path)
  if err != nil {
    return err
  }

  err = json.Unmarshal(content, &walletFile)
  if err != nil {
    return err
  }

  w.Salt = walletFile.Salt
  data, err := Decrypt(string(walletFile.Data), w.Passphrase, w.Salt)
  if err != nil {
    return err
  }

  err = json.Unmarshal(data, &w.Entries)
  if err != nil {
    return err
  }

  return nil
}

// Save the wallet on the disk
func (w *Wallet) Save() error {
  if w.Salt == "" {
    salt := make([]byte, 8)
    for i := 0; i < 8; i++ {
      salt[i] = byte(65 + rand.Intn(25))
    }
    w.Salt = string(salt)
  }

  data, err := json.Marshal(&w.Entries)
  if err != nil {
    return err
  }

  dataEncrypted, err := Encrypt(data, w.Passphrase, w.Salt)
  if err != nil {
    return err
  }

  walletFile := WalletFile{ Salt: w.Salt, Data: dataEncrypted }
  content, err := json.Marshal(&walletFile)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(w.Path, content, 0600)
  if err != nil {
    return err
  }

  return nil
}

// SearchEntry return an array with the array expected with the pattern
func (w *Wallet) SearchEntry(pattern string, group string) []Entry {
  var entries []Entry
  r := regexp.MustCompile(pattern)

  for _, entry := range w.Entries {
    if group != "" && strings.ToLower(entry.Group) != strings.ToLower(group) {
      continue
    }
    if r.Match([]byte(entry.Name)) || r.Match([]byte(entry.Comment)) || r.Match([]byte(entry.URI)) {
      entries = append(entries, entry)
    }
  }

  sort.Slice(entries, func(i, j int) bool {
    return entries[i].Group < entries[j].Group
  })

  return entries
}

// SearchEntryByID return an Entry
func (w *Wallet) SearchEntryByID(id string) Entry {
  for _, entry := range w.Entries {
    if entry.ID == id {
      return entry
    }
  }

  return Entry{}
}

// AddEntry append a new entry to wallet
func (w *Wallet) AddEntry(entry Entry) error {
  err := entry.Verify()
  if err != nil {
    return err
  }

 if w.SearchEntryByID(entry.ID) != (Entry{}) {
   return fmt.Errorf("the id already exists in wallet, can't add the entry")
 }

  w.Entries = append(w.Entries, entry)

  return nil
}

// DeleteEntry delete an entry to wallet
func (w *Wallet) DeleteEntry(id string) error {
  for index, entry := range w.Entries {
    if entry.ID == id {
      w.Entries = append(w.Entries[:index], w.Entries[index+1:]...)
      return nil
    }
  }

  return fmt.Errorf("entry not found with this id")
}

// UpdateEntry update an Entry to wallet
func (w *Wallet) UpdateEntry(entry Entry) error {
  oldEntry := w.SearchEntryByID(entry.ID)
  if oldEntry == (Entry{}) {
    return fmt.Errorf("entry not found with this id")
  }

  err := entry.Verify()
  if err != nil {
    return err
  }

  for index, i := range w.Entries {
    if entry.ID == i.ID {
      w.Entries[index] = entry
      return nil
    }
  }

  return fmt.Errorf("unknown error during the update")
}
