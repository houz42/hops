package tui

import "github.com/charmbracelet/lipgloss"

// ANSI 256 colors — maps to whatever the terminal theme defines.
// This means hops looks native in Catppuccin, Dracula, Gruvbox, Solarized, etc.
var (
	colorBlue    = lipgloss.Color("4")  // enabled profiles, title accents
	colorCyan    = lipgloss.Color("6")  // selected item, input labels
	colorGray    = lipgloss.Color("8")  // disabled profiles, muted text
	colorDim     = lipgloss.Color("8")  // help bar, dim hints
	colorYellow  = lipgloss.Color("3")  // warnings, confirm prompts
	colorRed     = lipgloss.Color("1")  // errors only
	colorGreen   = lipgloss.Color("2")  // success status messages
	colorDefault = lipgloss.Color("7")  // normal text

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
			Foreground(colorBlue)

	disabledStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(colorDefault)

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
