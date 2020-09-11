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
	boardUserConnectionsPath = "board-user-connection"
)

// BoardUserConnectionsService handles communication to Miro BoardUserConnections API.
//
// API doc: https://developers.miro.com/reference#board-user-connection-object
type BoardUserConnectionService service

// BoardUserConnection object represents Miro BoardUserConnection.
//
// API doc: https://developers.miro.com/reference#board-user-connection-object
//go:generate gomodifytags -file $GOFILE -struct BoardUserConnection -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct BoardUserConnection -add-tags json -w -transform camelcase
type BoardUserConnection struct {
	ID         string    `json:"id"`
	User       *MiniUser `json:"user"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	CreatedBy  *MiniUser `json:"createdBy"`
	ModifiedBy *MiniUser `json:"modifiedBy"`
}

// Get gets board user connection by BoardUserConnection ID.
//
// API doc: https://developers.miro.com/reference#get-board-user-connection
func (s *BoardUserConnectionService) Get(ctx context.Context, id string) (*BoardUserConnection, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s", boardUserConnectionsPath, id))
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

	conn := &BoardUserConnection{}
	if err := json.NewDecoder(resp.Body).Decode(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// UpdateBoardUserConnectionRequest represents request to update board user connection
//
//go:generate gomodifytags -file $GOFILE -struct UpdateBoardUserConnectionRequest -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct UpdateBoardUserConnectionRequest -add-tags json -w -transform camelcase
type UpdateBoardUserConnectionRequest struct {
	Role string `json:"role"`
}

// Update updates board user connection by BoardUserConnection ID.
//
// API doc: https://developers.miro.com/reference#update-board-user-connection
func (s *BoardUserConnectionService) Updates(ctx context.Context, id string, request *UpdateBoardUserConnectionRequest) (*BoardUserConnection, error) {
	req, err := s.client.NewPatchRequest(fmt.Sprintf("%s/%s", boardUserConnectionsPath, id), request)
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

	conn := &BoardUserConnection{}
	if err := json.NewDecoder(resp.Body).Decode(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// Delete deletes board user connection by BoardUserConnection ID.
//
// API doc: https://developers.miro.com/reference#delete-board-user-connection
func (s *BoardUserConnectionService) Delete(ctx context.Context, id string) error {
	req, err := s.client.NewDeleteRequest(fmt.Sprintf("%s/%s", boardUserConnectionsPath, id))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	conn := &BoardUserConnection{}
	if err := json.NewDecoder(resp.Body).Decode(conn); err != nil {
		return err
	}

	return nil
}

func (c *BoardUserConnection) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			c.ID = v.(string)
		}

		if strings.ToLower(k) == "role" {
			c.Role = v.(string)
		}

		if strings.ToLower(k) == "createdAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			c.CreatedAt = at
		}

		if strings.ToLower(k) == "modifiedAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			c.ModifiedAt = at
		}

		if strings.ToLower(k) == "user" {
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

			c.User = user
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

			c.CreatedBy = user
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

			c.ModifiedBy = user
		}
	}

	return nil
}
