package gpm

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/atotto/clipboard"
)

// Options
var (
	LENGTH  = flag.Int("length", 16, "specify the password length")
	CONFIG  = flag.String("config", "", "specify the config file")
	WALLET  = flag.String("wallet", "", "specify the wallet")
	PASSWD  = flag.Bool("password", false, "generate and print a random password")
	DIGIT   = flag.Bool("digit", false, "use digit to generate a random password")
	LETTER  = flag.Bool("letter", false, "use letter to generate a random password")
	SPECIAL = flag.Bool("special", false, "use special chars to generate a random password")
	EXPORT  = flag.String("export", "", "json file path to export a wallet")
	IMPORT  = flag.String("import", "", "json file path to import entries")
	HELP    = flag.Bool("help", false, "print this help message")
)

// Cli struct
type Cli struct {
	Config Config
	Wallet Wallet
}

// NotificationBox print a notification
func (c *Cli) NotificationBox(msg string, error bool) {
	p := widgets.NewParagraph()
	p.SetRect(25, 20, 80, 23)
	if error {
		p.Title = "Error"
		p.Text = fmt.Sprintf("[%s](fg:red) ", msg)
	} else {
		p.Title = "Notification"
		p.Text = fmt.Sprintf("[%s](fg:green) ", msg)
	}

	ui.Render(p)
}

// ChoiceBox is a boolean form
func (c *Cli) ChoiceBox(title string, choice bool) bool {
	t := widgets.NewTabPane("Yes", "No")
	t.SetRect(10, 10, 70, 5)
	t.Title = title
	t.Border = true
	if !choice {
		t.ActiveTabIndex = 1
	}

	uiEvents := ui.PollEvents()
	for {
		ui.Render(t)
		e := <-uiEvents
		switch e.ID {
		case "<Enter>":
			return choice
		case "<Left>", "h":
			t.FocusLeft()
			choice = true
		case "<Right>", "l":
			t.FocusRight()
			choice = false
		}
	}
}

// InputBox is string form
func (c *Cli) InputBox(title string, input string, hidden bool) string {
	var secret string

	p := widgets.NewParagraph()
	p.SetRect(10, 10, 70, 5)
	p.Title = title
	p.Text = input

	uiEvents := ui.PollEvents()
	for {
		ui.Render(p)
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return ""
		case "<Backspace>":
			if len(input) >= 1 {
				input = input[:len(input)-1]
			}
		case "<Enter>":
			return input
		case "<Space>":
			input = input + " "
		default:
			input = input + e.ID
		}

		if hidden {
			secret = ""
			for i := 1; i <=  int(float64(len(input)) * 1.75); i++ {
				secret = secret + "*"
			}
			p.Text = secret
		} else {
			p.Text = input
		}
	}
}

// EntryBox to add a new entry
func (c *Cli) EntryBox(entry Entry) {
	p := widgets.NewParagraph()
	p.SetRect(25, 0, 80, 20)
	p.Text = fmt.Sprintf("%s[Name:](fg:yellow) %s\n", p.Text, entry.Name)
	p.Text = fmt.Sprintf("%s[Group:](fg:yellow) %s\n", p.Text, entry.Group)
	p.Text = fmt.Sprintf("%s[URI:](fg:yellow) %s\n", p.Text, entry.URI)
	p.Text = fmt.Sprintf("%s[User:](fg:yellow) %s\n", p.Text, entry.User)
	if entry.OTP == "" {
		p.Text = fmt.Sprintf("%s[OTP:](fg:yellow) [no](fg:red)\n", p.Text)
	} else {
		p.Text = fmt.Sprintf("%s[OTP:](fg:yellow) [yes](fg:green)\n", p.Text)
	}
	p.Text = fmt.Sprintf("%s[Comment:](fg:yellow) %v\n", p.Text, entry.Comment)

	ui.Render(p)
}

// GroupBox to select a group
func (c *Cli) GroupsBox() string {
	l := widgets.NewList()
	l.Title = "Groups"
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)
	l.WrapText = false
	l.SetRect(0, 0, 25, 23)
	l.Rows = c.Wallet.Groups()

	uiEvents := ui.PollEvents()
	for {
		ui.Render(l)
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			return ""
		case "<Enter>":
			if len(l.Rows) == 0 {
				return ""
		  } else {
				return l.Rows[l.SelectedRow]
			}
		case "j", "<Down>":
			if len(l.Rows) > 0 {
				l.ScrollDown()
			}
		case "k", "<Up>":
			if len(l.Rows) > 0 {
				l.ScrollUp()
			}
		}
	}
}

// UnlockWallet to decrypt a wallet
func (c *Cli) UnlockWallet(wallet string) error {
	var walletName string
	var err error

	ui.Clear()
	if wallet == "" {
		walletName = c.Config.WalletDefault
	} else {
		walletName = wallet
	}

	c.Wallet = Wallet{
		Name: walletName,
		Path: fmt.Sprintf("%s/%s.gpm", c.Config.WalletDir, walletName),
	}

	for i := 0; i < 3; i++ {
		c.Wallet.Passphrase = c.InputBox("Passphrase to unlock the wallet", "", true)

		err = c.Wallet.Load()
		if err == nil {
			return nil
		}
		c.NotificationBox(fmt.Sprintf("%s", err), true)
	}

	return err
}

// DeleteEntry to delete an exisiting entry
func (c *Cli) DeleteEntry(entry Entry) bool {
	if !c.ChoiceBox("Do you want delete this entry ?", false) {
		return false
	}

	err := c.Wallet.DeleteEntry(entry.ID)
	if err != nil {
		c.NotificationBox(fmt.Sprintf("%s", err), true)
		return false
	}

	err = c.Wallet.Save()
	if err != nil {
		c.NotificationBox(fmt.Sprintf("%s", err), true)
		return false
	}

	return true
}

// UpdateEntry to update an existing entry
func (c *Cli) UpdateEntry(entry Entry) bool {
	entry.Name = c.InputBox("Name", entry.Name, false)
	entry.Group = c.InputBox("Group", entry.Group, false)
	entry.URI = c.InputBox("URI", entry.URI, false)
	entry.User = c.InputBox("Username", entry.User, false)
	if c.ChoiceBox("Generate a new random password ?", false) {
		entry.Password = RandomString(c.Config.PasswordLength,
			c.Config.PasswordLetter, c.Config.PasswordDigit, c.Config.PasswordSpecial)
	}
	entry.Password = c.InputBox("Password", "", true)
	entry.OTP = c.InputBox("OTP Key", entry.OTP, false)
	entry.Comment = c.InputBox("Comment", entry.Comment, false)

	err := c.Wallet.UpdateEntry(entry)
	if err != nil {
		c.NotificationBox(fmt.Sprintf("%s", err), true)
		return false
	}

	err = c.Wallet.Save()
	if err != nil {
		c.NotificationBox(fmt.Sprintf("%s", err), true)
		return false
	}

	return true
}

// AddEntry to add new entry
func (c *Cli) AddEntry() bool {
	entry := Entry{}
	entry.GenerateID()
	entry.Name = c.InputBox("Name", "", false)
	entry.Group = c.InputBox("Group", "", false)
	entry.URI = c.InputBox("URI", "", false)
	entry.User = c.InputBox("Username", "", false)
	if c.ChoiceBox("Generate a random password ?", true) {
		entry.Password = RandomString(c.Config.PasswordLength,
			c.Config.PasswordLetter, c.Config.PasswordDigit, c.Config.PasswordSpecial)
	} else {
		entry.Password = c.InputBox("Password", "", true)
	}
	entry.OTP = c.InputBox("OTP Key", "", false)
	entry.Comment = c.InputBox("Comment", "", false)

	err := c.Wallet.AddEntry(entry)
	if err != nil {
		c.NotificationBox(fmt.Sprintf("%s", err), true)
		return false
	}

	err = c.Wallet.Save()
	if err != nil {
		c.NotificationBox(fmt.Sprintf("%s", err), true)
		return false
	}

	return true
}

// ListEntries to list all entries
func (c *Cli) ListEntries(ch chan<- bool) {
	var pattern, group string
	var entries []Entry
	var selected bool

	refresh := true
	index := -1

	l := widgets.NewList()
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)
	l.WrapText = false
	l.SetRect(0, 0, 25, 23)

	ui.Clear()
	uiEvents := ui.PollEvents()
	for {
		if group != "" {
			l.Title = fmt.Sprintf("Group: %s", group)
		} else {
			l.Title = "Group: All"
		}

		if refresh {
			refresh = false
			index = -1
			entries = c.Wallet.SearchEntry(pattern, group)
			l.Rows = []string{}
			for _, entry := range entries {
				l.Rows = append(l.Rows, entry.Name)
			}
			ui.Clear()
		}

		if len(entries) > 0 && index >= 0 && index < len(entries) {
			selected = true
		} else {
			selected = false
		}

		if selected {
			c.EntryBox(entries[index])
		}

		ui.Render(l)
		e := <-uiEvents
		switch e.ID {
		case "q":
			clipboard.WriteAll("")
			ch <- true
		case "<Enter>":
			index = l.SelectedRow
		case "<Escape>":
			pattern = ""
			group = ""
			refresh = true
		case "n":
			refresh = c.AddEntry()
		case "u":
			if selected {
				refresh = c.UpdateEntry(entries[index])
			}
		case "d":
			if selected {
				refresh = c.DeleteEntry(entries[index])
			}
		case "/":
			pattern = c.InputBox("Search", pattern, false)
			refresh = true
		case "g":
			group = c.GroupsBox()
			refresh = true
		case "j", "<Down>":
			if len(entries) > 0 {
				l.ScrollDown()
			}
		case "k", "<Up>":
			if len(entries) > 0 {
				l.ScrollUp()
			}
		case "<C-b>":
			if selected {
				clipboard.WriteAll(entries[index].User)
			}
		case "<C-c>":
			if selected {
				clipboard.WriteAll(entries[index].Password)
			}
		case "<C-o>":
			if selected {
				code, time, _ := entries[index].OTPCode()
				c.NotificationBox(fmt.Sprintf("the OTP code is available for %d seconds", time), false)
				clipboard.WriteAll(code)
			}
		}

		ch <- false
	}
}

// Import entries from json file
func (c *Cli) ImportWallet() error {
	_, err := os.Stat(*IMPORT)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(*IMPORT)
	if err != nil {
		return err
	}

	err = c.Wallet.Import(data)
	if err != nil {
		return err
	}

	err = c.Wallet.Save()
	if err != nil {
		return err
	}

	return nil
}

// Export a wallet in json format
func (c *Cli) ExportWallet() error {
	data, err := c.Wallet.Export()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(*EXPORT, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Run the cli interface
func Run() {
	var c Cli

	flag.Parse()
	c.Config.Load(*CONFIG)

	if *HELP {
		flag.PrintDefaults()
		os.Exit(1)
	} else if *PASSWD {
		fmt.Println(RandomString(*LENGTH, *LETTER, *DIGIT, *SPECIAL))
		os.Exit(0)
	}

	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v\n", err)
		os.Exit(2)
	}
	defer ui.Close()

	err := c.UnlockWallet(*WALLET)
	if err != nil {
		ui.Close()
		fmt.Printf("failed to open the wallet: %v\n", err)
		os.Exit(2)
	}

	if *IMPORT != "" {
		err := c.ImportWallet()
		if err != nil {
			ui.Close()
			fmt.Printf("failed to import: %v\n", err)
			os.Exit(2)
		}
	} else if *EXPORT != "" {
		err := c.ExportWallet()
		if err != nil {
			ui.Close()
			fmt.Printf("failed to export: %v\n", err)
			os.Exit(2)
		}
	} else {
		c1 := make(chan bool)
		go c.ListEntries(c1)

		for {
			select {
				case res := <-c1:
					if res {
						return
					}
				case <-time.After(300 * time.Second):
					return
			}
		}
	}
}
