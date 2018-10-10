package main

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

type ArgsGetSummaries struct {
	Table string `validate:"required"`
}

// CmdGetSummaries show summaries
func CmdGetSummaries(args *ArgsGetSummaries) error {
	dao, err := NewAwsDao("ap-northeast-1")
	if err != nil {
		return errors.Wrap(err, "Fail to create aws client.")
	}

	summaries, err := dao.FetchSummaries(args.Table)
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
	Bucket       string `validate:"required"`
	BucketPrefix string
	Key          string `validate:"required"`
}

// CmdGetReport show report
func CmdGetReport(args *ArgsGetReport) error {
	dao, err := NewAwsDao("ap-northeast-1")
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
