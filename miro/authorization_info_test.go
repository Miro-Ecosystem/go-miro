package miro

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func getAuthorizationInfoJSON(id string) string {
	return fmt.Sprintf(`{
    "type": "team-user-connection",
	"scopes": [
      "boards:read",
      "team:read"
    ],
	"id": "%s",
    "createdAt": "1995-06-15T10:00:00Z"
}`, id)
}

func getAuthorizationInfo(id string) *AuthorizationInfo {
	createdAt, _ := time.Parse("1994-03-01T10:00:00Z", "1995-06-15T10:00:00Z")

	return &AuthorizationInfo{
		ID:        id,
		Scopes:    []string{"boards:read", "team:read"},
		CreatedAt: createdAt,
	}
}

func TestAuthzService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *AuthorizationInfo
	}{
		"ok": {"1", getAuthorizationInfo("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s", AuthorizationInfoPath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getAuthorizationInfoJSON(tc.id)))
			})

			got, err := client.AuthzInfo.Get(context.Background())
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
