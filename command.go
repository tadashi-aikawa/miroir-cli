package main

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

type ArgsGetSummaries struct {
	table string
}

func (r *ArgsGetSummaries) validate() error {
	if r.table == "" {
		return errors.New("Table is required")
	}

	return nil
}

// CmdGetSummaries show summaries
func CmdGetSummaries(args ArgsGetSummaries) error {
	dao, err := NewAwsDao("ap-northeast-1")
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	summaries, err := dao.FetchSummaries(args.table)
	if err != nil {
		return errors.Wrap(err, "Fail to fetch summaries.")
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].BeginTime > summaries[j].BeginTime
	})

	for _, x := range summaries {
		fmt.Printf("%v\t%v\t%v\n", x.BeginTime, x.Hashkey, x.Title)
	}

	return nil
}

type ArgsGetReport struct {
	bucket       string
	bucketPrefix string
	key          string
}

func (r *ArgsGetReport) validate() error {
	if r.bucket == "" {
		return errors.New("Bucket is required")
	}
	if r.key == "" {
		return errors.New("Key is required")
	}

	return nil
}

// CmdGetReport show report
func CmdGetReport(args ArgsGetReport) error {
	dao, err := NewAwsDao("ap-northeast-1")
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	report, err := dao.FetchReport(args.bucket, args.bucketPrefix, args.key)
	if err != nil {
		return errors.Wrap(err, "Fail to fetch report.")
	}

	fmt.Printf("%v\n", report)

	return nil
}
