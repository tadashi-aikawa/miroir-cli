package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func createArgsGetSummaries(args Args, config Config) *ArgsGetSummaries {
	table := config.Table
	if args.Table != "" {
		table = args.Table
	}
	r := &ArgsGetSummaries{
		Table: table,
	}

	err := validate.Struct(r)
	if err != nil {
		log.Fatal(err)
	}

	return r
}

func createArgsGetReport(args Args, config Config) *ArgsGetReport {
	bucket := config.Bucket
	if args.Bucket != "" {
		bucket = args.Bucket
	}
	bucketPrefix := config.BucketPrefix
	if args.BucketPrefix != "" {
		bucketPrefix = args.BucketPrefix
	}
	r := &ArgsGetReport{
		Bucket:       bucket,
		BucketPrefix: bucketPrefix,
		Key:          args.Key,
	}

	if err := validate.Struct(r); err != nil {
		log.Fatal(err)
	}

	return r
}

func createArgsPrune(args Args, config Config) *ArgsPrune {
	table := config.Table
	if args.Table != "" {
		table = args.Table
	}
	bucket := config.Bucket
	if args.Bucket != "" {
		bucket = args.Bucket
	}
	bucketPrefix := config.BucketPrefix
	if args.BucketPrefix != "" {
		bucketPrefix = args.BucketPrefix
	}
	r := &ArgsPrune{
		Table:        table,
		Bucket:       bucket,
		BucketPrefix: bucketPrefix,
		Dry:          args.Dry,
	}

	if err := validate.Struct(r); err != nil {
		log.Fatal(err)
	}

	return r
}

func main() {
	validate = validator.New()

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

	case args.CmdPrune:
		if err := CmdPrune(createArgsPrune(args, config)); err != nil {
			log.Fatal(errors.Wrap(err, "Fail to command `prune`"))
		}
	}
}
