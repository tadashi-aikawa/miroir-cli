package miroir

import "testing"

func TestBuildResultsPrefix(t *testing.T) {
	if got := buildResultsPrefix("", "abc"); got != "results/abc" {
		t.Fatalf("unexpected prefix without bucket prefix: %s", got)
	}

	if got := buildResultsPrefix("prod", "abc"); got != "prod/results/abc" {
		t.Fatalf("unexpected prefix with bucket prefix: %s", got)
	}
}

func TestBuildResponseBodyKey(t *testing.T) {
	got := buildResponseBodyKey("prod", "hashkey", "one/(3)example")
	want := "prod/results/hashkey/one/(3)example"
	if got != want {
		t.Fatalf("unexpected body key\nwant: %s\ngot:  %s", want, got)
	}
}

func TestFindTrialBySeq(t *testing.T) {
	trial, err := findTrialBySeq([]reportTrial{
		{Seq: 1},
		{Seq: 3},
	}, 3)
	if err != nil {
		t.Fatalf("findTrialBySeq returned error: %v", err)
	}

	if trial.Seq != 3 {
		t.Fatalf("unexpected trial: %+v", trial)
	}
}

func TestFindTrialBySeqNotFound(t *testing.T) {
	_, err := findTrialBySeq([]reportTrial{{Seq: 1}}, 2)
	if err == nil {
		t.Fatal("expected error when seq is not found")
	}
}

func TestGetTrialSideFile(t *testing.T) {
	trial := &reportTrial{
		One:   reportTrialSide{File: "one/(3)example"},
		Other: reportTrialSide{File: "other/(3)example"},
	}

	file, err := getTrialSideFile(trial, "one")
	if err != nil {
		t.Fatalf("getTrialSideFile returned error: %v", err)
	}
	if file != "one/(3)example" {
		t.Fatalf("unexpected one file: %s", file)
	}

	file, err = getTrialSideFile(trial, "other")
	if err != nil {
		t.Fatalf("getTrialSideFile returned error: %v", err)
	}
	if file != "other/(3)example" {
		t.Fatalf("unexpected other file: %s", file)
	}
}

func TestGetTrialSideFileRequiresStoredBody(t *testing.T) {
	_, err := getTrialSideFile(&reportTrial{}, "one")
	if err == nil {
		t.Fatal("expected error when file is empty")
	}
}

func TestGetTrialSideFileRejectsUnsupportedSide(t *testing.T) {
	_, err := getTrialSideFile(&reportTrial{}, "invalid")
	if err == nil {
		t.Fatal("expected error for unsupported side")
	}
}

func TestValidateResponseBodyFilePath(t *testing.T) {
	if err := validateResponseBodyFilePath("responses/one/(3)example.json"); err != nil {
		t.Fatalf("validateResponseBodyFilePath returned error: %v", err)
	}
}

func TestValidateResponseBodyFilePathRejectsAbsolutePath(t *testing.T) {
	err := validateResponseBodyFilePath("/responses/one.json")
	if err == nil {
		t.Fatal("expected error for absolute path")
	}
}

func TestValidateResponseBodyFilePathRejectsTraversal(t *testing.T) {
	err := validateResponseBodyFilePath("responses/../secret.json")
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
}

func TestValidateResponseBodyFilePathRejectsEmptySegment(t *testing.T) {
	err := validateResponseBodyFilePath("responses//one.json")
	if err == nil {
		t.Fatal("expected error for empty segment")
	}
}

func TestValidateResponseBodySize(t *testing.T) {
	if err := validateResponseBodySize(maxResponseBodySize); err != nil {
		t.Fatalf("validateResponseBodySize returned error: %v", err)
	}
}

func TestValidateResponseBodySizeRejectsLargeObject(t *testing.T) {
	err := validateResponseBodySize(maxResponseBodySize + 1)
	if err == nil {
		t.Fatal("expected error for oversized response body")
	}
}
