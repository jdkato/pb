package cli

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"github.com/spf13/pflag"
)

var intro = fmt.Sprintf(`pb - %s.

%s:	%s
	%s
	%s

pb is a tool for cross-posting Markdown content while preserving structural
elements (math typesetting, syntax highlighting, diagrams, etc.) across
multiple platforms.`,

	pterm.Italic.Sprintf("A multi-platform publishing workflow"),

	pterm.Bold.Sprintf("Usage"),
	pterm.Gray("pb [options] [command] [arguments...]"),
	pterm.Gray("pb --to medium file.md"),
	pterm.Gray("pb configure"),
)

func toFlag(name string) string {
	code := shortcodes[name]
	return fmt.Sprintf("%s, %s", pterm.Gray("-"+code), pterm.Gray("--"+name))
}

func init() {
	pflag.Usage = func() {
		fmt.Println(intro)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetAutoWrapText(false)

		fmt.Println(pterm.Bold.Sprintf("\nFlags:"))
		pflag.VisitAll(func(f *pflag.Flag) {
			table.Append([]string{toFlag(f.Name), f.Usage})
		})

		table.Render()
		table.ClearRows()

		fmt.Println(pterm.Bold.Sprintf("Commands:"))
		for cmd, use := range commandInfo {
			table.Append([]string{pterm.Gray(cmd), use})
		}
		table.Render()
	}
}
