package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.1.0"
const usage = `Miroir CLI.

Usage:
  miroir get summaries --table=<table>
  miroir --help

Options:
  -t --table=<table>                 DynamoDB table name
  -h --help                          Show this screen.
  -v --version                       Version
  `

// Args created by CLI args
type Args struct {
	CmdGet       bool `docopt:"get"`
	CmdSummaries bool `docopt:"summaries"`

	Table string `docopt:"--table"`
}

// CreateArgs creates Args
func CreateArgs(usage string, argv []string, version string) (Args, error) {
	parser := &docopt.Parser{
		HelpHandler:  docopt.PrintHelpOnly,
		OptionsFirst: false,
	}

	opts, err := parser.ParseArgs(usage, argv, version)
	if err != nil {
		return Args{}, errors.Wrap(err, "Fail to parse arguments.")
	}

	var args Args
	opts.Bind(&args)

	return args, nil
}
