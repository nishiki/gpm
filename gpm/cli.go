package gpm

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Cli struct {
	Config Config
	Wallet Wallet
}

func (c *Cli) ErrorBox(msg string) {
	p := widgets.NewParagraph()
	p.Title = "Notification"
	p.SetRect(10, 0, 70, 5)
	p.Text = fmt.Sprintf("[ERROR: %s](fg:red) ", msg)

	ui.Render(p)
}

func (c *Cli) InputBox(msg string, hidden bool) string {
	var input, secret string

	p := widgets.NewParagraph()
	p.SetRect(10, 10, 70, 5)
	p.Text = fmt.Sprintf("%s ", msg)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	for {
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
			p.Text = fmt.Sprintf("%s %s", msg, secret)
		} else {
			p.Text = fmt.Sprintf("%s %s", msg, input)
		}
		ui.Render(p)
	}
}

func (c *Cli) EntryBox(entry Entry) {
	p := widgets.NewParagraph()
	p.Title = "Entry"
	p.SetRect(25, 0, 80, 23)
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
		c.Wallet.Passphrase = c.InputBox("Enter the passphrase to unlock the wallet:\n", true)

		err = c.Wallet.Load()
		if err == nil {
			return nil
		}
		c.ErrorBox(fmt.Sprintf("%s", err))
	}

	return err
}

func (c *Cli) GroupsBox() string {
	return ""
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
		case "/":
			pattern = c.InputBox("Search", false)
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
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	err := c.UnlockWallet("test")
	if err != nil {
		return
	}

	c.ListEntries()
}
