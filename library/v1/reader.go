// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package todotxt reads and writes todo.txt files as defined in Gina Trapani's
// Todo.txt Format: https://github.com/ginatrapani/todo.txt-cli/wiki/The-Todo.txt-Format
//
// A todo.txt file contains zero or more tasks.
// A single line in your todo.txt text file represents a single task.
package todotxt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	//	"unicode"
	"github.com/toffanin/go-todo/utils"
)

// Task represents a todo.txt task entry
type Task struct {
	Id             uint64 // Internal task ID
	Raw            string // Raw task text
	Todo           string // Todo part of task text
	Priority       string
	Projects       []string
	Contexts       []string
	AdditionalTags map[string]string // Add-on tags will be available here.
	CreatedDate    time.Time
	DueDate        time.Time
	CompletedDate  time.Time
	Completed      bool
	Padding        uint
}

// TaskList represents a list of todo.txt task entries.
// It is usually loaded from a whole todo.txt file.
type TaskList []Task

// A ParseError is returned for parsing errors.
// The first line is 1.  The first column is 0.
type ParseError struct {
	Line   int   // Line where the error occurred
	Column int   // Column (rune index) where the error occurred
	Err    error // The actual error
}

// A Reader reads tasks from a todo.txt file.
//
// As returned by NewReader, a Reader expects input conforming to Todo.txt Format.
// The exported fields can be changed to customize the details before the first
// call to Read or ReadAll.
//
// Comment, if not 0, is the comment character. Lines beginning with the
// Comment character are ignored. It defaults to '#'.
type Reader struct {
	Comment rune          // character used for comments
	line    uint64        // holds the total number of lines parsed
	column  uint          // holds the scanner position for a token
	length  uint64        // holds the total number of tasks
	buffer  *bufio.Reader // buffer used for parsing and scanning io inputs
	tasks   TaskList
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	// Each task in todo.txt starts on a new line.
	// We increment our line number (lines start at 1) and set column to 0 so as
	// we increment in Read()/ReadAll() it points to the token we read.

	return &Reader{
		Comment: '#',
		line:    1,
		column:  0,
		length:  1,
		buffer:  bufio.NewReader(r),
		tasks:   TaskList{},
	}
}

// Read reads one task from r.
func (r *Reader) Read() (*TaskList, error) {

	rawTask, err := r.buffer.ReadString('\n')
	if err == io.EOF {
		return &r.tasks, err
	}
	utils.Check(err)

	// Set the split function for a Scanner that returns each line of text,
	// stripped of any trailing end-of-line marker
	scanner := bufio.NewScanner(strings.NewReader(rawTask))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		rawTask = scanner.Text()
		// skip blank lines and comments
		if rawTask == "" || (r.Comment != 0 && strings.HasPrefix(rawTask, "#")) {
			fmt.Println("****")
			break
		}
		//fmt.Printf("task: %s (test)\n", rawTask)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	return &r.tasks, nil
}

// Read reads all the tasks from r.
func (r *Reader) ReadAll() (TaskList, error) {

	scanner := bufio.NewScanner(r.buffer)
	for scanner.Scan() {
		// strip any leading and trailing white spaces
		rawTask := strings.TrimSpace(scanner.Text())

		// strip any trailing end-of-line markers
		scanner := bufio.NewScanner(strings.NewReader(rawTask))
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			rawTask = scanner.Text()
		}

		// skip blank lines
		if rawTask == "" {
			continue
		}

		// decode todo.txt cli format
		task, err := r.parseRecord(rawTask, r.length)
		if err != nil {
			return r.tasks, err
		}

		// append the new mangled task into the structured tasks list
		r.tasks = append(r.tasks, *task)
		//fmt.Printf("task (mangled): %st\n", task)

		r.length++
		r.line++
	}
	if err := scanner.Err(); err != nil {
		return r.tasks, err
	}

	return r.tasks, nil
}

// parseRecord reads and parses a single todo.txt task from r.
func (r *Reader) parseRecord(raw string, id uint64) (*Task, error) {

	var err error
	task := Task{}

	task.Raw = raw
	task.Todo = raw
	task.Id = id

	// TODO: check for data completed date

	// TODO: check for priority

	// TODO: check for created date

	// TODO: check for contexts and projects
	// Set the split function for a Scanner that returns each token inside the
	// line of text previously scanned
	scanner := bufio.NewScanner(strings.NewReader(raw))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		token := scanner.Text()

		if strings.IndexRune(token, '@') == 0 {
			task.Contexts = append(task.Contexts, token)
		}
		if strings.IndexRune(token, '+') == 0 {
			task.Projects = append(task.Projects, token)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	//fmt.Printf("task: %s (contexts: %d) (projects: %d)\n", raw, contexts, projects)

	// TODO: check for additional tags

	// trim any remaining white spaces
	task.Todo = strings.TrimSpace(task.Todo)
	//fmt.Println("task: ", raw)

	return &task, err
}

//
func (r *Reader) Len() uint64 {
	return r.length - 1
}
