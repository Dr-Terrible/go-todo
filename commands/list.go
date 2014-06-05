// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"

	"github.com/toffanin/go-todo/library/v1"
	"github.com/toffanin/go-todo/utils"
)

// scan todo.txt file and collect all the tasks
func listAllTasks(file *os.File) {
	reader := todotxt.NewReader(file)
	tasks, err := reader.ReadAll()
	utils.Check(err)
	//fmt.Printf("Tasks: %st\n", tasks)

	// TODO: run the filter command, then sort and mangle output some more

	// TODO: Build and apply the filter

	// print output
	ntasks := reader.Len()
	padding := len(strconv.FormatUint(ntasks, 10))
	i := uint64(0)
	for i < ntasks {
		s := strconv.FormatUint(tasks[i].Id, 10)
		// TODO: console colours
		fmt.Printf("%s: %s\n", utils.PaddingLeft(s, "0", padding), tasks[i].Todo)
		i++
	}

	// if required print verbose info
	verbose, err := strconv.Atoi(utils.GetSetting("TODOTXT_VERBOSE"))
	utils.Check(err)

	switch {
	case verbose == 0:
		return
	case verbose >= 1:
		if verbose > 1 {
			// TODO
			fmt.Printf("TODO DEBUG: Filter Command was: %#t\n", nil)
		}
		fmt.Println("--")
		fmt.Printf("TODO: %d of %d tasks shown\n", ntasks, ntasks)
	}
}

func GetList() cli.Command {

	return cli.Command{
		Name:      "list",
		ShortName: "ls",
		Usage:     "Displays all the tasks with line numbers",
		Description: `
   This command lists all the tasks inside a todo.txt file, one per line. For
   every single task the 'id' number is printed for each lines.

   If one or more TERM(s) is given as arguments, 'list' displays all the tasks
   whose text contains TERM(s), sorted by priority. Instead tasks are always
   sorted alphabetically if no TERM(s) is specified.

   The user can supplies TERM(s) as arguments separated by logical operators.
   These operators control the behaviour of the 'list' command (see section
   OPERATORS).

   Logical operator 'and' is always assumed where the operator is omitted.
   Quotation marks around a logical statement are optional.

OPERATORS:

   Logical operators listed in order of decreasing precedence:

   TERM1 TERM2 [...]
   TERM1, TERM2, [...]
   TERM1 and TERM2 and [...]
      This is the default behaviour and the syntax expresses a logical
      AND (conjunction). The command 'list' displays only the tasks that contain
      all the specified TERM(s).

   TERM1 | TERM2 | [...]
   TERM1 || TERM2 || [...]
   TERM1 or TERM2 or [...]
      This syntax expresses a logical OR (disjunction).
      The command 'list' displays only the tasks that contain any of the
      specified TERM(s).

   -TERM1
   !TERM1
   -not TERM1
      This syntax expresses a logical NOT (negation).
      The command 'list' hides all the task that contain TERM(s)

EXAMPLES:

   Given this todo.txt as a reference:
      Do a load of laundry +cleaning
      Vacuum the house +cleaning
      Buy eggs, cheese and milk @grocery
      Buy a cake for Friday's dinner party with friends @grocery
      Cook an omelet with eggs, cheese and veggies for Mary's @lunch


   Lists all tasks whose text contain the word 'milk':

      $ todo list milk
      > Buy eggs and milk @grocery

   Lists all tasks whose text contain both 'cheese' and 'egg':

      $ todo list cheese egg
      > Buy eggs, cheese and milk @grocery
      > Cook an omelet with eggs, cheese and veggies for Mary's @lunch

   Lists all tasks that belong to context '@grocery' but exclude the one that
   contain the word 'egg':

      $ todo list @grocery -egg
      > Buy a cake for Friday's dinner party with friends @grocery

   Lists all the tasks that belong to any of context '@grocery' or '@lunch':

      $ todo list @grocery or @lunch
      > Buy eggs, cheese and milk @grocery
      > Buy a cake for Friday's dinner party with friends @grocery
      > Cook an omelet with eggs, cheese and veggies for Mary's @lunch
`,
		Action: func(c *cli.Context) {
			// collect all the user-submitted arguments in an array
			args := c.Args()
			nargs := len(args)

			// open todo.txt file
			todoFile := utils.GetSetting("TODO_FILE")
			file, err := os.Open(todoFile)
			utils.Check(err)
			defer file.Close()

			switch nargs {
			case 0:
				listAllTasks(file)
			}

			// debugging
			/*fmt.Println("[todo:list] ConfPaths (filtered): ", utils.ConfPaths)
			fmt.Println("[todo:list] Settings: (filtered): ", utils.Settings)*/

		},
	}
}
