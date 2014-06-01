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
	env = map[string]string{
		"PWD":  "",
		"HOME": "",
	}

	/*
	 * This map stores all the environment variables used by the todo.txt CLI
	 */
	settings = map[string]string{
		"TODO_DIR":             "",
		"TODO_FILE":            "",
		"DONE_FILE":            "",
		"REPORT_FILE":          "",
		"TODO_ACTIONS_DIR":     "",
		"TODOTXT_SORT_COMMAND": "",
		"TODOTXT_FINAL_FILTER": "",
		"TODOTXT_DATE_ON_ADD":  "0",
		"TODOTXT_FORCE":        "0",
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
	ConfPaths = []string{
		"$HOME/todo.cfg",
		"$HOME/.todo.cfg",
		"./todo.cfg",
		"/etc/todo/config"}
)

func GetHome() string {
	return env["HOME"]
}

func GetPwd() string {
	return env["PWD"]
}

func SetSetting(name string, value string) error {
	return nil
}
func GetSetting(name string) string {
	if name == "" {
		return ""
	}
	return settings[name]
}

// TODO: Determines if the setting was actually set in todo.cfg
func HasSetting(name string) bool {
	return true
}

// TODO: Looks up the value of a setting, returns false if no bool value exists
func IsSettingBool(name string) bool {
	return true
}

// TODO: Returns the all the settings
func GetSettings() map[string]string {
	return settings
}

func LoadConfig() {
	// Retrieve environment variables $HOME and $PWD
	pwd, err := os.Getwd()
	Check(err)
	env["PWD"] = path.Clean(pwd)
	env["HOME"] = path.Clean(os.Getenv("HOME"))

	/*
	 * The original bash script Todo.txt CLI relies on the expansion facilities
	 * that are built-in into the shell and performed automatically when the
	 * script is invoked. Go lang doesn't perform such facilities, that means
	 * that certain Environment Variables (like $HOME) aren't expanded correctly.
	 *
	 * ${var} and $var expansion is mimicked to guarantee backward-compatibility
	 * with the original bash script.
	 */
	//fmt.Printf("ConfPaths (raw)     : %s\n", ConfPaths)
	if env["HOME"] != "" {
		// $HOME isn't empty, we expand $HOME for ConfPaths[0] and ConfPaths[1]
		ConfPaths[0] = strings.Replace(ConfPaths[0], "$HOME", env["HOME"], -1)
		ConfPaths[1] = strings.Replace(ConfPaths[1], "$HOME", env["HOME"], -1)
	} else {
		// $HOME is empty, we unset ConfPaths[0] and ConfPaths[1]
		ConfPaths[0] = ""
		ConfPaths[1] = ""
	}
	//fmt.Printf("ConfPaths (filtered): %s\n", ConfPaths)

	// Load environment variables from all the configuration files
	// specified in the slice 'ConfPaths'
	for _, filepath := range ConfPaths {
		if filepath != "" {
			ret, err := Exists(filepath)
			//fmt.Printf("checking: %s (%t) \n", filepath, ret)

			// if the conf file exists, load it with godotenv
			if ret {
				err := godotenv.Load(filepath)
				Check(err)
			}

			Check(err)
		}
	}

	// Populate settings map with environment variables
	for k := range settings {
		settings[k] = os.Getenv(k)
	}
	//fmt.Println("settings: (raw)     : ", settings)

	// Sanitize settings map by expanding $HOME bash variables
	for k, v := range settings {
		v = strings.Replace(v, "$HOME", env["HOME"], -1)
		v = strings.Replace(v, "${HOME}", env["HOME"], -1)

		// save the sanitized value
		settings[k] = v
	}
	// Sanitize settings map by expanding $TODO_DIR bash variables
	for k, v := range settings {
		v = strings.Replace(v, "$TODO_DIR", settings["TODO_DIR"], -1)
		v = strings.Replace(v, "${TODO_DIR}", settings["TODO_DIR"], -1)

		// save the sanitized value
		settings[k] = v
	}
	//fmt.Println("settings: (filtered): ", settings)
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

func InteractiveInput(prompt string) string {
	if prompt != "" {
		fmt.Printf("%s ", prompt)
	}
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	Check(err)

	// sanitize
	return SanitizeInput(input)
}

/* CleanInput applies the following rules iteratively until no further
 * processing can be done:
 *
 * - trim all the extra white spaces
 * - trim all return carriage chars
 * - trim leading / ending quotation marks
 * - trim leading / ending spaces
 */
func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return input
	}
	input = strings.TrimPrefix(input, "\"")
	input = strings.TrimSuffix(input, "\"")
	return strings.NewReplacer("  ", " ", "\n", " ", "\t", " ", "\r", " ").Replace(input)
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
