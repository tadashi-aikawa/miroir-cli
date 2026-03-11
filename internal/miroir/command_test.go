package miroir

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"
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

func TestPruneSummariesPropagatesError(t *testing.T) {
	dao := &stubDao{
		hasReportByKey: map[string]bool{
			"ok": true,
			"ng": false,
		},
		removeSummaryErrByKey: map[string]error{
			"ng": errStub,
		},
	}

	err := pruneSummaries(dao, "table", "bucket", "prefix", []Summary{
		{Hashkey: "ok"},
		{Hashkey: "ng"},
	}, false)
	if err == nil {
		t.Fatal("expected pruneSummaries to return an error")
	}
}

var errStub = errors.New("stub error")

type stubDao struct {
	hasReportByKey        map[string]bool
	hasReportErrByKey     map[string]error
	removeSummaryErrByKey map[string]error
	removedKeys           []string
}

func (s *stubDao) FetchSummaries(table string) ([]Summary, error) {
	return nil, nil
}

func (s *stubDao) RemoveSummary(table, key string) error {
	if err := s.removeSummaryErrByKey[key]; err != nil {
		return err
	}
	s.removedKeys = append(s.removedKeys, key)
	return nil
}

func (s *stubDao) FetchReport(bucket string, bucketPrefix string, key string) (string, error) {
	return "", nil
}

func (s *stubDao) FetchResponseBody(bucket string, bucketPrefix string, key string, seq int, side string) (string, error) {
	return "", nil
}

func (s *stubDao) HasReport(bucket string, bucketPrefix string, key string) (bool, error) {
	if err := s.hasReportErrByKey[key]; err != nil {
		return false, err
	}
	return s.hasReportByKey[key], nil
}
