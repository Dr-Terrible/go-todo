package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
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
		fmt.Errorf("%v", e)
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

	// determine the number of tasks in todo.txt
	// TODO: with NewReadWriter the code should be more compact
	//       buf := bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))
	fd, err := os.OpenFile(TODOENV["TODO_FILE"], os.O_RDONLY|os.O_CREATE, 0600)
	check(err)
	scanner := bufio.NewScanner(fd)
	ntasks := 1
	for scanner.Scan() {
		ntasks++
	}
	if err := scanner.Err(); err != nil {
		fd.Close()
		check(err)
	}
	//fmt.Printf("n. lines: %d\n", ntasks)
	err = fd.Close()
	check(err)

	// Open todo.txt in append mode only
	fd, err = os.OpenFile(TODOENV["TODO_FILE"], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	check(err)
	defer fd.Close()

	// use buffered I/O
	writer := bufio.NewWriter(fd)

	// add the task to todo.txt
	_, err = writer.WriteString(task + "\n")
	check(err)
	err = writer.Flush()
	check(err)

	// print summary
	fmt.Printf("%d: %s\n", ntasks, task)
	fmt.Printf("TODO: %d added\n", ntasks)
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

	envCommand := cli.Command{
		Name:        "env",
		Usage:       "Prints `todo` environment information.",
		Description: "By default env prints information as a shell script.\n   The environment info will be dumped in straight-forward\n   form suitable for sourcing into a shell script.\n\n   If one or more variable names is given as arguments, env\n   prints the value of each named variable on its own line.",
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
					fmt.Printf("%s=\"%s\"\n", arg, TODOENV[arg])
				}
			case false:
				// print all the environment variables
				for k, v := range TODOENV {
					fmt.Printf("%s=\"%s\"\n", k, v)
				}
			}
		},
	}

	initCommand := cli.Command{
		Name:  "init",
		Usage: "Initialize a new todo.txt structure with default values",
		Description: `
   This command creates all the template file required by the todo.txt and
   a configuration files with default values - basically, the values TODO_DIR,
   TODO_FILE, DONE_FILE, REPORT_FILE and TODO_ACTIONS_DIR are exported.

   If the option '--dest' is set then it specifies a path to use instead of
   the working directory as the destination path for the new structure.

   Running 'todo init' in a pre-initialized directory is safe; it will not
   overwrite things that are already there.`,
		Flags: []cli.Flag{
			cli.StringFlag{"dest, d", "/path/to/your/dir", "specifies a different destination path"},
		},
		Action: func(c *cli.Context) {
			destination := "."
			if c.IsSet("dest") {
				//fmt.Println("dest:", c.String("dest"))
				destination = c.String("dest")
			}
			initAction(destination)
		},
	}

	shorthelpCommand := cli.Command{
		Name:        "shorthelp",
		Usage:       "",
		Description: "",
		Action: func(c *cli.Context) {
			cli.ShowAppHelp(c)
		},
	}

	addCommand := cli.Command{
		Name:      "add",
		ShortName: "a",
		Usage:     "Adds a task to your todo.txt file.",
		Description: `
   This command can be used to add the specified task to your todo.txt file on
   its own line.

   Project and content notation is optional. Quotes are optional too.

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
			fmt.Printf("(add::Action) global flag: %s\n", c.GlobalString("t"))*/

			// check incorrect usage of the command
			switch {
			case len(args) == 0: // no options specified
				fmt.Print("\nIncorrect Usage: missing option with command \"add [task]\"\n\n")
				cli.ShowCommandHelp(c, "add")
				return
			}

			// TODO: validating input as a task

			/* task mangler
			 * - collect all the arguments into a single string
			 * - replace return carriage chars with spaces
			 * - trim the task from leading / ending spaces
			 * - honor the -t/-T global flag if present
			 */
			task := strings.Join(args[0:], " ")
			r := strings.NewReplacer("\n", " ", "\t", " ", "\r", " ")
			task = r.Replace(task)
			task = strings.TrimSpace(task)

			if c.GlobalBool("t") {
				date := time.Now().Format("2006-01-02 ")
				//fmt.Printf("(add::Action) date: %s\n", date)
				task = date + task
			}

			// save the new task
			//fmt.Printf("task: |%s|\n", task)
			addAction(task)
		},
	}

	// Initialize the app CLI
	app := cli.NewApp()
	app.Name = "go-todo"
	app.Usage = "A simple and extensible utility for managing your todo.txt files"
	app.Version = "1.0.0"
	app.Author = "Mauro Toffanin"
	app.Email = "toffanin.mauro@gmail.com"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{"t", "", "Prepend the current date to a task automatically when it's added"},
		cli.StringFlag{"T", "", "Do not prepend the current date to a task automatically when it's added."},
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()

		// debugging
		/*fmt.Printf("(app) args (%d): %s\n", len(args), args)
		fmt.Printf("(app) global flag: %s\n", c.GlobalString("t"))*/

		switch c.GlobalString("t") {
		case "add":
			// inject all the arguments into a new flag set
			set := flag.NewFlagSet("add", 0)
			set.Parse([]string{"add", strings.Join(args[0:], " ")})

			// preserve the Global Flag by pushing the flag name as its value
			// (it's a hack, but it works)
			gset := flag.NewFlagSet("add", 0)
			gset.Bool("t", true, "")

			// create a new Context and run the relative Command
			c := cli.NewContext(app, set, gset)
			err := addCommand.Run(c)
			check(err)
			return
		default:
			// TODO: print error about misuse of the option -t missing 'add'
		}
		cli.ShowAppHelp(c)

	}
	app.Commands = []cli.Command{
		envCommand,
		initCommand,
		shorthelpCommand,
		addCommand,
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
	err = app.Run(os.Args)
	check(err)
}
