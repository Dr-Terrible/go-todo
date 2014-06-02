// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

var (

	/*
	 * This map is an empty configuration representation used by todo.txt CLI.
	 * This representation can then be filled with settings from
	 * environment variables.
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
	 * The slice defines also the exact order at which the configuration files
	 * should be loaded by the godotenv library.
	 */
	// TODO: transform the slice into a map[string]bool where the boolean value defines if the config file exists
	cfgPath = []string{
		"$HOME/todo.cfg",
		"$HOME/.todo.cfg",
		"./todo.cfg",
		"/etc/todo/config"}
)

// SetSetting adds a setting and a value to the configuration.
// It returns true if the setting and value were inserted.
func SetSetting(name string, value string) bool {
	if name == "" {
		return false
	}
	settings[name] = value
	return HasSetting(name)
}

// GetSetting retrieves the value for the given setting.
func GetSetting(name string) string {
	if name == "" || !HasSetting(name) {
		return ""
	}

	return settings[name]
}

// HasSettings checks if the configuration has the given setting.
// It returns false if the setting does not exist.
func HasSetting(name string) bool {
	_, exist := settings[name]
	return exist
}

// Looks up the value of a setting, returns false if no bool value exists.
func IsSettingBool(name string) bool {
	switch settings[name] {
	case "0":
	case "1":
		return true
	}
	return false
}

// GetSettings returns a list of all the settings
func GetSettings() map[string]string {
	return settings
}

// LoadConfig reads all the configuration files (todo.cfg) and then
// creates a configuration representation filled with keys and values.
//
// Call this function as close as possible to the start of your
// application, ideally in main().
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
	if env["HOME"] != "" {
		// $HOME isn't empty, we expand $HOME for cfgPath[0] and cfgPath[1]
		cfgPath[0] = strings.Replace(cfgPath[0], "$HOME", env["HOME"], -1)
		cfgPath[1] = strings.Replace(cfgPath[1], "$HOME", env["HOME"], -1)
	} else {
		// $HOME is empty, we unset cfgPath[0] and cfgPath[1]
		cfgPath[0] = ""
		cfgPath[1] = ""
	}

	// Load environment variables from all the configuration files
	// specified in the slice 'cfgPath'
	for _, filepath := range cfgPath {
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
}
