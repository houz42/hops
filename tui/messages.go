package tui

import "github.com/houz42/hops/hosts"

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type statusMsg string

type profilesLoadedMsg struct {
	profiles []hosts.ProfileInfo
}

type entriesLoadedMsg struct {
	profile string
	entries []hosts.Entry
}

type applyDoneMsg struct{ err error }

type confirmAction int

const (
	confirmNone confirmAction = iota
	confirmRemoveProfile
	confirmRemoveEntry
	confirmApply
)
