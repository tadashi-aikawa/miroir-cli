package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/pkg/errors"
)

type ArgsGetSummaries struct {
	Table            string `validate:"required"`
	RoleARN          string
	S3Endpoint       string
	DynamoDBEndpoint string
	STSEndpoint      string
	JSON             bool
}

// CmdGetSummaries show summaries
func CmdGetSummaries(args *ArgsGetSummaries) error {
	dao, err := NewAwsDao("ap-northeast-1", args.RoleARN, args.S3Endpoint, args.DynamoDBEndpoint, args.STSEndpoint)
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	summaries, err := dao.FetchSummaries(args.Table)
	if err != nil {
		return errors.Wrap(err, "Fail to fetch summaries.")
	}

	sortSummaries(summaries)

	if err := writeSummaries(os.Stdout, summaries, args.JSON); err != nil {
		return errors.Wrap(err, "Fail to write summaries.")
	}

	return nil
}

func sortSummaries(summaries []Summary) {
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].BeginTime > summaries[j].BeginTime
	})
}

func writeSummaries(w io.Writer, summaries []Summary, outputJSON bool) error {
	if outputJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(summaries)
	}

	for _, x := range summaries {
		if _, err := fmt.Fprintf(
			w,
			"%v\t%v\t%v\t%v\t%v\t%v\n",
			x.BeginTime, x.Hashkey[0:7], x.SameCount, x.DifferentCount, x.FailureCount, x.Title,
		); err != nil {
			return err
		}
	}

	return nil
}

type ArgsGetReport struct {
	Bucket           string `validate:"required"`
	BucketPrefix     string
	Key              string `validate:"required"`
	RoleARN          string
	S3Endpoint       string
	DynamoDBEndpoint string
	STSEndpoint      string
}

// CmdGetReport show report
func CmdGetReport(args *ArgsGetReport) error {
	dao, err := NewAwsDao("ap-northeast-1", args.RoleARN, args.S3Endpoint, args.DynamoDBEndpoint, args.STSEndpoint)
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	report, err := dao.FetchReport(args.Bucket, args.BucketPrefix, args.Key)
	if err != nil {
		return errors.Wrap(err, "Fail to fetch report.")
	}

	fmt.Printf("%v\n", report)

	return nil
}

type ArgsPrune struct {
	Table            string `validate:"required"`
	Bucket           string `validate:"required"`
	BucketPrefix     string
	Dry              bool
	RoleARN          string
	S3Endpoint       string
	DynamoDBEndpoint string
	STSEndpoint      string
}

func pruneReport(dao Dao, table, bucket, bucketPrefix, key string, dryRun bool) error {
	hasReport, err := dao.HasReport(bucket, bucketPrefix, key)
	if err != nil {
		return errors.Wrap(err, "Fail to check whether report exists or not.")
	}

	if hasReport {
		log.Printf("[INFO] %v is fine.\n", key)
		return nil
	}

	if dryRun {
		log.Printf("[DRY RUN] %v is removed..\n", key)
	} else {
		if err := dao.RemoveSummary(table, key); err != nil {
			return errors.Wrap(err, "Fail to remove summary")
		}
		log.Printf("[RUN] %v is removed..\n", key)
	}

	return nil
}

// CmdPrune remove summaries if report that associated with key is not existed.
func CmdPrune(args *ArgsPrune) error {
	dao, err := NewAwsDao("ap-northeast-1", args.RoleARN, args.S3Endpoint, args.DynamoDBEndpoint, args.STSEndpoint)
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	summaries, err := dao.FetchSummaries(args.Table)
	if err != nil {
		return errors.Wrap(err, "Fail to fetch summaries.")
	}

	for _, s := range summaries {
		pruneReport(dao, args.Table, args.Bucket, args.BucketPrefix, s.Hashkey, args.Dry)
	}

	return nil
}
