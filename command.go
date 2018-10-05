package main

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// CmdGetSummaries show summaries
func CmdGetSummaries(args Args) error {
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
		fmt.Printf("%v\t%v\t%v\n", x.BeginTime, x.Hashkey[0:12], x.Title)
	}

	return nil
}
