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
	teamUserConnectionsPath = "team-user-connection"
)

// TeamUserConnectionsService handles communication to Miro TeamUserConnections API.
//
// API doc: https://developers.miro.com/reference#team-user-connection-object
type TeamUserConnectionService service

// TeamUserConnection object represents Miro TeamUserConnection.
//
// API doc: https://developers.miro.com/reference#team-user-connection-object
//go:generate gomodifytags -file $GOFILE -struct TeamUserConnection -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct TeamUserConnection -add-tags json -w -transform camelcase
type TeamUserConnection struct {
	ID         string    `json:"id"`
	User       *MiniUser `json:"user"`
	Team       *MiniTeam `json:"team"`
	Role       string    `json:"role"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	CreatedBy  *MiniUser `json:"createdBy"`
	ModifiedBy *MiniUser `json:"modifiedBy"`
}

//go:generate gomodifytags -file $GOFILE -struct MiniTeamUserConnection -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct MiniTeamUserConnection -add-tags json -w -transform camelcase
type MiniTeamUserConnection struct {
	ID   string    `json:"id"`
	User *MiniUser `json:"user"`
	Role string    `json:"role"`
}

// Get gets team user connection by TeamUserConnection ID.
//
// API doc: https://developers.miro.com/reference#get-team-user-connection
func (s *TeamUserConnectionService) Get(ctx context.Context, id string) (*TeamUserConnection, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s", teamUserConnectionsPath, id))
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

	conn := &TeamUserConnection{}
	if err := json.NewDecoder(resp.Body).Decode(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func (t *TeamUserConnection) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			t.ID = v.(string)
		}

		if strings.ToLower(k) == "name" {
			t.Name = v.(string)
		}

		if strings.ToLower(k) == "role" {
			t.Role = v.(string)
		}

		if strings.ToLower(k) == "createdAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			t.CreatedAt = at
		}

		if strings.ToLower(k) == "modifiedAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			t.ModifiedAt = at
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

			t.User = user
		}

		if strings.ToLower(k) == "team" {
			team := &MiniTeam{}
			u := v.(map[string]interface{})

			for k, v := range u {
				if strings.ToLower(k) == "id" {
					team.ID = v.(string)
				}

				if strings.ToLower(k) == "name" {
					team.Name = v.(string)
				}
			}

			t.Team = team
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

			t.CreatedBy = user
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

			t.ModifiedBy = user
		}
	}

	return nil
}

func (t *MiniTeamUserConnection) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			t.ID = v
		}

		if strings.ToLower(k) == "role" {
			t.Role = v
		}

		if strings.ToLower(k) == "user" {
			u := &MiniUser{}
			if err := json.Unmarshal([]byte(v), u); err != nil {
				return err
			}

			t.User = u
		}
	}

	return nil
}
