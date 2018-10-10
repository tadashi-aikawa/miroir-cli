package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func createArgsGetSummaries(args Args, config Config) ArgsGetSummaries {
	table := config.Table
	if args.Table != "" {
		table = args.Table
	}
	r := ArgsGetSummaries{
		table: table,
	}
	if err := r.validate(); err != nil {
		log.Fatal(err)
	}

	return r
}

func createArgsGetReport(args Args, config Config) ArgsGetReport {
	bucket := config.Bucket
	if args.Bucket != "" {
		bucket = args.Bucket
	}
	bucketPrefix := config.BucketPrefix
	if args.BucketPrefix != "" {
		bucketPrefix = args.BucketPrefix
	}
	r := ArgsGetReport{
		bucket:       bucket,
		bucketPrefix: bucketPrefix,
		key:          args.Key,
	}
	if err := r.validate(); err != nil {
		log.Fatal(err)
	}

	return r
}

func main() {
	args, err := CreateArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to create arguments."))
	}

	config, err := CreateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fail to load `.miroirconfig`."))
	}

	switch true {
	case args.CmdGet:
		switch true {
		case args.CmdSummaries:
			if err := CmdGetSummaries(createArgsGetSummaries(args, config)); err != nil {
				log.Fatal(errors.Wrap(err, "Fail to command `get summaries`"))
			}
		case args.CmdReport:
			if err := CmdGetReport(createArgsGetReport(args, config)); err != nil {
				log.Fatal(errors.Wrap(err, "Fail to command `get report`"))
			}
		}
	}
}
