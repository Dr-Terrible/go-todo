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

func main() {

	// Load Todo.txt CLI environment variables
	utils.LoadConfig()

	// Initialize the app CLI
	app := cli.NewApp()
	app.Name = "go-todo"
	app.Usage = "A simple and extensible utility for managing your todo.txt files"
	app.Version = "1.0.1"
	app.Author = "Mauro Toffanin"
	app.Email = "toffanin.mauro@gmail.com"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{"t", "Prefixes the current date to a task automatically when it's added"},
		cli.BoolFlag{"T", "Do not prefix the current date to a task automatically when it's added."},
		cli.BoolFlag{"f", "Forces actions without confirmation or interactive input."},
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
