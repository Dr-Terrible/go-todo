package commands

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/toffanin/go-todo/utils"

	"github.com/codegangsta/cli"
)

func GetAdd() cli.Command {

	return cli.Command{
		Name:      "add",
		ShortName: "a",
		Usage:     "Add a task to your todo.txt file",
		Description: `
   This command can be used to add the specified task to your todo.txt file on
   its own line.

   Project and content notation are optional. Quotation marks are optional too.

EXAMPLES

   Adds a simple task (quotes are optional):

	  $ todo add "Move out cardboard boxes from the garage"

   Adds tasks with a project notation (quotes are optional):

	  $ todo add "Move out cardboard boxes from the garage +cleaning"
	  $ todo add "Do a load of laundry +cleaning"
	  $ todo add "Vacuum the house +cleaning"

   Adds tasks with a context notation (quotes are optional):

	  $ todo add "Buy eggs and milk @grocery"
	  $ todo add "Buy a cake for friday's dinner party with friends @backery"

   Adds tasks with both project and context notation (quotes are optional):

	  $ todo add "Feed the kitten +BellyOfTheBeast"
	  $ todo add "Buy food with amino acid taurine @petshop +BellyOfTheBeast"
	  $ todo add "Buy huge amont of meat @butcher +BellyOfTheBeast"
	  $ todo add "Hire a bouncer to protect @kitchen cupboard from the cat +BellyOfTheBeast"
`,
		Action: func(c *cli.Context) {
			// collect all the user-submitted arguments in an array
			args := c.Args()

			// debugging
			/*fmt.Printf("(add::Action) args (%d): %s\n", len(args), args)
			fmt.Printf("(add::Action) global flag: %s (%t)\n", "-t", c.GlobalBool("t"))*/

			// task mangler
			task := ""
			switch {
			case len(args) == 0: // no options specified

				// check incorrect usage of the command
				if c.GlobalBool("f") {
					fmt.Print("\nDetected missing option with command \"add [task]\"\n")
					fmt.Print("Usage: todo -f add [task]\"\n\n")
					cli.ShowCommandHelp(c, "add")
					return
				}

				// invoke interactive input
				task = utils.InteractiveInput("Add:")

			default: // collect all the arguments into a single string
				task = strings.Join(args[0:], " ")
			}

			// TODO: validating input as a task

			// sanitize input
			task = utils.SanitizeInput(task)

			// honor the -t global flag
			if c.GlobalBool("t") {
				date := time.Now().Format("2006-01-02 ")
				task = date + task
			}

			// save the new task
			addAction(task)
		},
	}
}

func GetAddm() cli.Command {

	return cli.Command{
		Name:      "addm",
		ShortName: "",
		Usage:     "Add multiple tasks to your todo.txt file",
		Description: `
   This command can be used to add the specified tasks to your todo.txt file.

   Project and content notation are optional. Quotation marks are optional too.

EXAMPLES

   Adds some simple tasks (quotes are optional):

	  $ todo addm "Buy eggs and milk @grocery"
	  $ > Buy a cake for friday's dinner party with friends @backery
`,
		Action: func(c *cli.Context) {
			// collect all the user-submitted arguments in an array
			args := c.Args()

			// check incorrect usage of the command
			if len(args) == 0 {
				fmt.Print("\nDetected missing option with command \"addm [task]\"\n")
				fmt.Print("Usage: todo addm [task]\"\n\n")
				cli.ShowCommandHelp(c, "addm")
				return
			}

			// collect all the arguments into a single string
			firstTask := strings.Join(args[0:], " ")

			// invoke interactive input
			secondTask := utils.InteractiveInput(">")

			// TODO: validating input as a task

			// sanitize tasks
			firstTask = utils.SanitizeInput(firstTask)
			secondTask = utils.SanitizeInput(secondTask)

			// honor the -t global flag
			if c.GlobalBool("t") {
				date := time.Now().Format("2006-01-02 ")
				firstTask = date + firstTask
				secondTask = date + secondTask
			}

			// save taska
			addAction(firstTask)
			addAction(secondTask)
		},
	}
}

// Adds a task to a todo.txt file.
func addAction(task string) {

	todoFile := utils.GetSetting("TODO_FILE")
	todoDir := path.Dir(todoFile)
	//fmt.Printf("*DIR: %s\n", todoDir)

	// validate dir path
	finfo, err := os.Stat(todoDir)
	//fmt.Printf("FileMode IsDir: %t\n", finfo.Mode().IsDir())
	//fmt.Printf("FileMode IsRegular: %t\n", finfo.Mode().IsRegular())
	//fmt.Printf("FileMode Perm: %t\n", finfo.Mode().Perm())
	if err != nil {
		// path doesn't exists
		if os.IsNotExist(err) {
			fmt.Printf("DIR:%s doesn't exists.\n", todoDir)
			fmt.Printf("Please create the missing directory with: `mkdir -p %s`.\n\n", todoDir)
			os.Exit(1)
		}

		// brace yourself: unknown errors are coming
		utils.Check(err)
	}

	// path exists but is not a directory
	if !finfo.IsDir() {
		fmt.Printf("DIR:%s is not a directory.\n", todoDir)
		fmt.Println("Please fix your todo.cfg file and be sure to specify a directory with an absolute path.")
		os.Exit(1)
	}

	// determine the number of tasks in todo.txt
	// TODO: with NewReadWriter the code should be more compact
	//       buf := bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))
	fd, err := os.OpenFile(todoFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		// file isn't readable
		if os.IsPermission(err) {
			fmt.Printf("%s or %s don't have correct permission bits.\n", todoFile, todoDir)
			fmt.Println("Please fix the directory / file permissions")
			os.Exit(1)
		}

		// brace yourself: unknown errors are coming
		utils.Check(err)
	}

	scanner := bufio.NewScanner(fd)
	ntasks := 1
	for scanner.Scan() {
		ntasks++
	}
	if err := scanner.Err(); err != nil {
		fd.Close()
		utils.Check(err)
	}
	//fmt.Printf("n. lines: %d\n", ntasks)
	err = fd.Close()
	utils.Check(err)

	// Open todo.txt in append mode only
	fd, err = os.OpenFile(todoFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	utils.Check(err)
	defer fd.Close()

	// use buffered I/O
	writer := bufio.NewWriter(fd)

	// add the task to todo.txt
	_, err = writer.WriteString(task + "\n")
	utils.Check(err)
	err = writer.Flush()
	utils.Check(err)

	// print summary
	fmt.Printf("%d: %s\n", ntasks, task)
	fmt.Printf("TODO: %d added\n", ntasks)
}
