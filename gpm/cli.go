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

package gpm

import(
  "bufio"
  "fmt"
  "math/rand"
  "os"
  "strconv"
  "syscall"
  "time"
  "github.com/atotto/clipboard"
  "github.com/olekukonko/tablewriter"
  "github.com/pquerna/otp/totp"
  "golang.org/x/crypto/ssh/terminal"
)

// Cli contain config and wallet to use
type Cli struct {
  Config Config
  Wallet Wallet
}

// printEntries show entries with tables
func (c *Cli) printEntries(entries []Entry) {
  var otp string
  var tables map[string]*tablewriter.Table

  tables = make(map[string]*tablewriter.Table)

  for i, entry := range entries {
    if entry.OTP == "" { otp = "" } else { otp = "X" }
    if _, present := tables[entry.Group]; present == false  {
      tables[entry.Group] = tablewriter.NewWriter(os.Stdout)
      tables[entry.Group].SetHeader([]string{"", "Name", "URI", "User", "OTP", "Comment"})
      tables[entry.Group].SetBorder(false)
      tables[entry.Group].SetColumnColor(
        tablewriter.Colors{tablewriter.Normal, tablewriter.FgYellowColor},
        tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},
        tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor},
        tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor},
        tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},
        tablewriter.Colors{tablewriter.Normal, tablewriter.FgMagentaColor})
    }

    tables[entry.Group].Append([]string{ strconv.Itoa(i), entry.Name, entry.URI, entry.User, otp, entry.Comment })
  }

  for group, table := range tables {
    fmt.Printf("\n%s\n\n", group)
    table.Render()
    fmt.Println("")
  }
}

// generate a random password
func (c *Cli) generatePassword(length int, letter bool, digit bool, special bool) string {
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
  chars := ""
	password := make([]byte, length)

  if letter { chars = chars + letters }
  if digit { chars = chars + digits }
  if special { chars = chars + specials }
  if !letter && !digit && !special {
	  chars = digits + letters
  }

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
	    password[i] = chars[rand.Intn(len(chars))]
	}

  return string(password)
}

// error print a message and exit)
func (c *Cli) error(msg string) {
  fmt.Println(msg)
  os.Exit(2)
}

// input from the console
func (c *Cli) input(text string, defaultValue string, show bool) string {
  fmt.Print(text)

  if show == false {
    data, _ := terminal.ReadPassword(int(syscall.Stdin))
    text := string(data)
    fmt.Printf("\n")

    if text == "" {
      return defaultValue
    }
    return text
  }

  input := bufio.NewScanner(os.Stdin)
  input.Scan()
  if input.Text() == "" {
    fmt.Printf("\n")
    return defaultValue
  }
  return input.Text()
}

// selectEntry with a form
func (c *Cli) selectEntry() Entry {
  var index int

  entries := c.Wallet.SearchEntry(*PATTERN, *GROUP)
  if len(entries) == 0 {
    fmt.Println("no entry found")
    os.Exit(1)
  }

  c.printEntries(entries)
  if len(entries) == 1 {
    return entries[0]
  }

  for true {
    index, err := strconv.Atoi(c.input("Select the entry: ", "", true))
    if err == nil && index >= 0 && index + 1 <= len(entries) {
      break
    }
    fmt.Println("your choice is not an integer or is out of range")
  }

  return entries[index]
}

// loadWallet get and unlock the wallet
func (c *Cli) loadWallet() {
  var walletName string

  passphrase := c.input("Enter the passphrase to unlock the wallet: ", "", false)

  if *WALLET == "" {
    walletName = c.Config.WalletDefault
  } else {
    walletName = *WALLET
  }

  c.Wallet = Wallet{
    Name: walletName,
    Path: fmt.Sprintf("%s/%s.gpm", c.Config.WalletDir, c.Config.WalletDefault),
    Passphrase: passphrase,
  }

  err := c.Wallet.Load()
  if err != nil {
    c.error(fmt.Sprintf("%s", err))
  }
}

// List the entry of a wallet
func (c *Cli) listEntry() {
  c.loadWallet()
  entries := c.Wallet.SearchEntry(*PATTERN, *GROUP)
  if len(entries) == 0 {
    fmt.Println("no entry found")
    os.Exit(1)
  } else {
    c.printEntries(entries)
  }
}

// Delete an entry of a wallet
func (c *Cli) deleteEntry() {
  var entry Entry

  c.loadWallet()
  entry = c.selectEntry()
  confirm := c.input("are you sure you want to remove this entry [y/N] ?", "N", true)

  if confirm == "y" {
    err := c.Wallet.DeleteEntry(entry.ID)
    if err != nil {
      c.error(fmt.Sprintf("%s", err))
    }

    err = c.Wallet.Save()
    if err != nil {
      c.error(fmt.Sprintf("%s", err))
    }

    fmt.Println("the entry has been deleted")
  }
}

// Add a new entry in wallet
func (c *Cli) addEntry() {
  c.loadWallet()

  entry := Entry{}
  entry.GenerateID()
  entry.Name     = c.input("Enter the name: ", "", true)
  entry.Group    = c.input("Enter the group: ", "", true)
  entry.URI      = c.input("Enter the URI: ",  "", true)
  entry.User     = c.input("Enter the username: ", "", true)
  if *RANDOM {
    entry.Password = c.generatePassword(c.Config.PasswordLength,
      c.Config.PasswordLetter, c.Config.PasswordDigit, c.Config.PasswordSpecial)
  } else {
    entry.Password = c.input("Enter the new password: ", entry.Password, false)
  }
  entry.OTP      = c.input("Enter the OTP key: ", "", false)
  entry.Comment  = c.input("Enter a comment: ", "", true)

  err := c.Wallet.AddEntry(entry)
  if err != nil {
    c.error(fmt.Sprintf("%s", err))
  }
  c.Wallet.Save()
}

// Update an entry in wallet
func (c *Cli) updateEntry() {
  c.loadWallet()

  entry := c.selectEntry()
  entry.Name     = c.input("Enter the new name: ", entry.Name, true)
  entry.Group    = c.input("Enter the new group: ", entry.Group, true)
  entry.URI      = c.input("Enter the new URI: ", entry.URI, true)
  entry.User     = c.input("Enter the new username: ", entry.User, true)
  if *RANDOM {
    entry.Password = c.generatePassword(c.Config.PasswordLength,
      c.Config.PasswordLetter, c.Config.PasswordDigit, c.Config.PasswordSpecial)
  } else {
    entry.Password = c.input("Enter the new password: ", entry.Password, false)
  }
  entry.OTP      = c.input("Enter the new OTP key: ", entry.OTP, false)
  entry.Comment  = c.input("Enter a new comment: ", entry.Comment, true)

  err := c.Wallet.UpdateEntry(entry)
  if err != nil {
    c.error(fmt.Sprintf("%s", err))
  }
  c.Wallet.Save()
}

// Copy login and password from an entry
func (c *Cli) copyEntry() {
  c.loadWallet()
  entry := c.selectEntry()

  for true {
    choice := c.input("select one action: ", "", true)
    switch choice {
      case "l":
        clipboard.WriteAll(entry.User)
      case "p":
        clipboard.WriteAll(entry.Password)
      case "o":
        code, _ := totp.GenerateCode(entry.OTP, time.Now())
        clipboard.WriteAll(code)
      case "q":
        os.Exit(0)
      default:
        fmt.Println("l -> copy login")
        fmt.Println("p -> copy password")
        fmt.Println("o -> copy OTP code")
        fmt.Println("q -> quit")
    }
  }
}
