package tui

import (
	"fmt"
	"strings"

	"github.com/houz42/hops/hosts"
)

type profileItem struct {
	info hosts.ProfileInfo
}

func (p profileItem) FilterValue() string { return p.info.Name }

func renderProfileList(profiles []hosts.ProfileInfo, cursor int, height int) string {
	if len(profiles) == 0 {
		return dimStyle.Render("  No profiles found. Press 'a' to add one.")
	}

	var b strings.Builder

	visible := height - 4
	if visible < 1 {
		visible = 10
	}

	start := 0
	if cursor >= visible {
		start = cursor - visible + 1
	}
	end := start + visible
	if end > len(profiles) {
		end = len(profiles)
	}

	for i := start; i < end; i++ {
		p := profiles[i]
		selected := i == cursor

		icon := "○"
		statusText := "[off]"
		style := disabledStyle
		if p.Enabled {
			icon = "●"
			statusText = "[on] "
			style = enabledStyle
		}

		line := fmt.Sprintf("  %s %s  %s  (%d entries)", icon, p.Name, statusText, p.EntryCount)

		if selected {
			b.WriteString(selectedStyle.Render("▸"+line[1:]) + "\n")
		} else {
			b.WriteString(style.Render(line) + "\n")
		}
	}

	if start > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  ↑ %d more above", start)) + "\n")
	}
	if end < len(profiles) {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  ↓ %d more below", len(profiles)-end)) + "\n")
	}

	return b.String()
}
