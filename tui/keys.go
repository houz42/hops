package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Space   key.Binding
	Back    key.Binding
	Forward key.Binding
	Add     key.Binding
	Delete  key.Binding
	Import  key.Binding
	Apply   key.Binding
	Quit    key.Binding
	Escape  key.Binding
	Tab     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "toggle"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
	Back: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "back"),
	),
	Forward: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "detail"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "x"),
		key.WithHelp("d/x", "delete"),
	),
	Import: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "import"),
	),
	Apply: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "apply to /etc/hosts"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
}

func profileListHelp() string {
	return helpStyle.Render("↑/k ↓/j navigate · enter/space toggle · →/l detail · a add · d delete · i import · S apply · q quit")
}

func detailViewHelp() string {
	return helpStyle.Render("↑/k ↓/j navigate · ←/h back · a add entry · d delete entry · q quit")
}

func confirmHelp() string {
	return confirmStyle.Render("enter confirm · esc cancel")
}

func inputHelp() string {
	return helpStyle.Render("enter confirm · esc cancel · tab next field")
}
