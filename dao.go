package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pkg/errors"
)

// Summary of report
type Summary struct {
	Hashkey        string `json:"hashkey"`
	Title          string `json:"title"`
	OneHost        string `dynamodbav:"one_host" json:"one_host"`
	OtherHost      string `dynamodbav:"other_host" json:"other_host"`
	SameCount      int    `dynamodbav:"same_count" json:"same_count"`
	DifferentCount int    `dynamodbav:"different_count" json:"different_count"`
	FailureCount   int    `dynamodbav:"failure_count" json:"failure_count"`
	BeginTime      string `dynamodbav:"begin_time" json:"begin_time"`
	EndTime        string `dynamodbav:"end_time" json:"end_time"`
	ElapsedSec     int    `dynamodbav:"elapsed_sec" json:"elapsed_sec"`
	CheckStatus    string `dynamodbav:"check_status" json:"check_status"`
	RetryHash      string `dynamodbav:"retry_hash" json:"retry_hash"`
	WithZip        bool   `dynamodbav:"with_zip" json:"with_zip"`
}

// Dao can fetch data
type Dao interface {
	FetchSummaries(table string) ([]Summary, error)
	RemoveSummary(table, key string) error
	FetchReport(bucket string, BucketPrefix string, key string) (string, error)
	HasReport(bucket string, BucketPrefix string, key string) (bool, error)
}

type awsClient struct {
	dynamodb *dynamodb.Client
	s3       *s3.Client
}

type endpoints struct {
	S3       string
	DynamoDB string
	STS      string
}

func (e endpoints) hasCustomEndpoint() bool {
	return e.S3 != "" || e.DynamoDB != "" || e.STS != ""
}

func shouldUseDummyCredentials(e endpoints) bool {
	if !e.hasCustomEndpoint() {
		return false
	}

	return os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_SECRET_ACCESS_KEY") == ""
}

func (r *awsClient) fetchJSON(bucket string, key string) (interface{}, error) {
	resp, err := r.s3.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
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
func NewAwsDao(region string, roleARN string, s3Endpoint string, dynamodbEndpoint string, stsEndpoint string) (Dao, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, errors.Wrap(err, "unable to load SDK config")
	}

	endpoints := endpoints{
		S3:       s3Endpoint,
		DynamoDB: dynamodbEndpoint,
		STS:      stsEndpoint,
	}

	if shouldUseDummyCredentials(endpoints) {
		cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("test", "test", ""))
	}

	if roleARN != "" {
		provider := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg, func(o *sts.Options) {
			if endpoints.STS != "" {
				o.BaseEndpoint = aws.String(endpoints.STS)
			}
		}), roleARN)
		cfg.Credentials = aws.NewCredentialsCache(provider)
	}

	var client awsClient
	client.dynamodb = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoints.DynamoDB != "" {
			o.BaseEndpoint = aws.String(endpoints.DynamoDB)
		}
	})
	client.s3 = s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoints.S3 != "" {
			o.BaseEndpoint = aws.String(endpoints.S3)
			o.UsePathStyle = true
		}
	})

	return &client, nil
}

func (r *awsClient) FetchSummaries(table string) ([]Summary, error) {
	resp, err := r.dynamodb.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Fail to get summaries from "+table)
	}

	var summaries []Summary
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &summaries); err != nil {
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

func (r *awsClient) RemoveSummary(table, key string) error {
	_, err := r.dynamodb.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]dynamodbTypes.AttributeValue{
			"hashkey": &dynamodbTypes.AttributeValueMemberS{
				Value: key,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "Fail to delete summary "+table)
	}

	return nil
}

func (r *awsClient) HasReport(bucket string, BucketPrefix string, key string) (bool, error) {
	var prefix string
	if BucketPrefix != "" {
		prefix += BucketPrefix + "/"
	}

	hashDirKey := fmt.Sprintf("%sresults/%s", prefix, key)

	resp, err := r.s3.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(hashDirKey),
	})
	if err != nil {
		return false, errors.Wrap(err, "Fail to check whether report exists or not: "+key)
	}

	if resp.KeyCount == nil {
		return false, nil
	}

	return *resp.KeyCount > 0, nil
}
