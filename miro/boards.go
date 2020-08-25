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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	board := &Board{}
	if err := json.NewDecoder(resp.Body).Decode(board); err != nil {
		return nil, err
	}

	return board, nil
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
