package resolver

import (
	"strings"
	"unicode"
)

// codexJoin joins the items of Codex's title, which by default reads
// `<activity> <thread name> | <project>` and drops the name until it has one.
const codexJoin = " | "

// codexAttention are the two frames of the prefix Codex blinks in front of its
// title while it waits on the user.
var codexAttention = []string{"[ ! ] Action Required", "[ . ] Action Required"}

// withoutCodexFrame reports the thread name a Codex title carries, and false
// when it carries none. See docs/architecture/title-resolution.md for the
// layout this reads.
func withoutCodexFrame(title string) (string, bool) {
	// Herdr strips the spinner it recognizes, but not the voice mark before it.
	rest := strings.TrimLeftFunc(title, func(r rune) bool {
		return r == '●' || unicode.Is(unicode.Braille, r) || unicode.IsSpace(r)
	})

	for _, prefix := range codexAttention {
		if after, found := strings.CutPrefix(rest, prefix); found {
			rest = strings.TrimPrefix(after, codexJoin)

			break
		}
	}

	// The project always comes last, and alone it is the whole title.
	cut := strings.LastIndex(rest, codexJoin)
	if cut < 0 {
		return "", false
	}

	rest = strings.TrimSpace(rest[:cut])

	return rest, rest != ""
}
