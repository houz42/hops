package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/houz42/hops/hosts"
	"github.com/houz42/hops/tui"
)

func main() {
	hostsFile := flag.String("f", "/etc/hosts", "path to hosts file")
	applyWrite := flag.String("apply-write", "", "internal: privileged write target")
	blockFile := flag.String("block-file", "", "internal: temp file with block content")
	flag.Parse()

	if *applyWrite != "" && *blockFile != "" {
		if err := hosts.ApplyWrite(*applyWrite, *blockFile); err != nil {
			fmt.Fprintf(os.Stderr, "apply-write failed: %s\n", err)
			os.Exit(1)
		}
		return
	}

	manager, err := hosts.NewManager(*hostsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer manager.Close()

	model := tui.NewModel(manager)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %s\n", err)
		os.Exit(1)
	}
}
