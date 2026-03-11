package miroir

import "github.com/pkg/errors"

type appContext struct {
	Config Config
}

func (c *getSummariesCommand) Run(app *appContext) error {
	args, err := c.resolveArgs(app.Config)
	if err != nil {
		return errors.Wrap(err, "Fail to resolve `get summaries` arguments.")
	}

	if err := CmdGetSummaries(args); err != nil {
		return errors.Wrap(err, "Fail to command `get summaries`")
	}

	return nil
}

func (c *getReportCommand) Run(app *appContext) error {
	args, err := c.resolveArgs(app.Config)
	if err != nil {
		return errors.Wrap(err, "Fail to resolve `get report` arguments.")
	}

	if err := CmdGetReport(args); err != nil {
		return errors.Wrap(err, "Fail to command `get report`")
	}

	return nil
}

func (c *getResponseBodyCommand) Run(app *appContext) error {
	args, err := c.resolveArgs(app.Config)
	if err != nil {
		return errors.Wrap(err, "Fail to resolve `get response-body` arguments.")
	}

	if err := CmdGetResponseBody(args); err != nil {
		return errors.Wrap(err, "Fail to command `get response-body`")
	}

	return nil
}

func (c *pruneCommand) Run(app *appContext) error {
	args, err := c.resolveArgs(app.Config)
	if err != nil {
		return errors.Wrap(err, "Fail to resolve `prune` arguments.")
	}

	if err := CmdPrune(args); err != nil {
		return errors.Wrap(err, "Fail to command `prune`")
	}

	return nil
}

func (c *getResponseBodyCommand) Validate() error {
	if c.Seq < 1 {
		return errors.New("<seq> must be greater than or equal to 1")
	}

	return nil
}

func (c *getSummariesCommand) resolveArgs(config Config) (*ArgsGetSummaries, error) {
	resolvedAWS := resolveAWSFlags(c.awsFlags, config)
	resolvedDynamo := resolveDynamoFlags(c.dynamoFlags, config)

	if err := requireOption("table", resolvedDynamo.Table); err != nil {
		return nil, err
	}

	return &ArgsGetSummaries{
		Table:            resolvedDynamo.Table,
		RoleARN:          resolvedAWS.RoleARN,
		DynamoDBEndpoint: resolvedDynamo.DynamoDBEndpoint,
		STSEndpoint:      resolvedAWS.STSEndpoint,
		JSON:             c.JSON,
	}, nil
}

func (c *getReportCommand) resolveArgs(config Config) (*ArgsGetReport, error) {
	resolvedAWS := resolveAWSFlags(c.awsFlags, config)
	resolvedS3 := resolveS3Flags(c.s3Flags, config)

	if err := requireOption("bucket", resolvedS3.Bucket); err != nil {
		return nil, err
	}

	return &ArgsGetReport{
		Bucket:       resolvedS3.Bucket,
		BucketPrefix: resolvedS3.BucketPrefix,
		Key:          c.Key,
		RoleARN:      resolvedAWS.RoleARN,
		S3Endpoint:   resolvedS3.S3Endpoint,
		STSEndpoint:  resolvedAWS.STSEndpoint,
	}, nil
}

func (c *getResponseBodyCommand) resolveArgs(config Config) (*ArgsGetResponseBody, error) {
	resolvedAWS := resolveAWSFlags(c.awsFlags, config)
	resolvedS3 := resolveS3Flags(c.s3Flags, config)

	if err := requireOption("bucket", resolvedS3.Bucket); err != nil {
		return nil, err
	}

	side, err := c.resolveSide()
	if err != nil {
		return nil, err
	}

	return &ArgsGetResponseBody{
		Bucket:       resolvedS3.Bucket,
		BucketPrefix: resolvedS3.BucketPrefix,
		Key:          c.Key,
		Seq:          c.Seq,
		Side:         side,
		RoleARN:      resolvedAWS.RoleARN,
		S3Endpoint:   resolvedS3.S3Endpoint,
		STSEndpoint:  resolvedAWS.STSEndpoint,
	}, nil
}

func (c *pruneCommand) resolveArgs(config Config) (*ArgsPrune, error) {
	resolvedAWS := resolveAWSFlags(c.awsFlags, config)
	resolvedS3 := resolveS3Flags(c.s3Flags, config)
	resolvedDynamo := resolveDynamoFlags(c.dynamoFlags, config)

	if err := requireOption("table", resolvedDynamo.Table); err != nil {
		return nil, err
	}
	if err := requireOption("bucket", resolvedS3.Bucket); err != nil {
		return nil, err
	}

	return &ArgsPrune{
		Table:            resolvedDynamo.Table,
		Bucket:           resolvedS3.Bucket,
		BucketPrefix:     resolvedS3.BucketPrefix,
		Dry:              c.Dry,
		RoleARN:          resolvedAWS.RoleARN,
		S3Endpoint:       resolvedS3.S3Endpoint,
		DynamoDBEndpoint: resolvedDynamo.DynamoDBEndpoint,
		STSEndpoint:      resolvedAWS.STSEndpoint,
	}, nil
}

func (c *getResponseBodyCommand) resolveSide() (string, error) {
	switch {
	case c.One && c.Other:
		return "", errors.New("Either --one or --other must be specified, but not both.")
	case c.One:
		return "one", nil
	case c.Other:
		return "other", nil
	default:
		return "", errors.New("Either --one or --other must be specified.")
	}
}

func resolveAWSFlags(flags awsFlags, config Config) awsFlags {
	return awsFlags{
		RoleARN:     firstNonEmpty(flags.RoleARN, config.RoleARN),
		STSEndpoint: firstNonEmpty(flags.STSEndpoint, config.STSEndpoint),
	}
}

func resolveS3Flags(flags s3Flags, config Config) s3Flags {
	return s3Flags{
		Bucket:       firstNonEmpty(flags.Bucket, config.Bucket),
		BucketPrefix: firstNonEmpty(flags.BucketPrefix, config.BucketPrefix),
		S3Endpoint:   firstNonEmpty(flags.S3Endpoint, config.S3Endpoint),
	}
}

func resolveDynamoFlags(flags dynamoFlags, config Config) dynamoFlags {
	return dynamoFlags{
		Table:            firstNonEmpty(flags.Table, config.Table),
		DynamoDBEndpoint: firstNonEmpty(flags.DynamoDBEndpoint, config.DynamoDBEndpoint),
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func requireOption(name, value string) error {
	if value == "" {
		return errors.Errorf("%s is required either via CLI flag or .miroirconfig", name)
	}

	return nil
}
