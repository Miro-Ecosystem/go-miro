package miro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	picturesPath = "pictures"
)

// PicturesService handles communication to Miro Pictures API.
//
// API doc: https://developers.miro.com/reference#picture-object
type PicturesService service

// User object represents Miro User.
//
// API doc: https://developers.miro.com/reference#user-object
//go:generate gomodifytags -file $GOFILE -struct Picture -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Picture -add-tags json -w -transform camelcase
type Picture struct {
	ID       string `json:"id"`
	ImageURL string `json:"imageURL"`
}

// MiniPicture object represents Miro Mini picture.
//go:generate gomodifytags -file $GOFILE -struct MiniPicture -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct MiniPicture -add-tags json -w -transform camelcase
type MiniPicture struct {
	ID       string `json:"id"`
	ImageURL string `json:"imageURL"`
}

// Get gets picture by Picture ID.
//
// API doc: https://developers.miro.com/reference#get-user
func (s *PicturesService) Get(ctx context.Context, id string) (*Picture, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("type/%s/%s", id, picturesPath))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	picture := &Picture{}
	if err := json.NewDecoder(resp.Body).Decode(picture); err != nil {
		return nil, err
	}

	return picture, nil
}

// UpsertPictureRequest represents update board request payload.
//
//go:generate gomodifytags -file $GOFILE -struct UpsertPictureRequest -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct UpsertPictureRequest -add-tags json -w -transform camelcase
type UpsertPictureRequest struct {
	Image *bytes.Buffer
}

// Upsert upserts a picture by Picture ID.
//
// API doc: https://developers.miro.com/reference#create-or-update-picture
func (s *PicturesService) Upsert(ctx context.Context, id string, request *UpsertPictureRequest) (*Picture, error) {
	req, err := s.client.NewPostRequest(fmt.Sprintf("type/%s/%s", id, picturesPath), request)
	if err != nil {
		return nil, err
	}

	req.Header.Del("Content-Type")
	req.Header.Set("content-type", "multipart/form-data")

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	picture := &Picture{}
	if err := json.NewDecoder(resp.Body).Decode(picture); err != nil {
		return nil, err
	}

	return picture, nil
}

// Delete deletes picture by Picture ID.
//
// API doc: https://developers.miro.com/reference#delete-picture
func (s *PicturesService) Delete(ctx context.Context, id string) error {
	req, err := s.client.NewDeleteRequest(fmt.Sprintf("type/%s/%s", id, picturesPath))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}
	return nil
}

func (p *Picture) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			p.ID = v.(string)
		}

		if strings.ToLower(k) == "imageurl" {
			p.ImageURL = v.(string)
		}
	}

	return nil
}
