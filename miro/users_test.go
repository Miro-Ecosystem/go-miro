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
	testUserName     = "test-name"
	testUserState    = "test-state"
	testUserRole     = "test-role"
	testUserIndustry = "test-ind"
	testUserCompany  = "test-comp"
	testUserEmail    = "miro@test.com"
)

func getUserJSON(id string) string {
	return fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"state": "%s",
	"role": "%s",
	"industry": "%s",
	"company": "%s",
	"email": "%s",
	"createdAt": "1995-06-15T10:00:00Z"
}`, id, testUserName, testUserState, testUserRole, testUserIndustry, testUserCompany, testUserEmail)
}

func getUser(id string) *User {
	createdAt, _ := time.Parse("1994-03-01T10:00:00Z", "1995-06-15T10:00:00Z")

	return &User{
		ID:        id,
		State:     testUserState,
		Name:      testUserName,
		Role:      testUserRole,
		Industry:  testUserIndustry,
		Company:   testUserCompany,
		Email:     testUserEmail,
		CreatedAt: createdAt,
	}
}

func TestUsersService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *User
	}{
		"ok": {"1", getUser("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", usersPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getUserJSON(tc.id)))
			})

			got, err := client.Users.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
