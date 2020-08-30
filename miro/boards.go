package miro

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	boardsPath = "boards"
)

// BoardsService handles communication to Miro Boards API.
//
// API doc: https://developers.miro.com/reference#board-object
type BoardsService service

// Board object represents Miro Board.
//
// API doc: https://developers.miro.com/reference#board-object
//go:generate gomodifytags -file $GOFILE -struct Board -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Board -add-tags json -w -transform camelcase
type Board struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Description           string              `json:"description"`
	ImageURL              string              `json:"imageURL"`
	CreatedAt             time.Time           `json:"createdAt"`
	ModifiedAt            time.Time           `json:"modifiedAt"`
	CreatedBy             *MiniUser           `json:"createdBy"`
	ModifiedBy            *MiniUser           `json:"modifiedBy"`
	Owner                 *MiniUser           `json:"owner"`
	Picture               *MiniPicture        `json:"picture"`
	ViewLink              string              `json:"viewLink"`
	SharingPolicy         *SharingPolicy      `json:"sharingPolicy"`
	CurrentUserConnection *TeamUserConnection `json:"currentUserConnection"`
}

// SharingPolicy object represents the policy for the board
//go:generate gomodifytags -file $GOFILE -struct SharingPolicy -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct SharingPolicy -add-tags json -w -transform camelcase
type SharingPolicy struct {
	Access     string `json:"access"`
	TeamAccess string `json:"teamAccess"`
}

//go:generate gomodifytags -file $GOFILE -struct MiniBoard -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct MiniBoard -add-tags json -w -transform camelcase
type MiniBoard struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Description           string              `json:"description"`
	ImageURL              string              `json:"imageURL"`
	CreatedAt             time.Time           `json:"createdAt"`
	ModifiedAt            time.Time           `json:"modifiedAt"`
	CreatedBy             *User               `json:"createdBy"`
	ModifiedBy            *User               `json:"modifiedBy"`
	Owner                 *User               `json:"owner"`
	Picture               *Picture            `json:"picture"`
	ViewLink              string              `json:"viewLink"`
	SharingPolicy         *SharingPolicy      `json:"sharingPolicy"`
	CurrentUserConnection *TeamUserConnection `json:"currentUserConnection"`
}

// Get gets board by Board ID.
//
// API doc: https://developers.miro.com/reference#get-board
func (s *BoardsService) Get(ctx context.Context, id string) (*Board, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s", boardsPath, id))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}

	board := &Board{}
	if err := json.NewDecoder(resp.Body).Decode(board); err != nil {
		return nil, err
	}

	return board, nil
}

// CreateBoardRequest represents create board request payload.
//
//go:generate gomodifytags -file $GOFILE -struct CreateBoardRequest -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct CreateBoardRequest -add-tags json -w -transform camelcase
type CreateBoardRequest struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	SharingPolicy *SharingPolicy `json:"sharingPolicy"`
}

// Create creates board by Board Request.
//
// API doc: https://developers.miro.com/reference#create-board
func (s *BoardsService) Create(ctx context.Context, b *CreateBoardRequest) (*Board, error) {
	req, err := s.client.NewPostRequest(boardsPath, b)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}

	board := &Board{}
	if err := json.NewDecoder(resp.Body).Decode(board); err != nil {
		return nil, err
	}

	return board, nil
}

// UpdateBoardRequest represents update board request payload.
//
//go:generate gomodifytags -file $GOFILE -struct UpdateBoardRequest -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct UpdateBoardRequest -add-tags json -w -transform camelcase
type UpdateBoardRequest struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	SharingPolicy *SharingPolicy `json:"sharingPolicy"`
}

func (s *BoardsService) Update(ctx context.Context, b *UpdateBoardRequest) (*Board, error) {
	req, err := s.client.NewPatchRequest(boardsPath, b)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}

	board := &Board{}
	if err := json.NewDecoder(resp.Body).Decode(board); err != nil {
		return nil, err
	}

	return board, nil
}

// Delete deletes board by Board Request.
//
// API doc: No document yet
func (s *BoardsService) Delete(ctx context.Context, id string) error {
	req, err := s.client.NewDeleteRequest(fmt.Sprintf("%s/%s", boardsPath, id))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return err
		}
		return respErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respErr := &RespError{}
		if err := json.NewDecoder(resp.Body).Decode(respErr); err != nil {
			return err
		}
		return fmt.Errorf("status code not expected, got:%d, message:%s", resp.StatusCode, respErr.Message)
	}

	return nil
}

func (b *Board) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			b.ID = v.(string)
		}

		if strings.ToLower(k) == "name" {
			b.Name = v.(string)
		}

		if strings.ToLower(k) == "description" {
			b.Description = v.(string)
		}

		if strings.ToLower(k) == "viewlink" {
			b.ViewLink = v.(string)
		}

		if strings.ToLower(k) == "createdAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			b.CreatedAt = at
		}

		if strings.ToLower(k) == "modifiedAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			b.ModifiedAt = at
		}

		if strings.ToLower(k) == "owner" {
			user := &MiniUser{}
			u := v.(map[string]interface{})

			for k, v := range u {
				if strings.ToLower(k) == "id" {
					user.ID = v.(string)
				}

				if strings.ToLower(k) == "name" {
					user.Name = v.(string)
				}
			}

			b.Owner = user
		}

		if strings.ToLower(k) == "picture" {
			pic := &MiniPicture{}
			p := v.(map[string]interface{})

			for k, v := range p {
				if strings.ToLower(k) == "id" {
					pic.ID = v.(string)
				}

				if strings.ToLower(k) == "imageurl" {
					pic.ImageURL = v.(string)
				}
			}

			b.Picture = pic
		}

		if strings.ToLower(k) == "createdBy" {
			user := &MiniUser{}
			u := v.(map[string]interface{})

			for k, v := range u {
				if strings.ToLower(k) == "id" {
					user.ID = v.(string)
				}

				if strings.ToLower(k) == "name" {
					user.Name = v.(string)
				}
			}

			b.CreatedBy = user
		}

		if strings.ToLower(k) == "modifiedBy" {
			user := &MiniUser{}
			u := v.(map[string]string)

			for k, v := range u {
				if strings.ToLower(k) == "id" {
					user.ID = v
				}

				if strings.ToLower(k) == "name" {
					user.Name = v
				}
			}

			b.ModifiedBy = user
		}
	}

	return nil
}
