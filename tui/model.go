package tui

import (
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/houz42/hops/hosts"
)

type viewState int

const (
	viewProfiles viewState = iota
	viewDetail
	viewInput
	viewConfirm
)

type Model struct {
	manager *hosts.Manager

	view     viewState
	width    int
	height   int
	prevView viewState

	profiles      []hosts.ProfileInfo
	profileCursor int

	detailProfile string
	entries       []hosts.Entry
	entryCursor   int

	input inputState

	pendingEntries []hosts.Entry

	confirm       confirmAction
	confirmTarget string

	dirty bool

	statusText string
	statusErr  bool
}

func NewModel(manager *hosts.Manager) Model {
	profiles := manager.ListProfiles()
	return Model{
		manager:  manager,
		view:     viewProfiles,
		profiles: profiles,
		input:    newInputState(),
		dirty:    manager.HasDirtyState(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case errMsg:
		m.statusText = msg.Error()
		m.statusErr = true
		return m, nil

	case statusMsg:
		m.statusText = string(msg)
		m.statusErr = false
		return m, nil

	case applyDoneMsg:
		m.view = viewProfiles
		if msg.err != nil {
			m.statusText = fmt.Sprintf("Apply failed: %s", msg.err)
			m.statusErr = true
		} else {
			m.dirty = m.manager.HasDirtyState()
			m.statusText = "Applied to /etc/hosts"
			m.statusErr = false
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, keys.Quit) && m.view != viewInput {
		return m, tea.Quit
	}

	switch m.view {
	case viewProfiles:
		return m.handleProfileKeys(msg)
	case viewDetail:
		return m.handleDetailKeys(msg)
	case viewInput:
		return m.handleInputKeys(msg)
	case viewConfirm:
		return m.handleConfirmKeys(msg)
	}

	return m, nil
}

func (m Model) handleProfileKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.profileCursor > 0 {
			m.profileCursor--
		}
	case key.Matches(msg, keys.Down):
		if m.profileCursor < len(m.profiles)-1 {
			m.profileCursor++
		}
	case key.Matches(msg, keys.Enter) || key.Matches(msg, keys.Space):
		if len(m.profiles) > 0 {
			p := m.profiles[m.profileCursor]
			if err := m.manager.Toggle(p.Name); err != nil {
				m.statusText = fmt.Sprintf("Toggle failed: %s", err)
				m.statusErr = true
			} else {
				m.profiles = m.manager.ListProfiles()
				m.dirty = m.manager.HasDirtyState()
				action := "Enabled"
				if p.Enabled {
					action = "Disabled"
				}
				m.statusText = fmt.Sprintf("%s %s", action, p.Name)
				m.statusErr = false
			}
		}
	case key.Matches(msg, keys.Forward):
		if len(m.profiles) > 0 {
			p := m.profiles[m.profileCursor]
			m.detailProfile = p.Name
			entries, err := m.manager.GetEntries(p.Name)
			if err != nil {
				m.statusText = fmt.Sprintf("Failed to load entries: %s", err)
				m.statusErr = true
			} else {
				m.entries = entries
				m.entryCursor = 0
				m.view = viewDetail
			}
		}
	case key.Matches(msg, keys.Add):
		m.input = startAddProfile()
		m.prevView = viewProfiles
		m.view = viewInput
		m.pendingEntries = nil
	case key.Matches(msg, keys.Delete):
		if len(m.profiles) > 0 {
			p := m.profiles[m.profileCursor]
			m.confirm = confirmRemoveProfile
			m.confirmTarget = p.Name
			m.view = viewConfirm
		}
	case key.Matches(msg, keys.Import):
		m.input = startImport()
		m.prevView = viewProfiles
		m.view = viewInput
	case key.Matches(msg, keys.Apply):
		if m.dirty {
			m.confirm = confirmApply
			m.confirmTarget = ""
			m.view = viewConfirm
		} else {
			m.statusText = "Nothing to apply — /etc/hosts is up to date"
			m.statusErr = false
		}
	}

	return m, nil
}

func (m Model) handleDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if m.entryCursor > 0 {
			m.entryCursor--
		}
	case key.Matches(msg, keys.Down):
		if m.entryCursor < len(m.entries)-1 {
			m.entryCursor++
		}
	case key.Matches(msg, keys.Back):
		m.profiles = m.manager.ListProfiles()
		m.view = viewProfiles
	case key.Matches(msg, keys.Add):
		m.input = startAddEntry()
		m.prevView = viewDetail
		m.view = viewInput
	case key.Matches(msg, keys.Delete):
		if len(m.entries) > 0 {
			e := m.entries[m.entryCursor]
			m.confirm = confirmRemoveEntry
			m.confirmTarget = e.Hostname
			m.view = viewConfirm
		}
	case key.Matches(msg, keys.Escape):
		m.profiles = m.manager.ListProfiles()
		m.view = viewProfiles
	}

	return m, nil
}

func (m Model) handleInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Escape):
		if m.input.mode == inputAddProfileEntries && len(m.pendingEntries) > 0 {
			if err := m.manager.AddProfile(m.input.buf, m.pendingEntries); err != nil {
				m.statusText = fmt.Sprintf("Add profile failed: %s", err)
				m.statusErr = true
			} else {
				m.profiles = m.manager.ListProfiles()
				m.dirty = m.manager.HasDirtyState()
				m.statusText = fmt.Sprintf("Added profile %q with %d entries", m.input.buf, len(m.pendingEntries))
				m.statusErr = false
			}
			m.pendingEntries = nil
		}
		m.input = newInputState()
		m.view = m.prevView
		return m, nil

	case key.Matches(msg, keys.Tab):
		return m.handleInputTab(), nil

	case key.Matches(msg, keys.Enter):
		return m.handleInputEnter()
	}

	if len(m.input.fields) > 0 && m.input.focus < len(m.input.fields) {
		var cmd tea.Cmd
		m.input.fields[m.input.focus], cmd = m.input.fields[m.input.focus].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleInputTab() Model {
	if len(m.input.fields) <= 1 {
		return m
	}

	m.input.fields[m.input.focus].Blur()
	m.input.focus = (m.input.focus + 1) % len(m.input.fields)
	m.input.fields[m.input.focus].Focus()

	switch m.input.mode {
	case inputAddEntryIP:
		if m.input.focus == 1 {
			m.input.mode = inputAddEntryHostname
		}
	case inputAddEntryHostname:
		if m.input.focus == 0 {
			m.input.mode = inputAddEntryIP
		}
	case inputImportName:
		if m.input.focus == 1 {
			m.input.mode = inputImportURL
		}
	case inputImportURL:
		if m.input.focus == 0 {
			m.input.mode = inputImportName
		}
	}

	return m
}

func (m Model) handleInputEnter() (tea.Model, tea.Cmd) {
	switch m.input.mode {
	case inputAddProfileName:
		name := strings.TrimSpace(m.input.fields[0].Value())
		if name == "" {
			m.statusText = "Profile name cannot be empty"
			m.statusErr = true
			return m, nil
		}
		m.input = startAddProfileEntries(name)
		return m, nil

	case inputAddProfileEntries:
		line := strings.TrimSpace(m.input.fields[0].Value())
		if line == "" {
			if len(m.pendingEntries) > 0 {
				if err := m.manager.AddProfile(m.input.buf, m.pendingEntries); err != nil {
					m.statusText = fmt.Sprintf("Add profile failed: %s", err)
					m.statusErr = true
				} else {
					m.profiles = m.manager.ListProfiles()
					m.dirty = m.manager.HasDirtyState()
					m.statusText = fmt.Sprintf("Added profile %q with %d entries", m.input.buf, len(m.pendingEntries))
					m.statusErr = false
				}
				m.pendingEntries = nil
			}
			m.input = newInputState()
			m.view = m.prevView
			return m, nil
		}
		parts := strings.Fields(line)
		if len(parts) < 2 || net.ParseIP(parts[0]) == nil {
			m.statusText = "Invalid format. Use: IP hostname"
			m.statusErr = true
			m.input.fields[0].SetValue("")
			return m, nil
		}
		for _, h := range parts[1:] {
			m.pendingEntries = append(m.pendingEntries, hosts.Entry{IP: parts[0], Hostname: h})
		}
		m.statusText = fmt.Sprintf("Added entry: %s → %s (%d total)", parts[0], strings.Join(parts[1:], ", "), len(m.pendingEntries))
		m.statusErr = false
		m.input.fields[0].SetValue("")
		return m, nil

	case inputAddEntryIP, inputAddEntryHostname:
		ip := strings.TrimSpace(m.input.fields[0].Value())
		hostname := strings.TrimSpace(m.input.fields[1].Value())
		if ip == "" || hostname == "" {
			m.statusText = "Both IP and hostname are required"
			m.statusErr = true
			return m, nil
		}
		if net.ParseIP(ip) == nil {
			m.statusText = fmt.Sprintf("Invalid IP: %s", ip)
			m.statusErr = true
			return m, nil
		}
		if err := m.manager.AddEntry(m.detailProfile, hosts.Entry{IP: ip, Hostname: hostname}); err != nil {
			m.statusText = fmt.Sprintf("Add entry failed: %s", err)
			m.statusErr = true
		} else {
			entries, _ := m.manager.GetEntries(m.detailProfile)
			m.entries = entries
			m.dirty = m.manager.HasDirtyState()
			m.statusText = fmt.Sprintf("Added %s → %s", ip, hostname)
			m.statusErr = false
		}
		m.input = newInputState()
		m.view = m.prevView
		return m, nil

	case inputImportName, inputImportURL:
		name := strings.TrimSpace(m.input.fields[0].Value())
		url := strings.TrimSpace(m.input.fields[1].Value())
		if name == "" || url == "" {
			m.statusText = "Both name and URL are required"
			m.statusErr = true
			return m, nil
		}
		if err := m.manager.ImportFromURL(name, url); err != nil {
			m.statusText = fmt.Sprintf("Import failed: %s", err)
			m.statusErr = true
		} else {
			m.profiles = m.manager.ListProfiles()
			m.dirty = m.manager.HasDirtyState()
			m.statusText = fmt.Sprintf("Imported profile %q from URL", name)
			m.statusErr = false
		}
		m.input = newInputState()
		m.view = m.prevView
		return m, nil
	}

	return m, nil
}

func (m Model) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Enter):
		switch m.confirm {
		case confirmRemoveProfile:
			if err := m.manager.RemoveProfile(m.confirmTarget); err != nil {
				m.statusText = fmt.Sprintf("Remove failed: %s", err)
				m.statusErr = true
			} else {
				m.profiles = m.manager.ListProfiles()
				if m.profileCursor >= len(m.profiles) && m.profileCursor > 0 {
					m.profileCursor--
				}
				m.dirty = m.manager.HasDirtyState()
				m.statusText = fmt.Sprintf("Removed profile %q", m.confirmTarget)
				m.statusErr = false
			}
			m.view = viewProfiles

		case confirmRemoveEntry:
			if err := m.manager.RemoveEntry(m.detailProfile, m.confirmTarget); err != nil {
				m.statusText = fmt.Sprintf("Remove entry failed: %s", err)
				m.statusErr = true
			} else {
				entries, _ := m.manager.GetEntries(m.detailProfile)
				m.entries = entries
				if m.entryCursor >= len(m.entries) && m.entryCursor > 0 {
					m.entryCursor--
				}
				m.dirty = m.manager.HasDirtyState()
				m.statusText = fmt.Sprintf("Removed entry %q", m.confirmTarget)
				m.statusErr = false
			}
			m.view = viewDetail

		case confirmApply:
			cmd, err := m.manager.ApplyCmd()
			if err != nil {
				m.statusText = fmt.Sprintf("Apply failed: %s", err)
				m.statusErr = true
				m.view = viewProfiles
				m.confirm = confirmNone
				m.confirmTarget = ""
				return m, nil
			}
			m.confirm = confirmNone
			m.confirmTarget = ""
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return applyDoneMsg{err: err}
			})
		}
		m.confirm = confirmNone
		m.confirmTarget = ""

	case key.Matches(msg, keys.Escape):
		wasEntryConfirm := m.confirm == confirmRemoveEntry
		m.confirm = confirmNone
		m.confirmTarget = ""
		if wasEntryConfirm {
			m.view = viewDetail
		} else {
			m.view = viewProfiles
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	titleText := "hops"
	if m.dirty {
		titleText += "  " + statusErrStyle.Render("● unsaved")
	} else {
		titleText += "  " + statusOkStyle.Render("● synced")
	}
	b.WriteString(titleStyle.Render(titleText) + "\n")
	b.WriteString(dimStyle.Render(fmt.Sprintf("  profiles: %s", m.manager.DataDir())) + "\n\n")

	switch m.view {
	case viewProfiles:
		b.WriteString(renderProfileList(m.profiles, m.profileCursor, m.height))

	case viewDetail:
		b.WriteString(renderDetailView(m.detailProfile, m.entries, m.entryCursor, m.height))

	case viewInput:
		b.WriteString(renderInput(m.input))

	case viewConfirm:
		switch m.confirm {
		case confirmRemoveProfile:
			b.WriteString(confirmStyle.Render(fmt.Sprintf("Remove profile %q? (enter/esc)", m.confirmTarget)))
		case confirmRemoveEntry:
			b.WriteString(confirmStyle.Render(fmt.Sprintf("Remove entry %q? (enter/esc)", m.confirmTarget)))
		case confirmApply:
			b.WriteString(confirmStyle.Render("Apply enabled profiles to /etc/hosts? (sudo required) (enter/esc)"))
		}
	}

	content := b.String()
	contentLines := strings.Count(content, "\n")
	padding := m.height - contentLines - 3
	if padding > 0 {
		b.WriteString(strings.Repeat("\n", padding))
	}

	if m.statusText != "" {
		if m.statusErr {
			b.WriteString(statusErrStyle.Render("✗ " + m.statusText))
		} else {
			b.WriteString(statusOkStyle.Render("✓ " + m.statusText))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	switch m.view {
	case viewProfiles:
		b.WriteString(profileListHelp())
	case viewDetail:
		b.WriteString(detailViewHelp())
	case viewInput:
		b.WriteString(inputHelp())
	case viewConfirm:
		b.WriteString(confirmHelp())
	}

	return b.String()
}
