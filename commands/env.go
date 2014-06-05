// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/toffanin/go-todo/utils"

	"github.com/codegangsta/cli"
)

func GetEnv() cli.Command {

	return cli.Command{
		Name:  "env",
		Usage: "Display information about the `todo` environment",
		Description: `
   By default 'env' displays information as a shell script.

   The environment info will be dumped in a straight-forward form suitable for
   sourcing into a shell script.

   If one or more variable names is given as arguments, 'env' displays the value
   of each named variable on its own line.

   The 'env' environment can be controlled through todo.cfg files (see section
   CONFIGURATION FILES) and environment variables.

CONFIGURATION FILES:

   Command line argument defaults can be set globally in a ~/.todo.cfg file or
   set individual in a .todo.cfg for a specific directory.

   Configuration files are simple text files with the following syntax:

   # This is just an example
   TODO_DIR="$HOME/todo"
   TODO_FILE="$TODO_DIR/todo.txt"
   DONE_FILE="$TODO_DIR/done.txt"
   REPORT_FILE="$TODO_DIR/report.txt"

   For backward compatibility, the following syntax is accepted too:

   # This is just an example
   export TODO_DIR="$HOME/todo"
`,
		Action: func(c *cli.Context) {
			// collect all the user-submitted arguments in an array
			args := c.Args()

			// debugging
			/*fmt.Printf("pwd: %s\n", pwd)
			fmt.Printf("HOME=\"%s\"\n", ENV["HOME"])*/

			switch args.Present() {
			case true:
				// print only the required environment variables
				for _, arg := range args {
					fmt.Printf("%s=\"%s\"\n", arg, utils.GetSetting(arg))
				}
			case false:
				// print all the environment variables
				for k, v := range utils.GetSettings() {
					fmt.Printf("%s=\"%s\"\n", k, v)
				}
			}
		},
	}
}
