package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorGreen  = lipgloss.Color("#00cc66")
	colorRed    = lipgloss.Color("#cc3333")
	colorCyan   = lipgloss.Color("#00cccc")
	colorDim    = lipgloss.Color("#666666")
	colorWhite  = lipgloss.Color("#eeeeee")
	colorYellow = lipgloss.Color("#cccc00")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Padding(0, 1)

	statusErrStyle = lipgloss.NewStyle().
			Foreground(colorRed).
			Padding(0, 1)

	statusOkStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Padding(0, 1)

	enabledStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	disabledStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	confirmStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true).
			Padding(0, 1)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDim).
			Padding(1, 2)
)
