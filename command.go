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
	Table            string
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
	Bucket       string
	BucketPrefix string
	Key          string
	RoleARN      string
	S3Endpoint   string
	STSEndpoint  string
}

type ArgsGetResponseBody struct {
	Bucket       string
	BucketPrefix string
	Key          string
	Seq          int
	Side         string
	RoleARN      string
	S3Endpoint   string
	STSEndpoint  string
}

// CmdGetReport show report
func CmdGetReport(args *ArgsGetReport) error {
	dao, err := NewAwsDao("ap-northeast-1", args.RoleARN, args.S3Endpoint, "", args.STSEndpoint)
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

// CmdGetResponseBody show response body
func CmdGetResponseBody(args *ArgsGetResponseBody) error {
	dao, err := NewAwsDao("ap-northeast-1", args.RoleARN, args.S3Endpoint, "", args.STSEndpoint)
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	body, err := dao.FetchResponseBody(args.Bucket, args.BucketPrefix, args.Key, args.Seq, args.Side)
	if err != nil {
		return errors.Wrap(err, "Fail to fetch response body.")
	}

	if err := writeResponseBody(os.Stdout, body); err != nil {
		return errors.Wrap(err, "Fail to write response body.")
	}

	return nil
}

func writeResponseBody(w io.Writer, body string) error {
	_, err := fmt.Fprintf(w, "%v\n", body)
	return err
}

type ArgsPrune struct {
	Table            string
	Bucket           string
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

	if err := pruneSummaries(dao, args.Table, args.Bucket, args.BucketPrefix, summaries, args.Dry); err != nil {
		return err
	}

	return nil
}

func pruneSummaries(dao Dao, table, bucket, bucketPrefix string, summaries []Summary, dryRun bool) error {
	for _, s := range summaries {
		if err := pruneReport(dao, table, bucket, bucketPrefix, s.Hashkey, dryRun); err != nil {
			return errors.Wrap(err, "Fail to prune summary")
		}
	}

	return nil
}
