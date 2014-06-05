// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/toffanin/go-todo/commands"
	"github.com/toffanin/go-todo/utils"

	"github.com/codegangsta/cli"
	//"github.com/fatih/color"
)

var (
	appName = "todo"

	// The text template for the Default help topic
	appHelpTemplate = `
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   [environment variables] {{.Name}} [global options] command [...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
ENVIRONMENT VARIABLES:
   TODOTXT_AUTO_ARCHIVE=0,1{{"\t"}}is equivalent to global options -a (0) / -A (1)
   TODOTXT_CFG_FILE=CONFIG_FILE{{"\t"}}is equivalent to global option -d CONFIG_FILE
   TODOTXT_FORCE=1{{"\t"}}is equivalent to global option -f
   TODOTXT_PRESERVE_LINE_NUMBERS=0,1{{"\t"}}is equivalent to global options -n (0) / -N (1)
   TODOTXT_PLAIN=0,1{{"\t"}}is equivalent to global options -p (1) / -c (0)
   TODOTXT_DATE_ON_ADD=0,1{{"\t"}}is equivalent to global options -t (1) / -T (0)
   TODOTXT_VERBOSE=1{{ "\t" }}is equivalent to global option -v
   TODOTXT_DISABLE_FILTER=1{{ "\t" }}is equivalent to global option -x
   TODOTXT_DEFAULT_ACTION=""{{ "\t" }}run this when called with no arguments
   TODOTXT_SORT_COMMAND="sort ..."{{ "\t" }}customize list output
   TODOTXT_FINAL_FILTER="sed ..."{{ "\t" }}customize list after color, P@+ hiding
   TODOTXT_SOURCEVAR=\$DONE_FILE{{ "\t" }}use another source for listcon, listproj

`

	// The text template for the command help topic.
	commandHelpTemplate = `
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   ` + appName + ` {{.Name}} [options] [arguments...]

DESCRIPTION:
   {{.Description}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}

`
)

func main() {

	// Load Todo.txt CLI environment variables
	utils.LoadConfig()

	// Initialize the templates for help sections
	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	// Initialize the app CLI
	app := cli.NewApp()

	app.Name = appName
	app.Usage = "A simple and extensible utility for managing your todo.txt files"
	app.Version = "1.0.1"
	app.Author = "Mauro Toffanin"
	app.Email = "toffanin.mauro@gmail.com"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{"t", "Prefixes the current date to a task automatically when it's added"},
		cli.BoolFlag{"T", "Do not prefix the current date to a task automatically when it's added"},
		cli.BoolFlag{"f", "Forces actions without confirmation or interactive input"},
	}
	app.Commands = []cli.Command{
		commands.GetEnv(),
		commands.GetInit(),
		commands.GetShorthelp(),
		commands.GetAdd(),
		commands.GetAddm(),
		commands.GetList(),
		/*{
			Name:  "status",
			Usage: "Obtain a summary of the todo.txt structure",
			Flags: []cli.Flag{
				cli.StringFlag{"dest, d", "", "specifies a different destination path"},
			},
			Action: func(c *cli.Context) {
				//fmt.Println("status: ", c.Args().First())
			},
		},*/
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
