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
	AuthorizationInfoPath = "oauth-token"
)

// AuthzService handles communication to Miro BoardUserConnections API.
//
// API doc: https://developers.miro.com/reference#authorization-object
type AuthzInfoService service

// AuthorizationInfo object represents Miro Authorization info.
//
// API doc: https://developers.miro.com/reference#authorization-object
//go:generate gomodifytags -file $GOFILE -struct AuthorizationInfo -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct AuthorizationInfo -add-tags json -w -transform camelcase
type AuthorizationInfo struct {
	ID        string    `json:"id"`
	Scopes    []string  `json:"scopes"`
	User      *MiniUser `json:"user"`
	Team      *MiniTeam `json:"team"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *MiniUser `json:"createdBy"`
}

// Get gets OAuth token.
//
// API doc: https://developers.miro.com/reference#get-authorization
func (s *AuthzInfoService) Get(ctx context.Context) (*AuthorizationInfo, error) {
	req, err := s.client.NewGetRequest(AuthorizationInfoPath)
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

	i := &AuthorizationInfo{}
	if err := json.NewDecoder(resp.Body).Decode(i); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *AuthorizationInfo) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			i.ID = v.(string)
		}

		if strings.ToLower(k) == "createdAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			i.CreatedAt = at
		}

		if strings.ToLower(k) == "scopes" {
			input := v.([]interface{})
			scopes := make([]string, len(input))
			for i, scope := range input {
				scopes[i] = scope.(string)
			}

			i.Scopes = scopes
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

			i.CreatedBy = user
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

			i.Team = team
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

			i.CreatedBy = user
		}
	}

	return nil
}
