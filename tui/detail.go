package tui

import (
	"fmt"
	"strings"

	"github.com/houz42/hops/hosts"
)

func renderDetailView(profile string, entries []hosts.Entry, cursor int, height int) string {
	var b strings.Builder

	header := titleStyle.Render(fmt.Sprintf("Profile: %s", profile))
	b.WriteString(header + "\n\n")

	if len(entries) == 0 {
		b.WriteString(dimStyle.Render("  No entries. Press 'a' to add one.") + "\n")
		return b.String()
	}

	visible := height - 6
	if visible < 1 {
		visible = 10
	}

	start := 0
	if cursor >= visible {
		start = cursor - visible + 1
	}
	end := start + visible
	if end > len(entries) {
		end = len(entries)
	}

	maxIPLen := 0
	for _, e := range entries {
		if len(e.IP) > maxIPLen {
			maxIPLen = len(e.IP)
		}
	}

	for i := start; i < end; i++ {
		e := entries[i]
		selected := i == cursor

		line := fmt.Sprintf("  %-*s  %s", maxIPLen, e.IP, e.Hostname)

		if selected {
			b.WriteString(selectedStyle.Render("▸"+line[1:]) + "\n")
		} else {
			b.WriteString(normalStyle.Render(line) + "\n")
		}
	}

	if start > 0 {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  ↑ %d more above", start)) + "\n")
	}
	if end < len(entries) {
		b.WriteString(dimStyle.Render(fmt.Sprintf("  ↓ %d more below", len(entries)-end)) + "\n")
	}

	return b.String()
}
