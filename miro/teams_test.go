package miro

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testTeamName = "test-team"
)

func getTeamJSON(id string) string {
	return fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"modifiedAt": "1995-06-15T10:00:00Z",
	"createdAt": "1995-06-15T10:00:00Z"
}`, id, testTeamName)
}

func getTeam(id string) *Team {
	modifiedAt, _ := time.Parse("1994-03-01T10:00:00Z", "1995-06-15T10:00:00Z")
	createdAt, _ := time.Parse("1994-03-01T10:00:00Z", "1995-06-15T10:00:00Z")

	return &Team{
		ID:         id,
		Name:       testTeamName,
		ModifiedAt: modifiedAt,
		CreatedAt:  createdAt,
	}
}

func TestTeamsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *Team
	}{
		"ok": {"1", getTeam("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", teamsPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getTeamJSON(tc.id)))
			})

			got, err := client.Teams.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}

func TestTeamsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *Team
		req  *UpdateTeamRequest
	}{
		"ok": {"1", getTeam("1"), &UpdateTeamRequest{
			"miro",
		}},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", teamsPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getTeamJSON(tc.id)))
			})

			got, err := client.Teams.Update(context.Background(), tc.id, tc.req)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}

func TestTeamsService_GetCurrentUserConnection(t *testing.T) {
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
			mux.HandleFunc(fmt.Sprintf("/%s/%s/%s/me", teamsPath, tc.id, userConnectionsPath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getTeamUserConnectionJSON(tc.id)))
			})

			got, err := client.Teams.GetCurrentUserConnection(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}

func TestTeamsService_Invite(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id    string
		email string
		want  []*TeamUserConnection
	}{
		"ok": {"1", "miro@test.com", getTeamUserConnections("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s/%s/%s", teamsPath, tc.id, userConnectionsPath, teamInvitePath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getTeamUserConnectionsJSON(tc.id)))
			})

			got, err := client.Teams.Invite(context.Background(), tc.id, tc.email)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
