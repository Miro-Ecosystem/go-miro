package miro

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getBoardUserConnectionJSON(id string) string {
	return fmt.Sprintf(`{
    "type": "team-user-connection",
	"user": {
      "type": "user",
      "name": "Sergey",
      "id": "user"
    },
    "role": "admin",
    "id": "%s"
}`, id)
}

func getBoardUserConnection(id string) *BoardUserConnection {
	return &BoardUserConnection{
		ID:   id,
		Role: "admin",
		User: &MiniUser{
			ID:   "user",
			Name: "Sergey",
		},
	}
}

func TestBoardUserConnectionService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *BoardUserConnection
	}{
		"ok": {"1", getBoardUserConnection("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", boardUserConnectionsPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getBoardUserConnectionJSON(tc.id)))
			})

			got, err := client.BoardUserConnection.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
