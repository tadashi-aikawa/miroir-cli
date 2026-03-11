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
