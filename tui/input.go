package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

type inputMode int

const (
	inputNone inputMode = iota
	inputAddProfileName
	inputAddProfileEntries
	inputAddEntryIP
	inputAddEntryHostname
	inputImportName
	inputImportURL
)

type inputState struct {
	mode   inputMode
	fields []textinput.Model
	focus  int
	buf    string
}

func newInputState() inputState {
	return inputState{mode: inputNone}
}

func newTextInput(placeholder string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = width
	ti.CharLimit = 256
	return ti
}

func startAddProfile() inputState {
	ti := newTextInput("profile-name", 40)
	ti.Focus()
	return inputState{
		mode:   inputAddProfileName,
		fields: []textinput.Model{ti},
		focus:  0,
	}
}

func startAddProfileEntries(name string) inputState {
	ti := newTextInput("192.168.1.1 example.local (one per line, empty to finish)", 60)
	ti.Focus()
	return inputState{
		mode:   inputAddProfileEntries,
		fields: []textinput.Model{ti},
		focus:  0,
		buf:    name,
	}
}

func startAddEntry() inputState {
	ipField := newTextInput("192.168.1.1", 20)
	hostField := newTextInput("example.local", 40)
	ipField.Focus()
	return inputState{
		mode:   inputAddEntryIP,
		fields: []textinput.Model{ipField, hostField},
		focus:  0,
	}
}

func startImport() inputState {
	nameField := newTextInput("profile-name", 40)
	urlField := newTextInput("https://example.com/hosts.txt", 60)
	nameField.Focus()
	return inputState{
		mode:   inputImportName,
		fields: []textinput.Model{nameField, urlField},
		focus:  0,
	}
}

func renderInput(is inputState) string {
	var b strings.Builder

	switch is.mode {
	case inputAddProfileName:
		b.WriteString(inputLabelStyle.Render("Add Profile") + "\n\n")
		b.WriteString("  Name: " + is.fields[0].View() + "\n")

	case inputAddProfileEntries:
		b.WriteString(inputLabelStyle.Render("Add Profile: "+is.buf) + "\n\n")
		b.WriteString("  Entry (IP hostname): " + is.fields[0].View() + "\n")
		b.WriteString(dimStyle.Render("  Press enter to add, esc when done") + "\n")

	case inputAddEntryIP, inputAddEntryHostname:
		b.WriteString(inputLabelStyle.Render("Add Entry") + "\n\n")
		ipFocus := " "
		hostFocus := " "
		if is.focus == 0 {
			ipFocus = "▸"
		} else {
			hostFocus = "▸"
		}
		b.WriteString(fmt.Sprintf("  %s IP:       %s\n", ipFocus, is.fields[0].View()))
		b.WriteString(fmt.Sprintf("  %s Hostname: %s\n", hostFocus, is.fields[1].View()))

	case inputImportName, inputImportURL:
		b.WriteString(inputLabelStyle.Render("Import from URL") + "\n\n")
		nameFocus := " "
		urlFocus := " "
		if is.focus == 0 {
			nameFocus = "▸"
		} else {
			urlFocus = "▸"
		}
		b.WriteString(fmt.Sprintf("  %s Name: %s\n", nameFocus, is.fields[0].View()))
		b.WriteString(fmt.Sprintf("  %s URL:  %s\n", urlFocus, is.fields[1].View()))
	}

	return b.String()
}
