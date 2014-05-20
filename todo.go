package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strings"
)

var (
	/*
	 * This map stores some useful global environment variables used
	 * internally by the logic of the tool (we avoid to clutter the source
	 * code with redundant calls to pkg/os functions).
	 */
	ENV = map[string]string{
		"PWD":  "",
		"HOME": "",
	}

	/*
	 * This map stores all the environment variables used by the todo.txt CLI
	 */
	TODOENV = map[string]string{
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
func addAction(task string) {
	// TODO: TODOENV["TODO_FILE"] path should be validated somehow
	// before to be stated by os.OpenFile
	//path.Clean(TODOENV["TODO_FILE"])

	// Open todo.txt in append mode only
	todoFile, err := os.OpenFile(TODOENV["TODO_FILE"], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	defer todoFile.Close()
	check(err)

	// add the task to the todo.txt file
	ret, err := todoFile.WriteString(task + "\n")
	check(err)
	fmt.Printf("added new task (%d bytes): \"%s\"\n", ret, task)

	// sync / flush todo.txt
	todoFile.Sync()
}

// Create a todo.txt structure at the specified location (default destination is ".")
func initAction(destination string) {

	var (
		FileName = map[string]string{
			"cfg":    "/todo.cfg",
			"todo":   "/todo.txt",
			"done":   "/done.txt",
			"report": "/report.txt",
		}

		FileTemplate = map[string]string{
			"cfg": `
# === EDIT FILE LOCATIONS BELOW ===

# Your todo.txt directory
#export TODO_DIR="$HOME/todo"
export TODO_DIR="."

# Your todo/done/report.txt locations
export TODO_FILE="$TODO_DIR/todo.txt"
export DONE_FILE="$TODO_DIR/done.txt"
export REPORT_FILE="$TODO_DIR/report.txt"

# You can customize your actions directory location
#export TODO_ACTIONS_DIR="$HOME/.todo.actions.d"`,
			"todo":   "",
			"done":   "",
			"report": "",
		}

		initiated = false
		message   = "Initialized a new"
	)

	// try to guess if the destination is an existing and
	// pre-configured todo.txt structure.
	for _, filename := range FileName {
		// sanitize the absolute path of the file
		filePath, err := filepath.Abs(destination + filename)
		//fmt.Printf("absolute path: %s\n", cfgFilePath)
		check(err)

		ret, _ := Exists(filePath)
		if ret {
			initiated = true
			break
		}
	}

	// if the destination is an existing todo.txt structure then
	// change the message accordingly
	if initiated {
		message = "Reinitialized an existing"
	}

	// sanitize the absolute path of the destination
	filePath, err := filepath.Abs(destination)
	check(err)

	// print first line of the action's summary
	fmt.Printf("%s todo.txt structure in %s\n", message, filePath)

	// create the missing files of the todo.txt structure
	for k, filename := range FileName {
		// sanitize the absolute path of the file
		filePath, err := filepath.Abs(destination + FileName[k])
		check(err)
		//fmt.Printf("absolute path: %s\n", filePath)

		// Open file
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		defer file.Close()

		// if there aren't errors, write a new file with default values
		if err == nil {
			size, err := file.WriteString(FileTemplate[k])
			check(err)

			// sync / flush file
			file.Sync()

			// print a small summary
			fmt.Printf("%s [%s] (%d bytes)\n", filename, "new", size)
			continue
		}

		// file exists, there is nothing to write
		if os.IsExist(err) {
			// print a small summary
			fmt.Printf("%s [%s]\n", filename, "exists")
			continue
		}
	}
}

func main() {

	// Retrieve environment variables $HOME and $PWD
	pwd, err := os.Getwd()
	check(err)
	ENV["HOME"] = os.Getenv("HOME")
	ENV["PWD"] = pwd

	// Expand $HOME var inside CFG slice
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

	// Populate TODOENV map with environment variables
	for k, _ := range TODOENV {
		TODOENV[k] = os.Getenv(k)
	}
	//fmt.Println("map:", TODOENV)

	// Sanitize TODOENV map by expanding bash variables
	for k, v := range TODOENV {
		// expand $HOME var
		v = strings.Replace(v, "$HOME", ENV["HOME"], -1)
		v = strings.Replace(v, "${HOME}", ENV["HOME"], -1)

		// Expand $TODO_DIR var
		v = strings.Replace(v, "$TODO_DIR", TODOENV["TODO_DIR"], -1)
		v = strings.Replace(v, "${TODO_DIR}", TODOENV["TODO_DIR"], -1)

		// save the sanitized value
		TODOENV[k] = v
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

				// debugging
				/*fmt.Printf("pwd: %s\n", pwd)
				fmt.Printf("HOME=\"%s\"\n", ENV["HOME"])*/

				// print only the required environment variables
				if 0 < nargs {
					for _, arg := range args {
						fmt.Printf("%s=\"%s\"\n", arg, TODOENV[arg])
					}
					// return since there is nothing else to do
					// (this reduce code size as the go compiler can optimize the branching)
					return
				}

				// print all the environment variables
				if 0 == nargs {
					for k, v := range TODOENV {
						fmt.Printf("%s=\"%s\"\n", k, v)
					}
				}
			},
		},
		{
			Name:        "init",
			Usage:       "Create a configuration file with default values",
			Description: "This command creates a configuration file with default values - basically\n   a TODO_DIR, TODO_FILE, DONE_FILE, REPORT_FILE and TODO_ACTIONS_DIR.\n\n   If the option `-d` is set then it specifies a path to use instead of\n   ./todo.cfg as the destination path for the configuration file\n\n   Running `todo init` in a pre-initialized directory is safe; it will not\n   overwrite things that are already there.",
			Action: func(c *cli.Context) {
				initAction(".")
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
				addAction(task)
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
