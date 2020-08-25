package miro

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testAuditLogsLimit    = 10
	testAuditLogsOffset   = 0
	testAuditLogsSize     = 8000
	testAuditLogsNextLink = "https://test-test.com/next"
)

func getAuditLogJSON() string {
	return fmt.Sprintf(`{
	"limit": %d,
	"offset": %d,
    "size": %d,
	"nextLink": "%s",
	"data": [
		{
			"type": "event",
			"event": "board_opened",
			"details": {
				"role": "owner"
			},
			"context": {
				"organization": {
					"type": "organization",
					"name": "miro",
					"id": "miro"
				}
			},
			"id": "log",
			"createdAt": "1994-03-01T10:00:00Z"
		}
	]
}`, testAuditLogsLimit, testAuditLogsOffset, testAuditLogsSize, testAuditLogsNextLink)
}

func getAuditLog() *AuditLog {

	return &AuditLog{
		Limit:    testAuditLogsLimit,
		Offset:   testAuditLogsOffset,
		Size:     testAuditLogsSize,
		NextLink: testAuditLogsNextLink,
		Data: []Data{{
			ID:    "log",
			Event: "board_opened",
		}},
	}
}

func TestAuditLogsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		want *AuditLog
	}{
		"ok": {getAuditLog()},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s", auditLogsPath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getAuditLogJSON()))
			})

			got, err := client.AuditLogs.Get(context.Background())
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
