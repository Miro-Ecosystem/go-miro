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

func getBoardJSON(id string) string {
	return fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"viewLink": "%s",
	"description": "%s",
	"picture": null,
	"createdAt": "1995-06-15T10:00:00Z",
	"modifiedAt": "1995-06-15T10:00:00Z"
}`, id, testBoardName, testBoardViewLink, testBoardDesc)
}

func getBoardListJSON() string {
	return fmt.Sprintf(`{
	"limit": 10,
	"offset": 0
}`)

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

func getBoardList() *ListBoardsResponse {
	return &ListBoardsResponse{
		Limit: 10,
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

func TestBoardsService_Share(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id     string
		emails []string
		want   *ListBoardsResponse
	}{
		"ok": {"1", []string{"keke@miro.com", "miro@keke.com"}, getBoardList()},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s/share", boardsPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getBoardListJSON()))
			})

			got, err := client.Boards.Share(context.Background(), tc.id, &ShareBoardRequest{
				tc.emails,
			})

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
				addHeader(w)
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
