package hosts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Entry struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

type ProfileInfo struct {
	Name       string
	Enabled    bool
	EntryCount int
}

type stateFile struct {
	Enabled map[string]bool `json:"enabled"`
}

type Manager struct {
	dataDir   string
	hostsPath string
	state     stateFile
}

func NewManager(hostsPath string) (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine home directory: %w", err)
	}
	dataDir := filepath.Join(home, ".local", "share", "hops", "profiles")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create data directory: %w", err)
	}

	m := &Manager{
		dataDir:   dataDir,
		hostsPath: hostsPath,
		state:     stateFile{Enabled: map[string]bool{}},
	}
	m.loadState()
	return m, nil
}

func (m *Manager) stateFilePath() string {
	return filepath.Join(filepath.Dir(m.dataDir), "state.json")
}

func (m *Manager) profilePath(name string) string {
	return filepath.Join(m.dataDir, name+".hosts")
}

func (m *Manager) loadState() {
	data, err := os.ReadFile(m.stateFilePath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &m.state)
	if m.state.Enabled == nil {
		m.state.Enabled = map[string]bool{}
	}
}

func (m *Manager) saveState() error {
	data, err := json.MarshalIndent(m.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.stateFilePath(), data, 0644)
}

func (m *Manager) ListProfiles() []ProfileInfo {
	entries, err := os.ReadDir(m.dataDir)
	if err != nil {
		return nil
	}
	var profiles []ProfileInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".hosts") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".hosts")
		hostEntries, _ := m.readProfileFile(name)
		profiles = append(profiles, ProfileInfo{
			Name:       name,
			Enabled:    m.state.Enabled[name],
			EntryCount: len(hostEntries),
		})
	}
	return profiles
}

func (m *Manager) GetEntries(name string) ([]Entry, error) {
	return m.readProfileFile(name)
}

func (m *Manager) readProfileFile(name string) ([]Entry, error) {
	f, err := os.Open(m.profilePath(name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseHostsReader(f)
}

func (m *Manager) writeProfileFile(name string, entries []Entry) error {
	var b strings.Builder
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("%s %s\n", e.IP, e.Hostname))
	}
	return os.WriteFile(m.profilePath(name), []byte(b.String()), 0644)
}

func (m *Manager) Toggle(name string) error {
	m.state.Enabled[name] = !m.state.Enabled[name]
	return m.saveState()
}

func (m *Manager) AddProfile(name string, entries []Entry) error {
	if err := m.writeProfileFile(name, entries); err != nil {
		return err
	}
	m.state.Enabled[name] = false
	return m.saveState()
}

func (m *Manager) RemoveProfile(name string) error {
	if err := os.Remove(m.profilePath(name)); err != nil && !os.IsNotExist(err) {
		return err
	}
	delete(m.state.Enabled, name)
	return m.saveState()
}

func (m *Manager) AddEntry(profileName string, e Entry) error {
	entries, err := m.readProfileFile(profileName)
	if err != nil {
		return err
	}
	entries = append(entries, e)
	return m.writeProfileFile(profileName, entries)
}

func (m *Manager) RemoveEntry(profileName, hostname string) error {
	entries, err := m.readProfileFile(profileName)
	if err != nil {
		return err
	}
	var filtered []Entry
	for _, e := range entries {
		if e.Hostname != hostname {
			filtered = append(filtered, e)
		}
	}
	return m.writeProfileFile(profileName, filtered)
}

func (m *Manager) ImportFromURL(profileName, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	entries, err := parseHostsReader(resp.Body)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("no valid host entries found at %s", url)
	}
	return m.AddProfile(profileName, entries)
}

// ApplyCmd builds the exec.Cmd that writes enabled profiles to /etc/hosts via sudo.
// The block is written to a temp file so sudo can prompt interactively for password.
func (m *Manager) ApplyCmd() (*exec.Cmd, error) {
	block := m.buildApplyBlock()

	tmpFile := filepath.Join(os.TempDir(), "hops-apply.tmp")
	if err := os.WriteFile(tmpFile, []byte(block), 0644); err != nil {
		return nil, fmt.Errorf("cannot write temp file: %w", err)
	}

	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("cannot find own executable: %w", err)
	}

	return exec.Command("sudo", exe, "--apply-write", m.hostsPath, "--block-file", tmpFile), nil
}

// ApplyWrite is the privileged helper: reads the managed block from a file,
// splices it into the hosts file replacing the old managed section.
func ApplyWrite(hostsPath, blockFile string) error {
	block, err := os.ReadFile(blockFile)
	if err != nil {
		return err
	}
	defer os.Remove(blockFile)

	existing, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	result := spliceBlock(string(existing), string(block))
	return os.WriteFile(hostsPath, []byte(result), 0644)
}

const (
	markerBegin = "# --- BEGIN hops managed ---"
	markerEnd   = "# --- END hops managed ---"
)

func (m *Manager) buildApplyBlock() string {
	var b strings.Builder
	b.WriteString(markerBegin + "\n")

	profiles := m.ListProfiles()
	for _, p := range profiles {
		if !p.Enabled {
			continue
		}
		entries, err := m.readProfileFile(p.Name)
		if err != nil || len(entries) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("# profile: %s\n", p.Name))
		for _, e := range entries {
			b.WriteString(fmt.Sprintf("%s %s\n", e.IP, e.Hostname))
		}
	}

	b.WriteString(markerEnd + "\n")
	return b.String()
}

func spliceBlock(existing, block string) string {
	beginIdx := strings.Index(existing, markerBegin)
	endIdx := strings.Index(existing, markerEnd)

	if beginIdx >= 0 && endIdx >= 0 {
		endIdx += len(markerEnd)
		// skip trailing newline after end marker
		if endIdx < len(existing) && existing[endIdx] == '\n' {
			endIdx++
		}
		return existing[:beginIdx] + block + existing[endIdx:]
	}

	// No existing block — append
	s := strings.TrimRight(existing, "\n") + "\n\n" + block
	return s
}

func (m *Manager) HasDirtyState() bool {
	// Compare what's in /etc/hosts managed block vs what we'd write
	existing, err := os.ReadFile(m.hostsPath)
	if err != nil {
		return true
	}
	currentBlock := extractBlock(string(existing))
	desiredBlock := m.buildApplyBlock()
	return currentBlock != desiredBlock
}

func extractBlock(content string) string {
	beginIdx := strings.Index(content, markerBegin)
	endIdx := strings.Index(content, markerEnd)
	if beginIdx < 0 || endIdx < 0 {
		return ""
	}
	endIdx += len(markerEnd)
	if endIdx < len(content) && content[endIdx] == '\n' {
		endIdx++
	}
	return content[beginIdx:endIdx]
}

func parseHostsReader(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		ip := fields[0]
		if net.ParseIP(ip) == nil {
			continue
		}
		for _, h := range fields[1:] {
			if strings.HasPrefix(h, "#") {
				break
			}
			entries = append(entries, Entry{IP: ip, Hostname: h})
		}
	}
	return entries, scanner.Err()
}

func (m *Manager) Close() {}

func (m *Manager) DataDir() string {
	return filepath.Dir(m.dataDir)
}
