# How to contribute

I really love pull requests and third-party patches are essential for keeping
`go-todo` in shape, but I simply can not access the huge number of platforms and
myriad configurations for running the tool. Therefore I want to keep it as easy
as possible to contribute changes that get things working in your environment.

Here below there are a few guidelines that I need contributors to follow so that
I can have a chance of keeping on top of things.


## Getting Started

1. Fork the repository on GitHub.
2. Create a topic branch from where you want to base your work.
3. Make sure you have added the necessary tests for your changes and make the
   tests pass. Only re-factoring and documentation changes require no new tests.
   If you are adding functionality or fixing a bug, I need a test.
4. Run _all_ the tests to assure nothing else was accidentally broken. **I only
   take pull requests with passing tests**.
5. Make sure you run `go fmt`, `go vet` and `golint` before to commit your changes.
6. Make sure your commit messages are in the proper format.
7. Push to your fork and submit a [pull request](http://help.github.com/send-pull-requests/).

At this point you are waiting for my feedbacks. I look at pull requests within
few days. I may suggest some improvements or alternatives.

Some things that will increase the chance that your pull request is accepted:

* Use Go idioms and helpers
* Include tests that fail without your code, and pass with it
* Update the documentation, the surrounding one, examples elsewhere, guides, and
  whatever is affected by your contribution

And in case I didn't emphasize it enough: I really ♥ tests!