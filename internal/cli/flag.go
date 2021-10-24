package cli

import (
	"github.com/spf13/pflag"
)

// CLIFlags holds the values that are defined at rumtime by the user.
//
// For example, `fsm --json pipeline.yml`.
type CLIFlags struct {
	To       string
	Help     bool
	Version  bool
	ImageDir string
}

var Flags CLIFlags

var shortcodes = map[string]string{
	"help":      "h",
	"to":        "t",
	"version":   "v",
	"image-dir": "d",
}

func init() {
	pflag.StringVarP(&Flags.To, "to", "t", "medium", "Comma-delimited list of destination platforms.")
	pflag.StringVarP(&Flags.ImageDir, "image-dir", "d", "", "Search directory for local images.")
	pflag.BoolVarP(&Flags.Help, "help", "h", false, "Print this help message.")
	pflag.BoolVarP(&Flags.Version, "version", "v", false, "Print the current version.")
}
