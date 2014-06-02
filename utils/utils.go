// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package utils provides helper functions for applications that need to
// interact with Todo.txt CLI files.
//
// A complete Todo.txt configuration representation can be retrieved as follows:
//
//   package main
//
//   import (
//     "fmt"
//     "github.com/toffanin/go-todo/utils"
//   )
//
//   func main() {
//     // Load Todo.txt CLI settings from todo.cfg
//     utils.LoadConfig()
//
//     settings := utils.GetSettings()
//     fmt.Printf("Settings: %#v\n", settings)
//
//     // now continue with your app and do something useful with settings
//   }
//
// Single configuration values can be retrieved as follow:
//
//  func main() {
//     todoDir := utils.GetSetting("TODO_DIR")
//     todoActionsDir := utils.GetSetting("TODO_ACTIONS_DIR")
//
//     // do something with todoDir and todoActionsDir
//  }
//
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
)

// GetHome retrieves the value of the environment variable named $HOME.
func GetHome() string {
	return env["HOME"]
}

// GetPwd retrieves the path name corresponding to the current directory.
func GetPwd() string {
	return env["PWD"]
}

// InteractiveInput shows a prompt and then reads a String provided by a user at
// a command-line.
func InteractiveInput(prompt string) string {
	if prompt != "" {
		fmt.Printf("%s ", prompt)
	}
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	Check(err)

	// sanitize input
	return SanitizeInput(input)
}

// SanitizeInput applies the following rules iteratively until no further
// processing can be done:
//
// - trim all the extra white spaces
// - trim all return carriage chars
// - trim leading / ending quotation marks
// - trim leading / ending spaces
//
func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return input
	}
	input = strings.TrimPrefix(input, "\"")
	input = strings.TrimSuffix(input, "\"")
	return strings.NewReplacer("  ", " ", "\n", " ", "\t", " ", "\r", " ").Replace(input)
}
