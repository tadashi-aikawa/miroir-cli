package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

// Summary of report
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

// Dao can fetch data
type Dao interface {
	FetchSummaries(table string) ([]Summary, error)
	FetchReport(bucket string, BucketPrefix string, key string) (string, error)
}

type awsClient struct {
	dynamodb *dynamodb.DynamoDB
	s3       *s3.S3
}

func (r *awsClient) fetchJSON(bucket string, key string) (interface{}, error) {
	req := r.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	resp, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to get report: "+key)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, errors.Wrap(err, "Fail to read report: "+key)
	}

	var jsonMap interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonMap); err != nil {
		return nil, errors.Wrap(err, "Fail to parse as json.")
	}

	return jsonMap, nil
}

// NewAwsDao creates dao instance
func NewAwsDao(region string) (Dao, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load SDK config")
	}
	cfg.Region = region

	var client awsClient
	client.dynamodb = dynamodb.New(cfg)
	client.s3 = s3.New(cfg)

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

func (r *awsClient) FetchReport(bucket string, BucketPrefix string, key string) (string, error) {
	var prefix string
	if BucketPrefix != "" {
		prefix += BucketPrefix + "/"
	}

	trialsKey := fmt.Sprintf("%sresults/%s/trials.json", prefix, key)
	trials, err := r.fetchJSON(bucket, trialsKey)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Fail to fetch json: %s (%s)", trialsKey, bucket))
	}

	withoutTrialsKey := fmt.Sprintf("%sresults/%s/report-without-trials.json", prefix, key)
	report, err := r.fetchJSON(bucket, withoutTrialsKey)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Fail to fetch json: %s (%s)", withoutTrialsKey, bucket))
	}

	report.(map[string]interface{})["trials"] = trials

	bs, err := json.Marshal(report)
	if err != nil {
		return "", errors.Wrap(err, "Fail to parse json to string")
	}

	return string(bs), nil
}
