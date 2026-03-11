package main

import "testing"

func TestCreateArgs_GetSummariesJSON(t *testing.T) {
	args, err := CreateArgs(usage, []string{"get", "summaries", "--json"}, version)
	if err != nil {
		t.Fatalf("CreateArgs returned error: %v", err)
	}

	if !args.CmdGet || !args.CmdSummaries {
		t.Fatalf("expected get summaries command to be parsed: %+v", args)
	}

	if !args.JSON {
		t.Fatalf("expected --json to be true: %+v", args)
	}
}

func TestCreateArgs_GetResponseBodyOne(t *testing.T) {
	args, err := CreateArgs(usage, []string{"get", "response-body", "report-key", "3", "--one"}, version)
	if err != nil {
		t.Fatalf("CreateArgs returned error: %v", err)
	}

	if !args.CmdGet || !args.CmdResponseBody {
		t.Fatalf("expected get response-body command to be parsed: %+v", args)
	}

	if args.Key != "report-key" || args.Seq != "3" || !args.One || args.Other {
		t.Fatalf("unexpected parsed args: %+v", args)
	}
}

func TestCreateArgs_GetResponseBodyOther(t *testing.T) {
	args, err := CreateArgs(usage, []string{"get", "response-body", "report-key", "4", "--other"}, version)
	if err != nil {
		t.Fatalf("CreateArgs returned error: %v", err)
	}

	if !args.Other || args.One {
		t.Fatalf("unexpected side flags: %+v", args)
	}
}
