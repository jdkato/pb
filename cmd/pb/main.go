package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jdkato/pb/internal/cli"
	"github.com/jdkato/pb/internal/config"
	"github.com/jdkato/pb/internal/platform"
	"github.com/pterm/pterm"
	"github.com/spf13/pflag"
)

// version is set during the release build process.
var version = "master"

func doConvert(c platform.Converter) error {
	spinner, err := pterm.DefaultSpinner.Start("Uploading ...")
	if err != nil {
		return err
	}

	sites := strings.Split(cli.Flags.To, ",")
	for _, dest := range sites {
		if err = c.To(dest); err != nil {
			return err
		}
	}
	done := fmt.Sprintf("Uploaded draft to %s", cli.ToSentence(sites, "and"))

	spinner.Success(done)
	return nil
}

func main() {
	pflag.Parse()

	args := pflag.Args()
	argc := len(args)

	if cli.Flags.Version {
		fmt.Println("pb version " + version)
		os.Exit(0)
	} else if cli.Flags.Help || argc == 0 {
		pflag.Usage()
		os.Exit(0)
	}

	cmd, exists := cli.Actions[args[0]]
	if exists {
		if err := cmd(args[1:], config.Auth); err != nil {
			cli.ShowError(
				fmt.Sprintf("Command '%s' failed", args[0]), err.Error())
		}
		os.Exit(0)
	}

	_, err := exec.LookPath("inkscape")
	if err != nil && cli.Flags.Inkscape != "" {
		cli.ShowError(
			"Please add 'inkscape' to your $PATH or specify its location with `--inkscape`.",
			err.Error())
	}

	converter, err := platform.NewConverter(args[0])
	if err != nil {
		cli.ShowError("Failed to read document", err.Error())
	}

	if err = doConvert(converter); err != nil {
		cli.ShowError("Failed to upload document", err.Error())
	}

	os.Exit(0)
}
