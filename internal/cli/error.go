package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/pterm/pterm"
)

var logger = log.New(os.Stderr, "", 0)

// NewError creates a colored error from the given information.
//
// The standard format is
//
// ```
// <code> [<context>] <title>
//
// <msg>
// ```
func NewError(title, msg string) error {
	return fmt.Errorf(
		"%s\n\n%s\n\n%s",
		pterm.Error.Sprintf(title),
		msg,
		pterm.Gray(pterm.Italic.Sprintf("Execution stopped with code 1.")),
	)
}

// ShowError displays the given error in the user-specified format.
func ShowError(title, msg string) {
	err := NewError(title, msg)

	logger.SetOutput(os.Stderr)
	logger.Println(err)

	os.Exit(1)
}
