package main

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/go-playground/validator.v9"
)

func TestWriteSummariesText(t *testing.T) {
	summaries := []Summary{
		{
			Hashkey:        "1234567890abcdef",
			Title:          "example",
			SameCount:      1,
			DifferentCount: 2,
			FailureCount:   3,
			BeginTime:      "2026-03-11T10:00:00Z",
		},
	}

	var buf bytes.Buffer
	if err := writeSummaries(&buf, summaries, false); err != nil {
		t.Fatalf("writeSummaries returned error: %v", err)
	}

	want := "2026-03-11T10:00:00Z\t1234567\t1\t2\t3\texample\n"
	if buf.String() != want {
		t.Fatalf("unexpected text output\nwant: %q\ngot:  %q", want, buf.String())
	}
}

func TestWriteSummariesJSON(t *testing.T) {
	summaries := []Summary{
		{
			Hashkey:        "1234567890abcdef",
			Title:          "example",
			OneHost:        "host-a",
			OtherHost:      "host-b",
			SameCount:      0,
			DifferentCount: 2,
			FailureCount:   0,
			BeginTime:      "2026-03-11T10:00:00Z",
			EndTime:        "2026-03-11T10:01:00Z",
			ElapsedSec:     60,
			CheckStatus:    "",
			RetryHash:      "",
			WithZip:        false,
		},
	}

	var buf bytes.Buffer
	if err := writeSummaries(&buf, summaries, true); err != nil {
		t.Fatalf("writeSummaries returned error: %v", err)
	}

	got := buf.String()
	contains := []string{
		`"hashkey": "1234567890abcdef"`,
		`"title": "example"`,
		`"one_host": "host-a"`,
		`"other_host": "host-b"`,
		`"same_count": 0`,
		`"different_count": 2`,
		`"failure_count": 0`,
		`"begin_time": "2026-03-11T10:00:00Z"`,
		`"end_time": "2026-03-11T10:01:00Z"`,
		`"elapsed_sec": 60`,
		`"check_status": ""`,
		`"retry_hash": ""`,
		`"with_zip": false`,
	}

	for _, expected := range contains {
		if !strings.Contains(got, expected) {
			t.Fatalf("expected JSON output to contain %q\noutput: %s", expected, got)
		}
	}

	if strings.Contains(got, `"hashkey": "1234567"`) {
		t.Fatalf("expected full hashkey in JSON output: %s", got)
	}
}

func TestSortSummaries(t *testing.T) {
	summaries := []Summary{
		{Hashkey: "old", BeginTime: "2026-03-11T09:00:00Z"},
		{Hashkey: "new", BeginTime: "2026-03-11T10:00:00Z"},
	}

	sortSummaries(summaries)

	if summaries[0].Hashkey != "new" || summaries[1].Hashkey != "old" {
		t.Fatalf("unexpected order after sort: %+v", summaries)
	}
}

func TestWriteResponseBody(t *testing.T) {
	var buf bytes.Buffer
	if err := writeResponseBody(&buf, "{\"foo\":\"bar\"}"); err != nil {
		t.Fatalf("writeResponseBody returned error: %v", err)
	}

	want := "{\"foo\":\"bar\"}\n"
	if buf.String() != want {
		t.Fatalf("unexpected response body output\nwant: %q\ngot:  %q", want, buf.String())
	}
}

func TestParseResponseBodyArgsOne(t *testing.T) {
	seq, side, err := parseResponseBodyArgs(Args{Seq: "3", One: true})
	if err != nil {
		t.Fatalf("parseResponseBodyArgs returned error: %v", err)
	}

	if seq != 3 || side != "one" {
		t.Fatalf("unexpected parse result: seq=%d side=%s", seq, side)
	}
}

func TestParseResponseBodyArgsOther(t *testing.T) {
	seq, side, err := parseResponseBodyArgs(Args{Seq: "4", Other: true})
	if err != nil {
		t.Fatalf("parseResponseBodyArgs returned error: %v", err)
	}

	if seq != 4 || side != "other" {
		t.Fatalf("unexpected parse result: seq=%d side=%s", seq, side)
	}
}

func TestParseResponseBodyArgsRequiresSide(t *testing.T) {
	_, _, err := parseResponseBodyArgs(Args{Seq: "1"})
	if err == nil {
		t.Fatal("expected error when side is not specified")
	}
}

func TestParseResponseBodyArgsRejectsBothSides(t *testing.T) {
	_, _, err := parseResponseBodyArgs(Args{Seq: "1", One: true, Other: true})
	if err == nil {
		t.Fatal("expected error when both sides are specified")
	}
}

func TestArgsGetResponseBodyValidationRejectsZeroSeq(t *testing.T) {
	v := validator.New()

	args := &ArgsGetResponseBody{
		Bucket: "bucket",
		Key:    "key",
		Seq:    0,
		Side:   "one",
	}

	if err := v.Struct(args); err == nil {
		t.Fatal("expected validation error for seq=0")
	}
}
