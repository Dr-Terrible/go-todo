# go-gtd

[![Build Status](https://travis-ci.org/toffanin/go-gtd.svg?branch=master)](https://travis-ci.org/toffanin/go-gtd) [![GoDoc](https://godoc.org/github.com/toffanin/go-gtd?status.png)](https://godoc.org/github.com/toffanin/go-gtd) 

`Go - Getting Things Done` provides a more advanced drop in replacement for the bash version of [Todo.txt CLI](https://github.com/ginatrapani/todo.txt-cli) and is meant to handle the following features:

- parsing and manipulating of task lists;
- stores tasks hierarchically, with each task given one of five priority levels;

## Purpose

This package is meant to be a [Golang](http://golang.org) implementation of Gina Trapani's [todo.txt-cli](https://github.com/ginatrapani/todo.txt-cli) which doesn't rely on other external dependencies (such as bash, grep or sed) and can integrate into 3rd party systems and APIs.


## Requirements

This command-line tool requires Go 1.2 ( or higher) and you also need to have `git` installed to build this utility.

Once you have Go installed and your ``GOPATH`` set, do the following to build `Go-GTD`:

```
git clone https://github.com/toffanin/go-gtd.git
cd go-gtd
go build todo.go
```

If you just want to run the code without building it, you can call:

```
go build todo.go
```


# TODO

- [ ] full compatibility with the Todo.txt CLI sintax
  - [ ] add|a
  - [ ] addm
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
- [ ] extra commands not part of the original CLI sintax
  - [x] env - prints `todo` environment information
  - [ ] init - create a configuration file with default values
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
