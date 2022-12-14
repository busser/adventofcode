# Advent of Code

![tests-passing](https://github.com/busser/adventofcode/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/busser/adventofcode)](https://goreportcard.com/report/github.com/busser/adventofcode)

[Advent of Code](https://adventofcode.com) is an Advent calendar of small
programming puzzles for a variety of skill sets and skill levels that can be
solved in any programming language you like.

In the Cloud Native ecosystem, [Go](https://golang.org/) has emerged as the
language of choice. As a Site Reliability Engineer, knowing how to write
software in Go is a useful skill.

Whether you are new to Go or have already written countless lines of it, this
repository aims at giving you everything you need to improve.

## Getting started

First and foremost, go to the [Advent of Code website](https://adventofcode.com/)
and log in.

To get started on a solution, clone this repository:

```bash
git clone git@github.com:busser/adventofcode.git
cd adventofcode
```

Remove existing solutions to start from scratch:

```bash
rm -r y*
```

Build the `adventofcode` command-line tool:

```bash
make build
```

Get started on your first puzzle:

```bash
bin/adventofcode scaffold --day 1 --workdir "$(pwd)"
```

The command above will create a package where your code will go. Your next steps
should be:

1. Implement your solution in the `solution.go` file. Use this command to test
   it:

   ```bash
   go test ./y2021/d01 -run ExamplePartOne
   ```

2. Once you think you have found the answer to the problem, submit it on the
   adventofcode.com website. If it's the right answer, congrats!

3. Update your tests by adding the answer to `ExamplePartOne` in
   `solution_test.go`.
4. Repeat steps 1 to 3 for the second part of the Advent of Code problem.
5. Now that you have finished, run all tests to make sure everything is ready
   for your pull request:

   ```bash
   make test
   ```

Once you are done with the first day of the Advent of Code, have a look at the
[Configuration](#configuration) section below on how you can configure the
`adventofcode` CLI, or the [Tests and benchmarks](#tests-and-benchmarks) section
for tips on how to measure your solution's performance.

## Scaffolding

The `adventofcode` can help you get started quickly on your solution to the
daily Advent of Code problems. The `scaffold` subcommand builds the following
for you:

- A new package to write your solution in;
- A `solution.go` file with a basic code skeleton to get started quickly;
- A `solution_test.go` file with basic unit tests and benchmarks, for when you
  have found the answer to the daily problem;

It can also download your input for the day's problem, granted you have provided
your adventofcode.com session cookie (see [Session cookie](#session-cookie) for
details).

### Session cookie

When logged in to the adventofcode.com website, your browser has a cookie called
`session`. Retrieve this cookie's value and provide it to the `adventofcode` CLI
to automatically download your input for the day.

## Helpers

This repository includes a `helpers` package with useful functions for
implementing solutions to Advent of Code problems. Feel free to use any of them.

For examples on how to use them, look for functions that start with `Example`.
These are actually unit tests, so you can be sure that they work as described.

## Tests and benchmarks

The scaffolding provided by the `adventofcode` CLI includes unit tests and
benchmarks. To run them, make sure they are properly uncommented in the
`solution_test.go` file, then run these commands:

```bash
# Run all units tests
go test ./y2022/d01
# Run all benchmarks
go test ./y2022/d01 -bench . -benchmem -cpu 1,2,4,8
```

## Configuration

To configure the `adventofcode` CLI, you can use flags, environment variables,
or a configuration file.

### Flags

The CLI is entirely configurable with flags. For a list of available flags, use
these commands:

```bash
adventofcode --help
adventofcode scaffold --help
```

### Environment variables

You can replace any flag in the `adventofcode` CLI with an environment variable.
Environment variables must simply start with `ADVENTOFCODE_`. For example, you
can replace the `--workdir` flag by setting the `ADVENTOFCODE_WORKDIR` variable.

If you have [direnv](https://direnv.net/) installed, you can add a `.envrc` file
to the `adventofcode` directory that looks like this:

```bash
export ADVENTOFCODE_WORKDIR="$(git rev-parse --show-toplevel)"
export ADVENTOFCODE_COOKIE="abcdef0123456789..."
```

### Configuration file

The `adventofcode` CLI automatically looks for a configuration file located in
your home directory: `$HOME/.adventofcode.yaml`. In this file, you can set
values for any flag. For example, your configuration file could look like this:

```yaml
workdir: /Users/arthur/workspace/adventofcode
cookie: abdefg0123456789...
```

## Troubleshooting

If you encounter any problems while using the `adventofcode` CLI, let us know
with a GitHub issue.

---

Made with 💜 by @busser.
