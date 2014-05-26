package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"../utils"

	"github.com/codegangsta/cli"
)

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
		utils.Check(err)

		ret, _ := utils.Exists(filePath)
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
	utils.Check(err)

	// print first line of the action's summary
	fmt.Printf("%s todo.txt structure in %s\n", message, filePath)

	// create the missing files of the todo.txt structure
	for k, filename := range FileName {
		// sanitize the absolute path of the file
		filePath, err := filepath.Abs(destination + FileName[k])
		utils.Check(err)
		//fmt.Printf("absolute path: %s\n", filePath)

		// Open file
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		defer file.Close()

		// if there aren't errors, write a new file with default values
		if err == nil {
			size, err := file.WriteString(FileTemplate[k])
			utils.Check(err)

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

func GetInitCommand() cli.Command {

	return cli.Command{
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
}
