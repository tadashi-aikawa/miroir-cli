package miroir

import (
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
)

type CLI struct {
	Version kong.VersionFlag `name:"version" short:"v" help:"Version."`
	Get     getCommand       `cmd:"" help:"Get reports and summaries."`
	Prune   pruneCommand     `cmd:"" help:"Remove summaries without reports."`
}

type getCommand struct {
	Summaries    getSummariesCommand    `cmd:"" help:"Show summaries."`
	Report       getReportCommand       `cmd:"" help:"Show report."`
	ResponseBody getResponseBodyCommand `cmd:"" name:"response-body" help:"Show response body."`
}

type awsFlags struct {
	RoleARN     string `name:"role-arn" short:"a" help:"Assume role ARN."`
	STSEndpoint string `name:"sts-endpoint" help:"STS endpoint URL."`
}

type s3Flags struct {
	Bucket       string `name:"bucket" short:"b" help:"S3 bucket name."`
	BucketPrefix string `name:"bucket-prefix" short:"B" help:"S3 bucket prefix (directory)."`
	S3Endpoint   string `name:"s3-endpoint" help:"S3 endpoint URL."`
}

type dynamoFlags struct {
	Table            string `name:"table" short:"t" help:"DynamoDB table name."`
	DynamoDBEndpoint string `name:"dynamodb-endpoint" help:"DynamoDB endpoint URL."`
}

type getSummariesCommand struct {
	awsFlags    `embed:""`
	dynamoFlags `embed:""`
	JSON        bool `name:"json" help:"Output summaries as JSON."`
}

type getReportCommand struct {
	awsFlags `embed:""`
	s3Flags  `embed:""`
	Key      string `arg:"" name:"key" help:"Report key."`
}

type getResponseBodyCommand struct {
	awsFlags `embed:""`
	s3Flags  `embed:""`
	Key      string `arg:"" name:"key" help:"Report key."`
	Seq      int    `arg:"" name:"seq" help:"Trial sequence number."`
	One      bool   `name:"one" xor:"side" required:"" help:"Fetch body for the one side."`
	Other    bool   `name:"other" xor:"side" required:"" help:"Fetch body for the other side."`
}

type pruneCommand struct {
	awsFlags    `embed:""`
	s3Flags     `embed:""`
	dynamoFlags `embed:""`
	Dry         bool `name:"dry" short:"d" help:"Dry run."`
}

func ParseCLI(argv []string) (*CLI, *kong.Kong, *kong.Context, error) {
	cli := &CLI{}

	parser, err := kong.New(
		cli,
		kong.Name("miroir"),
		kong.Description("Miroir CLI."),
		kong.Vars{
			"version": version,
		},
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "Fail to initialize arguments parser.")
	}

	ctx, err := parser.Parse(argv)
	return cli, parser, ctx, err
}
