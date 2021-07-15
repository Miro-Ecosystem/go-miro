package miro

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testPictureImageURL = "test-image-url"
)

func getPictureJSON(pictureType string, id string) string {
	return fmt.Sprintf(`{
	"type": "%s",
	"id": "%s",
	"imageUrl": "%s"
}`, pictureType, id, testPictureImageURL)
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
		pictureType string
		id   string
		want *Picture
	}{
		"boards": {"boards", "1", getPicture("1")},
		"teams": {"teams", "1", getPicture("1")},
		"users": {"users", "1", getPicture("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s/%s", tc.pictureType, tc.id, picturesPath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getPictureJSON(tc.pictureType, tc.id)))
			})

			got, err := client.Picture.Get(context.Background(), pictureType(tc.pictureType), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("Diff: %s(-got +want)", diff)
			}
		})
	}
}

func TestPicturesService_Upsert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		pictureType string
		id    string
		Image *bytes.Buffer
		want  *Picture
	}{
		"boards": {"boards", "1", bytes.NewBuffer(make([]byte, 0, 10)), getPicture("1")},
		"teams": {"teams", "1", bytes.NewBuffer(make([]byte, 0, 10)), getPicture("1")},
		"users": {"users", "1", bytes.NewBuffer(make([]byte, 0, 10)), getPicture("1")},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s/%s", tc.pictureType, tc.id, picturesPath), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, fmt.Sprintf(getPictureJSON(tc.pictureType, tc.id)))
			})

			got, err := client.Picture.Upsert(context.Background(), pictureType(tc.pictureType), tc.id, &UpsertPictureRequest{
				tc.Image,
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

func TestPicturesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	tcs := map[string]struct {
		pictureType string
		id   string
		want *Picture
	}{
		"boards": {"boards", "1", getPicture("1")},
		"teams": {"teams", "1", getPicture("1")},
		"users": {"users", "1", getPicture("1")},
	}


	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			mux.HandleFunc(fmt.Sprintf("/%s/%s/%s", tc.pictureType, tc.id, picturesPath), func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
				w.Write([]byte("{}"))
			})

			err := client.Picture.Delete(context.Background(), pictureType(tc.pictureType), tc.id)
			if err != nil {
				t.Fatalf("Failed: %v", err)
			}

		})
	}
}
