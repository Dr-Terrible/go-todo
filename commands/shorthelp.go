package commands

import (
	"github.com/codegangsta/cli"
)

func GetShorthelpCommand() cli.Command {

	return cli.Command{
		Name:  "shorthelp",
		Usage: "Show a usage message briefly summarizing all commands (a synonym for -h)",
		Description: `
   This command prints a summary of the command-line usage of 'todo' and all its
   add-ons, then exit.

   The 'shorthelp' command is supported only for backward-compatibility with the
   original Todo.txt CLI from Gina Trapani, and falls back on the standard
   help command.

   You should use the POSIX-compliant option 'help, -h' instead.
`,
		Action: func(c *cli.Context) {
			cli.ShowAppHelp(c)
		},
	}
}
