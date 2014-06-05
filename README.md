# go-todo

[![Build Status](https://travis-ci.org/toffanin/go-todo.svg?branch=master)](https://travis-ci.org/toffanin/go-todo) [![GoDoc](https://godoc.org/github.com/toffanin/go-todo?status.png)](https://godoc.org/github.com/toffanin/go-todo)

`go-todo` provides a more advanced drop in replacement for the bash version of [Todo.txt CLI](https://github.com/ginatrapani/todo.txt-cli) and is meant to handle the following features:

- parsing and manipulating of task lists;
- stores tasks hierarchically, with each task given one of five priority levels;

## Purpose

This package is meant to be a [Golang](http://golang.org) implementation of Gina Trapani's [todo.txt-cli](https://github.com/ginatrapani/todo.txt-cli) which doesn't rely on other external dependencies (such as bash, cat, grep, awk, sort and sed) and can integrate into 3rd party systems and APIs.


## Requirements

This command-line tool requires Go 1.2 ( or higher) and you also need to have `git` installed to build this utility.

Once you have Go installed and your ``GOPATH`` set, do the following to build `go-todo`:

```
git clone https://github.com/toffanin/go-todo.git
go install todo.go
```

If you just want to run the code without building it, you can call:

```
go run todo.go
```


# TODO

- [ ] full compatibility with the Todo.txt CLI sintax
  - [x] add|a
  - [x] addm
  - [ ] addto
  - [ ] append|app
  - [ ] archive
  - [ ] command
  - [ ] deduplicate
  - [ ] del|rm
  - [ ] depri|dp
  - [ ] do
  - [x] help
  - [ ] list|ls
    - [ ] TERMS
    - [ ] logical operators
    - [x] TODOTXT_VERBOSE
  - [ ] listall|lsa
  - [ ] listaddons
  - [ ] listcon|lsc
  - [ ] listfile|lf
  - [ ] listpri|lsp
  - [ ] listproj|lsprj
  - [ ] move|mv
  - [ ] prepend|prep
  - [ ] pri|p
  - [ ] replace
  - [ ] resort
  - [x] shorthelp
  - [ ] -@ | -@@
  - [ ] -+ | -++
  - [ ] -c
  - [ ] -d | TODOTXT_CFG_FILE
  - [ ] -f | TODOTXT_FORCE
  - [ ] -h
  - [ ] -p | -P | TODOTXT_PLAIN
  - [ ] -a | -A | TODOTXT_AUTO_ARCHIVE
  - [ ] -n | -N | TODOTXT_PRESERVE_LINE_NUMBERS
  - [x] -t | -T | TODOTXT_DATE_ON_ADD
  - [ ] -v | -vv | TODOTXT_VERBOSE
  - [ ] -V
  - [ ] -x | TODOTXT_DISABLE_FILTER
- [ ] extra commands not part of the original CLI sintax
  - [x] env - prints `go-todo` environment information
  - [ ] init - create a configuration file with default values
  - [ ] status - can be used to obtain a status summary
- [ ] full compatibility with the [Todo.txt Format](https://github.com/ginatrapani/todo.txt-cli/wiki/The-Todo.txt-Format)
  - [ ] filters (completed tasks are hidden by default, but may be displayed with -A)
- [ ] full compatibility with the [Todo.txt Add-ons](https://github.com/ginatrapani/todo.txt-cli/wiki/Creating-and-Installing-Add-ons)
- [ ] stores tasks hierarchically;
- [ ] integrate with third party systems
- [ ] integrate with third party APIs
- [ ] readline-based editing of task text and priority
- [ ] linked files
- [x] todo.cfg configuration file
- [ ] colour customisation
- [ ] custom task formatting