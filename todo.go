package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var (
	/* This map stores all the environment variables related to todo.txt CLI
	 */
	ENV = map[string]string{
		"HOME":                 "",
		"TODO_DIR":             "",
		"TODO_FILE":            "",
		"DONE_FILE":            "",
		"REPORT_FILE":          "",
		"TODO_ACTIONS_DIR":     "",
		"TODOTXT_SORT_COMMAND": "",
		"TODOTXT_FINAL_FILTER": "",
	}

	/* This slice defines all the possible paths for the configuration files.
	 * The slice defines also the exact order at which the conf files should be
	 * loaded by the godotenv library.
	 *
	 * $HOME/todo.cfg
	 * $HOME/.todo.cfg
	 * ./todo.cfg
	 * /etc/todo/config
	 */
	// TODO: transform the slice into a map[string]bool where the boolean value defines if the config file exists
	CFG = []string{
		"$HOME/todo.cfg",
		"$HOME/.todo.cfg",
		"./todo.cfg",
		"/etc/todo/config"}
)

// TODO: this is an ugly and hackish error handler that needs to be improved
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Exists returns true if the given path exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Adds a task to a todo.txt file.
func addTask(task string) {
	// TODO: ENV["TODO_FILE"] path should be validated somehow
	// before to be stated by os.OpenFile
	//path.Clean(ENV["TODO_FILE"])

	// Open todo.txt in append mode only
	todoFile, err := os.OpenFile(ENV["TODO_FILE"], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	defer todoFile.Close()
	check(err)

	// add the task to the todo.txt file
	ret, err := todoFile.WriteString(task + "\n")
	fmt.Printf("added new task (%d bytes): \"%s\"\n", ret, task)

	// sync / flush todo.txt
	todoFile.Sync()
}

func main() {

	// Retrieve environment variable $HOME
	ENV["HOME"] = os.Getenv("HOME")
	if ENV["HOME"] != "" {
		// $HOME isn't empty, we expand $HOME for CFG[0] and CFG[1]
		CFG[0] = strings.Replace(CFG[0], "$HOME", ENV["HOME"], -1)
		CFG[1] = strings.Replace(CFG[1], "$HOME", ENV["HOME"], -1)
	} else {
		// $HOME is empty, we unset CFG[0] and CFG[1]
		CFG[0] = ""
		CFG[1] = ""
	}

	// Load environment variables from all the configuration files
	// specified in the slice 'CFG'
	for _, filepath := range CFG {
		if filepath != "" {
			ret, err := Exists(filepath)
			check(err)
			//fmt.Printf("checking: %s (%t) \n", filepath, ret)

			// if the conf file exists, load it with godotenv
			if ret {
				err := godotenv.Load(filepath)
				check(err)
			}
		}
	}

	// Populate ENV map with environment variables
	for k, _ := range ENV {
		if k == "HOME" {
			continue
		}
		ENV[k] = os.Getenv(k)
	}
	//fmt.Println("map:", ENV)

	// Filter ENV map to expand bash variables
	for k, v := range ENV {
		if k == "HOME" {
			continue
		}
		// expand $HOME var
		v = strings.Replace(v, "$HOME", ENV["HOME"], -1)
		v = strings.Replace(v, "${HOME}", ENV["HOME"], -1)

		// Expand $TODO_DIR var
		v = strings.Replace(v, "$TODO_DIR", ENV["TODO_DIR"], -1)
		v = strings.Replace(v, "${TODO_DIR}", ENV["TODO_DIR"], -1)

		// save the filtered value
		ENV[k] = v
	}

	// Initialize the app CLI
	app := cli.NewApp()
	app.Name = "todo"
	app.Usage = "A simple and extensible utility for managing your todo.txt files"
	app.Version = "1.0.0"
	app.Author = "Mauro Toffanin"
	app.Email = "toffanin.mauro@gmail.com"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:        "env",
			Usage:       "Prints `todo` environment information.",
			Description: "By default env prints information as a shell script.\n   The environment info will be dumped in straight-forward\n   form suitable for sourcing into a shell script.\n\n   If one or more variable names is given as arguments, env\n   prints the value of each named variable on its own line.",
			Action: func(c *cli.Context) {
				// collect all the user-submitted arguments in an array
				// and store its lenght for later usage
				args := c.Args()
				nargs := len(args)

				// print only the required environment variables
				if 0 < nargs {
					for _, arg := range args {
						fmt.Printf("%s=\"%s\"\n", arg, ENV[arg])
					}
					// return since there is nothing else to do
					// (this reduce code size as the go compiler can optimize the branching)
					return
				}

				// print all the environment variables
				if 0 == nargs {
					for k, v := range ENV {
						fmt.Printf("%s=\"%s\"\n", k, v)
					}
				}
			},
		},
		{
			Name:        "add",
			ShortName:   "a",
			Usage:       "Adds a task to your todo.txt file.",
			Description: "add \"feed the cat\"\n   Adds \"feed the cat\" to your todo.txt file on its own line.",
			Action: func(c *cli.Context) {
				// TODO: validating input as a task

				// TODO: task mangler
				task := fmt.Sprintf("%s", c.Args().First())

				// Add the new task
				addTask(task)
			},
		},
		/*{
			Name:  "complete",
			Usage: "Completes a task",
			Action: func(c *cli.Context) {
				fmt.Println("completed task: ", c.Args().First())
			},
		},*/
	}
	app.Run(os.Args)
}
