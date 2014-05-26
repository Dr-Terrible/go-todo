package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

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

func GetHome() string {
	return ENV["HOME"]
}

func GetPwd() string {
	return ENV["PWD"]
}

func SetTodoEnv(key string, value string) error {
	return nil
}
func GetTodoEnv(key string) string {
	if key == "" {
		return ""
	}
	return TODOENV[key]
}

func LoadTodoEnv() {
	// Retrieve environment variables $HOME and $PWD
	pwd, err := os.Getwd()
	Check(err)
	ENV["PWD"] = path.Clean(pwd)
	ENV["HOME"] = path.Clean(os.Getenv("HOME"))

	/*
	 * The original bash script Todo.txt CLI relies on the expansion facilities
	 * that are built-in into the shell and performed automatically when the
	 * script is invoked. Go lang doesn't perform such facilities, that means
	 * that certain Environment Variables (like $HOME) aren't expanded correctly.
	 *
	 * ${var} and $var expansion is mimicked to guarantee backward-compatibility
	 * with the original bash script.
	 */
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
			Check(err)
			//fmt.Printf("checking: %s (%t) \n", filepath, ret)

			// if the conf file exists, load it with godotenv
			if ret {
				err := godotenv.Load(filepath)
				Check(err)
			}
		}
	}

	// Populate TODOENV map with environment variables
	for k := range TODOENV {
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

/*var (
	InputFile   *os.File = os.Stdin
	inputBuffer *bufio.Reader
)*/

func InteractiveInput() string {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	Check(err)

	// replace return carriage chars with spaces
	// and trim leading / ending spaces
	input = strings.NewReplacer("\n", " ", "\t", " ", "\r", " ").Replace(input)
	return strings.TrimSpace(input)
	//return strings.TrimSpace(string(line))
}

/*func buffer() *bufio.Reader {
	if inputBuffer == nil {
		inputBuffer = bufio.NewReader(InputFile)
	}
	return inputBuffer
}*/

// TODO: this is an ugly and hackish error handler that needs to be improved
func Check(e error) {
	if e != nil {
		fmt.Errorf("%v", e)
		panic(e)
	}
}
