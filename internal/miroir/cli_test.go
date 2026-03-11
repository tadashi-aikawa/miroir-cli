package miroir

import (
	"strings"
	"testing"
)

func TestParseCLIGetSummariesJSON(t *testing.T) {
	cli, _, ctx, err := ParseCLI([]string{"get", "summaries", "--json"})
	if err != nil {
		t.Fatalf("ParseCLI returned error: %v", err)
	}

	if ctx.Command() != "get summaries" {
		t.Fatalf("unexpected command: %s", ctx.Command())
	}

	if !cli.Get.Summaries.JSON {
		t.Fatalf("expected --json to be true: %+v", cli.Get.Summaries)
	}
}

func TestParseCLIGetResponseBodyOne(t *testing.T) {
	cli, _, ctx, err := ParseCLI([]string{"get", "response-body", "report-key", "3", "--one"})
	if err != nil {
		t.Fatalf("ParseCLI returned error: %v", err)
	}

	if ctx.Command() != "get response-body <key> <seq>" {
		t.Fatalf("unexpected command: %s", ctx.Command())
	}

	if cli.Get.ResponseBody.Key != "report-key" || cli.Get.ResponseBody.Seq != 3 || !cli.Get.ResponseBody.One || cli.Get.ResponseBody.Other {
		t.Fatalf("unexpected parsed args: %+v", cli.Get.ResponseBody)
	}
}

func TestParseCLIGetResponseBodyOther(t *testing.T) {
	cli, _, _, err := ParseCLI([]string{"get", "response-body", "report-key", "4", "--other"})
	if err != nil {
		t.Fatalf("ParseCLI returned error: %v", err)
	}

	if !cli.Get.ResponseBody.Other || cli.Get.ResponseBody.One {
		t.Fatalf("unexpected side flags: %+v", cli.Get.ResponseBody)
	}
}

func TestParseCLIGetResponseBodyRequiresSide(t *testing.T) {
	_, _, _, err := ParseCLI([]string{"get", "response-body", "report-key", "1"})
	if err == nil {
		t.Fatal("expected error when side is not specified")
	}
}

func TestParseCLIGetResponseBodyRejectsBothSides(t *testing.T) {
	_, _, _, err := ParseCLI([]string{"get", "response-body", "report-key", "1", "--one", "--other"})
	if err == nil {
		t.Fatal("expected error when both sides are specified")
	}
}

func TestParseCLIGetResponseBodyRejectsZeroSeq(t *testing.T) {
	_, _, _, err := ParseCLI([]string{"get", "response-body", "report-key", "0", "--one"})
	if err == nil {
		t.Fatal("expected error for seq=0")
	}

	if !strings.Contains(err.Error(), "greater than or equal to 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetReportCommandResolveArgsPrefersCLIOverConfig(t *testing.T) {
	cmd := getReportCommand{
		awsFlags: awsFlags{
			RoleARN: "cli-role",
		},
		s3Flags: s3Flags{
			Bucket:       "cli-bucket",
			BucketPrefix: "cli-prefix",
			S3Endpoint:   "http://cli-s3",
		},
		Key: "report-key",
	}

	args, err := cmd.resolveArgs(Config{
		Bucket:       "config-bucket",
		BucketPrefix: "config-prefix",
		RoleARN:      "config-role",
		S3Endpoint:   "http://config-s3",
		STSEndpoint:  "http://config-sts",
	})
	if err != nil {
		t.Fatalf("resolveArgs returned error: %v", err)
	}

	if args.Bucket != "cli-bucket" || args.BucketPrefix != "cli-prefix" || args.RoleARN != "cli-role" || args.S3Endpoint != "http://cli-s3" || args.STSEndpoint != "http://config-sts" {
		t.Fatalf("unexpected resolved args: %+v", args)
	}
}

func TestGetSummariesCommandResolveArgsUsesConfigFallback(t *testing.T) {
	cmd := getSummariesCommand{
		JSON: true,
	}

	args, err := cmd.resolveArgs(Config{
		Table:            "config-table",
		RoleARN:          "config-role",
		DynamoDBEndpoint: "http://config-dynamodb",
		STSEndpoint:      "http://config-sts",
	})
	if err != nil {
		t.Fatalf("resolveArgs returned error: %v", err)
	}

	if args.Table != "config-table" || args.RoleARN != "config-role" || args.DynamoDBEndpoint != "http://config-dynamodb" || args.STSEndpoint != "http://config-sts" || !args.JSON {
		t.Fatalf("unexpected resolved args: %+v", args)
	}
}
