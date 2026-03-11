package main

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

const version = "0.5.0"
const usage = `Miroir CLI.

Usage:
  miroir get summaries [--json] [--table=<table>] [--role-arn=<role_arn>] [--dynamodb-endpoint=<url>] [--sts-endpoint=<url>]
  miroir get report <key> [--bucket=<bucket>] [--bucket-prefix=<bucket-prefix>] [--role-arn=<role_arn>] [--s3-endpoint=<url>] [--sts-endpoint=<url>]
  miroir prune [--table=<table>] [--bucket=<bucket>] [--bucket-prefix=<bucket-prefix>] [--dry] [--role-arn=<role_arn>] [--s3-endpoint=<url>] [--dynamodb-endpoint=<url>] [--sts-endpoint=<url>]
  miroir --help

Options:
  <key>                                 Report key
  -t --table=<table>                    DynamoDB table name
  -b --bucket=<bucket>                  S3 bucket name
  -B --bucket-prefix=<bucket-prefix>    S3 bucket prefix (directory)
  -a --role-arn=<role_arn>              Assume role ARN
  --s3-endpoint=<url>                   S3 endpoint URL
  --dynamodb-endpoint=<url>             DynamoDB endpoint URL
  --sts-endpoint=<url>                  STS endpoint URL
  -d --dry                              Dry run
  --json                                Output summaries as JSON

  -h --help                             Show this screen.
  -v --version                          Version
  `

// Args created by CLI args
type Args struct {
	CmdGet   bool `docopt:"get"`
	CmdPrune bool `docopt:"prune"`

	CmdSummaries bool `docopt:"summaries"`
	CmdReport    bool `docopt:"report"`

	Table            string `docopt:"--table"`
	Bucket           string `docopt:"--bucket"`
	BucketPrefix     string `docopt:"--bucket-prefix"`
	RoleARN          string `docopt:"--role-arn"`
	S3Endpoint       string `docopt:"--s3-endpoint"`
	DynamoDBEndpoint string `docopt:"--dynamodb-endpoint"`
	STSEndpoint      string `docopt:"--sts-endpoint"`
	Key              string `docopt:"<key>"`

	Dry bool `docopt:"--dry"`
	JSON bool `docopt:"--json"`
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
