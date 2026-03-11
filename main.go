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
	dynamodbEndpoint := config.DynamoDBEndpoint
	if args.DynamoDBEndpoint != "" {
		dynamodbEndpoint = args.DynamoDBEndpoint
	}
	stsEndpoint := config.STSEndpoint
	if args.STSEndpoint != "" {
		stsEndpoint = args.STSEndpoint
	}
	r := &ArgsGetSummaries{
		Table:            table,
		RoleARN:          roleARN,
		DynamoDBEndpoint: dynamodbEndpoint,
		STSEndpoint:      stsEndpoint,
		JSON:             args.JSON,
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
	s3Endpoint := config.S3Endpoint
	if args.S3Endpoint != "" {
		s3Endpoint = args.S3Endpoint
	}
	stsEndpoint := config.STSEndpoint
	if args.STSEndpoint != "" {
		stsEndpoint = args.STSEndpoint
	}
	r := &ArgsGetReport{
		Bucket:       bucket,
		BucketPrefix: bucketPrefix,
		Key:          args.Key,
		RoleARN:      roleARN,
		S3Endpoint:   s3Endpoint,
		STSEndpoint:  stsEndpoint,
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
	s3Endpoint := config.S3Endpoint
	if args.S3Endpoint != "" {
		s3Endpoint = args.S3Endpoint
	}
	dynamodbEndpoint := config.DynamoDBEndpoint
	if args.DynamoDBEndpoint != "" {
		dynamodbEndpoint = args.DynamoDBEndpoint
	}
	stsEndpoint := config.STSEndpoint
	if args.STSEndpoint != "" {
		stsEndpoint = args.STSEndpoint
	}
	r := &ArgsPrune{
		Table:            table,
		Bucket:           bucket,
		BucketPrefix:     bucketPrefix,
		Dry:              args.Dry,
		RoleARN:          roleARN,
		S3Endpoint:       s3Endpoint,
		DynamoDBEndpoint: dynamodbEndpoint,
		STSEndpoint:      stsEndpoint,
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
