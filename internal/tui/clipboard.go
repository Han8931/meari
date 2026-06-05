package tui

import (
	"os"

	"github.com/atotto/clipboard"
	osc52 "github.com/aymanbagabas/go-osc52/v2"
)

// copyToClipboard puts text on the system clipboard two ways at once: the
// native clipboard (pbcopy / xclip / wl-copy via atotto) for local sessions,
// and an OSC 52 escape written to the terminal so copying also works over SSH
// in terminals that support it. It is a package variable so tests can stub it
// out instead of touching the real clipboard.
//
// The returned error reflects only the native path — OSC 52 is fire-and-forget
// (a terminal that doesn't support it silently ignores the sequence).
var copyToClipboard = func(text string) error {
	_, _ = osc52.New(text).WriteTo(os.Stderr)
	return clipboard.WriteAll(text)
}

// pasteFromClipboard reads the system clipboard (native only — terminals don't
// allow OSC 52 reads for security). A package variable so tests can stub it.
var pasteFromClipboard = func() (string, error) {
	return clipboard.ReadAll()
}
