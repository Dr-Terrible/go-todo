package commands

import (
	"github.com/codegangsta/cli"
)

func GetList() cli.Command {

	return cli.Command{
		Name:        "list",
		ShortName:   "ls",
		Usage:       "Displays all tasks that contain TERM sorted by priority with line numbers.",
		Description: "",
		Action: func(c *cli.Context) {
			// collect all the user-submitted arguments in an array
			//args := c.Args()

			// TODO: read todo.txt

			// TODO: calculate padding

			// TODO: Number the file, then run the filter command,
			// then sort and mangle output some more

			// TODO: Build and apply the filter

			// print output
			/*tasklist := TaskList{}
			file, err := os.Open(filename)*/
		},
	}
}
