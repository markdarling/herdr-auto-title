package reads_test

import (
	"testing"
	"time"

	"github.com/kryptamine/herdr-auto-title/internal/herdr"
	"github.com/kryptamine/herdr-auto-title/internal/reads"
	"github.com/kryptamine/herdr-auto-title/internal/reads/readstest"
	"github.com/kryptamine/herdr-auto-title/internal/state"
)

// codexPaneAt is a pane sitting in dir and holding a Codex session, which
// Herdr reports by id alone.
func codexPaneAt(paneID, sessionID, dir string) *state.PaneState {
	return state.PaneFrom(herdr.PaneInfo{
		PaneID: paneID, TabID: "wE:t1", CWD: dir, Agent: "codex",
		AgentSession: &herdr.AgentSessionInfo{
			Agent: "codex", Kind: herdr.SessionRefID, Value: sessionID,
		},
	}, "", time.Time{})
}

func codexOptions(t *testing.T) reads.Options {
	t.Helper()

	opts := transcriptOptions()
	opts.CodexHome = t.TempDir()

	return opts
}

func TestACodexPaneIsNamedFromItsSession(t *testing.T) {
	t.Parallel()

	repo := readstest.Repo(t, "main")
	opts := codexOptions(t)
	readstest.CodexRollout(t, opts.CodexHome, testSession,
		readstest.CodexStartedIn(repo),
		readstest.CodexPrompt("rework the poll loop"),
	)
	readstest.CodexName(t, opts.CodexHome, testSession, "Rework poll loop")

	pane := codexPaneAt("wE:p1", testSession, repo)
	readOne(opts, pane)

	if pane.AgentTopic != "Rework poll loop" {
		t.Errorf("topic = %q, want the session's name", pane.AgentTopic)
	}

	if pane.AgentDir != repo {
		t.Errorf("agent dir = %q, want where the session started, %q", pane.AgentDir, repo)
	}
}

func TestACodexSessionInAWorktreeNamesItsBranch(t *testing.T) {
	t.Parallel()

	// The pane sits at the repository root on the trunk while its Codex was
	// started in a worktree, so the branch worth showing is the session's.
	repo := readstest.Repo(t, "main")
	worktree := readstest.Worktree(t, repo, "oauth", "feat/oauth")
	opts := codexOptions(t)
	readstest.CodexRollout(t, opts.CodexHome, testSession, readstest.CodexStartedIn(worktree))

	pane := codexPaneAt("wE:p1", testSession, repo)
	readOne(opts, pane)

	if got := branchFor(pane); got != "feat/oauth" {
		t.Errorf("branch = %q, want feat/oauth", got)
	}
}

func TestTranscriptsSwitchedOffLeaveACodexSessionUnread(t *testing.T) {
	t.Parallel()

	opts := branchOptions()
	opts.CodexHome = t.TempDir()
	readstest.CodexRollout(t, opts.CodexHome, testSession,
		readstest.CodexStartedIn(t.TempDir()),
		readstest.CodexPrompt("rework the poll loop"),
	)

	pane := codexPaneAt("wE:p1", testSession, t.TempDir())
	readOne(opts, pane)

	if pane.AgentTopic != "" || pane.AgentDir != "" {
		t.Errorf("topic = %q, dir = %q, want neither read", pane.AgentTopic, pane.AgentDir)
	}
}
