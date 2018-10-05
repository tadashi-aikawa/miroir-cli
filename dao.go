package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type Summary struct {
	Hashkey        string
	Title          string
	OneHost        string `dynamodbav:"one_host"`
	OtherHost      string `dynamodbav:"other_host"`
	SameCount      int    `dynamodbav:"same_count"`
	DifferentCount int    `dynamodbav:"different_count"`
	FailureCount   int    `dynamodbav:"failure_count"`
	BeginTime      string `dynamodbav:"begin_time"`
	EndTime        string `dynamodbav:"end_time"`
	ElapsedSec     int    `dynamodbav:"elapsed_sec"`
	CheckStatus    string `dynamodbav:"check_status"`
	RetryHash      string `dynamodbav:"retry_hash"`
	WithZip        bool   `dynamodbav:"with_zip"`
}

type Dao interface {
	FetchSummaries(table string) ([]Summary, error)
}

type awsClient struct {
	dynamodb *dynamodb.DynamoDB
}

func NewAwsDao(region string) (Dao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load SDK config")
	}
	cfg.Region = region

	var client awsClient
	client.dynamodb = dynamodb.New(cfg)

	return &client, nil
}

func (r *awsClient) FetchSummaries(table string) ([]Summary, error) {
	req := r.dynamodb.ScanRequest(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})

	resp, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to get summaries from "+table)
	}

	var summaries []Summary
	if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &summaries); err != nil {
		return nil, errors.Wrap(err, "Fail to parse summaries with struct `Summary`")
	}

	return summaries, nil
}
