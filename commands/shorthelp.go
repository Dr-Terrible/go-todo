// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"github.com/codegangsta/cli"
)

func GetShorthelp() cli.Command {

	return cli.Command{
		Name:  "shorthelp",
		Usage: "Shows a usage message briefly summarizing all commands (a synonym for -h)",
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
