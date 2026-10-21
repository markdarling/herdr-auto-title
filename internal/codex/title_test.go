package codex

import (
	"path/filepath"
	"testing"

	"github.com/kryptamine/herdr-auto-title/internal/transcript"
)

func TestATitleShowsTheSessionItNames(t *testing.T) {
	t.Parallel()

	named := transcript.Topic{Name: "Rework the poll loop so a pane that moved is read again"}

	for _, title := range []string{
		"Rework the poll loop so a pane that moved is read... | dashboard",
		"[ ! ] Action Required | Rework the poll loop | dashboard",
		"",
	} {
		if !shows(title, named) {
			t.Errorf("title %q does not show %q", title, named.Name)
		}
	}

	// A session Codex has not named yet titles its terminal with the project.
	if !shows("dashboard", transcript.Topic{Opening: "rework the poll loop"}) {
		t.Error("an unnamed title does not show an unnamed session")
	}
}

func TestATitleThatMovedOnShowsAnotherSession(t *testing.T) {
	t.Parallel()

	// After /new or /resume Herdr still reports the session before, until the
	// next prompt, while the title already belongs to the next one.
	named := transcript.Topic{Name: "Reply with ok"}

	if shows("dashboard", named) {
		t.Error("a fresh session's title shows the named session before it")
	}

	if shows("Reply yes | dashboard", named) {
		t.Error("a resumed session's title shows the session before it")
	}

	if shows(
		"Reply with ok | dashboard",
		transcript.Topic{Opening: "reply with the single word ok"},
	) {
		t.Error("a named title shows an unnamed session")
	}
}

func TestAConfiguredTitleKeepsEverySessionRead(t *testing.T) {
	t.Parallel()

	// A title Codex was told to set otherwise carries no name to compare, and
	// reading it as stale would drop every session's name.
	h := newHome(t)
	h.rollout(meta(started), prompt("reply with ok"))
	h.index(named(session, "Reply with ok"))

	r := NewReader(h.dir)
	topic := r.Topic(session)

	if r.Shows("dashboard", topic) {
		t.Fatal("a default title of another session shows this one")
	}

	h.write(filepath.Join(h.dir, "config.toml"), "[tui]", `terminal_title = ["spinner", "project"]`)

	if !r.Shows("dashboard", topic) {
		t.Error("a configured title is read as another session's")
	}
}

func TestOnlyAnAssignedTerminalTitleCountsAsConfigured(t *testing.T) {
	t.Parallel()

	for config, want := range map[string]bool{
		"[tui]\nterminal_title = [\"project\"]": true,
		"tui.terminal_title = []":               true,
		"tui = { terminal_title = [] }":         true,
		"# terminal_title = [\"project\"]":      false,
		"[tui]\nnotifications = true":           false,
	} {
		if got := setsTerminalTitle(config); got != want {
			t.Errorf("setsTerminalTitle(%q) = %v, want %v", config, got, want)
		}
	}
}
