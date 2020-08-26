package miro

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTeamUserConnectionJSON(id string) string {
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

func getTeamUserConnectionsJSON(id string) string {
	return fmt.Sprintf(`[
	{
		"type": "team-user-connection",
		"user": {
		  "type": "user",
		  "name": "Sergey",
		  "id": "user"
		},
		"role": "admin",
		"id": "%s"
	},
	{
		"type": "team-user-connection",
		"user": {
		  "type": "user",
		  "name": "Sergey",
		  "id": "user"
		},
		"role": "admin",
		"id": "%s"
	}
]`, id, id)
}

func getTeamUserConnection(id string) *TeamUserConnection {
	return &TeamUserConnection{
		ID:   id,
		Role: "admin",
		User: &MiniUser{
			ID:   "user",
			Name: "Sergey",
		},
	}
}

func getTeamUserConnections(id string) []*TeamUserConnection {
	return []*TeamUserConnection{getTeamUserConnection(id), getTeamUserConnection(id)}
}

func TestTeamUserConnectionService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *TeamUserConnection
	}{
		"ok": {"1", getTeamUserConnection("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/1", teamUserConnectionsPath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getTeamUserConnectionJSON(tc.id)))
			})

			got, err := client.TeamUserConnection.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
