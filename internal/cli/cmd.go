package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jdkato/pb/internal/config"
	"github.com/pterm/pterm"
)

// Actions are the available CLI commands.
var Actions = map[string]func([]string, *config.Config) error{
	"configure": configure,
}

var commandInfo = map[string]string{
	"configure": "Run an interactive configuration wizard.",
}

func configure(args []string, cfg *config.Config) error {
	var token string
	var platform string

	p1 := &survey.Select{
		Message: "Choose a platform:",
		Options: []string{"Medium"},
		Default: "Medium",
	}

	err := survey.AskOne(p1, &platform)
	if err != nil {
		return err
	}

	p2 := &survey.Password{}

	switch platform {
	case "Medium":
		p2.Message = "Enter your Integration Token:"
	case "DEV":
		p2.Message = "Enter your Community API Key:"
	case "Hashnode":
		p2.Message = "Enter your Personal Access Token:"
	}
	survey.AskOne(p2, &token)

	err = cfg.AddToken(token, platform)
	if err != nil {
		return err
	}

	err = cfg.Save()
	if err != nil {
		return err
	}

	pterm.Success.Print("Updated configuration file.")
	pterm.FgGray.Print("\n└ " + fmt.Sprintf("(%s)\n", cfg.Path))

	return nil
}
