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
	roleARN := config.RoleARN
	if args.RoleARN != "" {
		roleARN = args.RoleARN
	}
	r := &ArgsGetSummaries{
		Table:   table,
		RoleARN: roleARN,
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
	roleARN := config.RoleARN
	if args.RoleARN != "" {
		roleARN = args.RoleARN
	}
	r := &ArgsGetReport{
		Bucket:       bucket,
		BucketPrefix: bucketPrefix,
		Key:          args.Key,
		RoleARN:      roleARN,
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
	roleARN := config.RoleARN
	if args.RoleARN != "" {
		roleARN = args.RoleARN
	}
	r := &ArgsPrune{
		Table:        table,
		Bucket:       bucket,
		BucketPrefix: bucketPrefix,
		Dry:          args.Dry,
		RoleARN:      roleARN,
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
		switch err {
		case ErrorHomeDirIsNotFound:
			log.Printf("[WARN] Home directory is not found and can't load .miroirconfig.... continue...")
		case ErrorConfigIsNotFound:
			log.Printf("[WARN] .miroirconfig is not found.... continue...")
		default:
			log.Fatal(errors.Wrap(err, "Fail to load `.miroirconfig`."))
		}
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
