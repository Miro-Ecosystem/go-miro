package miro

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testPictureImageURL = "test-image-url"
)

func getPictureJSON(id string) string {
	return fmt.Sprintf(`{
	"id": "%s",
	"imageUrl": "%s"
}`, id, testPictureImageURL)
}

func getPicture(id string) *Picture {
	return &Picture{
		ID:       id,
		ImageURL: testPictureImageURL,
	}
}

func TestPicturesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		id   string
		want *Picture
	}{
		"ok": {"1", getPicture("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s", picturesPath, tc.id), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getPictureJSON(tc.id)))
			})

			got, err := client.Picture.Get(context.Background(), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}
