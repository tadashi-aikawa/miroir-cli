package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.1.0"
const usage = `Miroir CLI.

Usage:
  miroir get summaries [--table=<table>] [--role-arn=<role_arn>]
  miroir get report <key> [--bucket=<bucket>] [--bucket-prefix=<bucket-prefix>] [--role-arn=<role_arn>]
  miroir prune [--table=<table>] [--bucket=<bucket>] [--bucket-prefix=<bucket-prefix>] [--dry] [--role-arn=<role_arn>]
  miroir --help

Options:
  <key>                                 Report key
  -t --table=<table>                    DynamoDB table name
  -b --bucket=<bucket>                  S3 bucket name
  -B --bucket-prefix=<bucket-prefix>    S3 bucket prefix (directory)
  -a --role-arn=<role_arn>              Assume role ARN
  -d --dry                              Dry run

  -h --help                             Show this screen.
  -v --version                          Version
  `

// Args created by CLI args
type Args struct {
	CmdGet   bool `docopt:"get"`
	CmdPrune bool `docopt:"prune"`

	CmdSummaries bool `docopt:"summaries"`
	CmdReport    bool `docopt:"report"`

	Table        string `docopt:"--table"`
	Bucket       string `docopt:"--bucket"`
	BucketPrefix string `docopt:"--bucket-prefix"`
	RoleARN      string `docopt:"--role-arn"`
	Key          string `docopt:"<key>"`

	Dry bool `docopt:"--dry"`
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
