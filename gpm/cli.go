package gpm

import (
	"fmt"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Cli struct {
	Config Config
	Wallet Wallet
}

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

func (c *Cli) EntryBox(entry Entry) {
	p := widgets.NewParagraph()
	p.Title = "Entry"
	p.SetRect(25, 0, 80, 20)
	p.Text = fmt.Sprintf("%s[Name:](fg:yellow) %s\n", p.Text, entry.Name)
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
			return l.Rows[l.SelectedRow]
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

func (c *Cli) UpdateEntry(entry Entry) error {
	entry.Name = c.InputBox("Name", entry.Name, false)
	entry.Group = c.InputBox("Group", entry.Group, false)
	entry.URI = c.InputBox("URI", entry.URI, false)
	entry.User = c.InputBox("Username", entry.User, false)
	entry.Password = c.InputBox("Password", "", true)
	entry.OTP = c.InputBox("OTP Key", entry.OTP, false)
	entry.Comment = c.InputBox("Comment", entry.Comment, false)

	err := c.Wallet.UpdateEntry(entry)
	if err != nil {
		return err
	}

	err = c.Wallet.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cli) AddEntry() error {
	entry := Entry{}
	entry.GenerateID()
	entry.Name = c.InputBox("Name", "", false)
	entry.Group = c.InputBox("Group", "", false)
	entry.URI = c.InputBox("URI", "", false)
	entry.User = c.InputBox("Username", "", false)
	entry.Password = c.InputBox("Password", "", true)
	entry.OTP = c.InputBox("OTP Key", "", false)
	entry.Comment = c.InputBox("Comment", "", false)

	err := c.Wallet.AddEntry(entry)
	if err != nil {
		return err
	}

	err = c.Wallet.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cli) ListEntries() {
	var pattern, group string
	var entries []Entry

	refresh := true
	index := -1

	l := widgets.NewList()
	l.Title = "Entries"
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen, ui.ColorClear, ui.ModifierBold)
	l.WrapText = false
	l.SetRect(0, 0, 25, 23)

	ui.Clear()
	uiEvents := ui.PollEvents()
	for {
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

		if index >= 0 {
			c.EntryBox(entries[index])
		}

		ui.Render(l)
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Enter>":
			index = l.SelectedRow
		case "<Escape>":
			pattern = ""
			refresh = true
		case "n":
			err := c.AddEntry()
			if err == nil {
				refresh = true
			} else {
				c.NotificationBox(fmt.Sprintf("%s", err), true)
			}
		case "u":
			if len(entries) > 0 && index >= 0 && index < len(entries) {
				err := c.UpdateEntry(entries[index])
				if err == nil {
					refresh = true
				} else {
					c.NotificationBox(fmt.Sprintf("%s", err), true)
				}
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
		}
	}
}

func Run() {
	var c Cli
	c.Config.Load("")

	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v\n", err)
		os.Exit(2)
	}
	defer ui.Close()

	err := c.UnlockWallet("test")
	if err != nil {
		ui.Close()
		fmt.Printf("failed to open the wallet: %v\n", err)
		os.Exit(2)
	}

	c.ListEntries()
}
