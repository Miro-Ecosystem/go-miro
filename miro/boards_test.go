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
	testBoardName     = "test-name"
	testBoardViewLink = "https://test-test.com"
	testBoardDesc     = ""
)

func getErrorJSON(status int) string {
	return fmt.Sprintf(`{
	"status": %d,
	"message": "error",
	"type": "error"
}`, status)
}

func getError(status int) *RespError {
	return &RespError{
		Status:  status,
		Message: "error",
		Type:    "error",
	}
}

func getBoardJSON(id string) string {
	return fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"viewLink": "%s",
	"description": "%s",
	"createdAt": "1995-06-15T10:00:00Z",
	"modifiedAt": "1995-06-15T10:00:00Z"
}`, id, testBoardName, testBoardViewLink, testBoardDesc)
}

func getBoard(id string) *Board {
	modifiedAt, _ := time.Parse("1994-03-01T10:00:00Z", "1995-06-15T10:00:00Z")
	createdAt, _ := time.Parse("1994-03-01T10:00:00Z", "1995-06-15T10:00:00Z")

	return &Board{
		ID:         id,
		ViewLink:   testBoardViewLink,
		Name:       testBoardName,
		ModifiedAt: modifiedAt,
		CreatedAt:  createdAt,
	}
}

func TestBoardsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *Board
	}{
		"ok": {"1", getBoard("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", boardsPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getBoardJSON(tc.id)))
			})

			got, err := client.Boards.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}

func TestBoardsService_Get_Error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want string
	}{
		"ok": {"1", "status code not expected, got:404, message:error"},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", boardsPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, fmt.Sprintf(getErrorJSON(http.StatusNotFound)), http.StatusNotFound)
			})

			_, err := client.Boards.Get(context.Background(), tc.id)
			if err == nil {
				t.Fatalf("Should failed")
			}

			if diff := cmp.Diff(err.Error(), tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
