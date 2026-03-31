package scaffold

import (
	"strings"
	"testing"
)

func TestSummarizeResults(t *testing.T) {
	results := []WriteResult{
		{Status: StatusCreated},
		{Status: StatusCreated},
		{Status: StatusDirCreated},
		{Status: StatusMerged},
		{Status: StatusSkipped},
		{Status: StatusFailed, Error: "boom"},
	}

	summary := SummarizeResults(results)
	if summary.Created != 2 || summary.Directories != 1 || summary.Merged != 1 || summary.Skipped != 1 || summary.Failed != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if summary.Successful() {
		t.Fatal("summary with failures should not be successful")
	}
	if !strings.Contains(summary.Detail(), "2 files created") || !strings.Contains(summary.Detail(), "1 writes failed") {
		t.Fatalf("unexpected detail string: %q", summary.Detail())
	}
	if summary.FailureReason() == "" {
		t.Fatal("failure reason should be populated when writes fail")
	}
}
