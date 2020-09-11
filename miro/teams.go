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
	teamsPath           = "teams"
	userConnectionsPath = "user-connections"
	teamInvitePath      = "invite"
)

// TeamsService handles communication to Miro Teams API.
//
// API doc: https://developers.miro.com/reference#team-object
type TeamsService service

// Team object represents Miro Team.
//
// API doc: https://developers.miro.com/reference#team-object
//go:generate gomodifytags -file $GOFILE -struct Team -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Team -add-tags json -w -transform camelcase
type Team struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	CreatedAt  time.Time    `json:"createdAt"`
	ModifiedAt time.Time    `json:"modifiedAt"`
	CreatedBy  *MiniUser    `json:"createdBy"`
	ModifiedBy *MiniUser    `json:"modifiedBy"`
	Picture    *MiniPicture `json:"picture"`
}

func (t *Team) GetType() string {
	return "team"
}

type MiniTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t *MiniTeam) GetType() string {
	return "team"
}

// Get gets team by Team ID.
//
// API doc: https://developers.miro.com/reference#get-team
func (s *TeamsService) Get(ctx context.Context, id string) (*Team, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s", teamsPath, id))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if http.StatusBadRequest <= resp.StatusCode && resp.StatusCode <= http.StatusInsufficientStorage {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	t := &Team{}
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

// UpdateTeamRequest represents request to update team user connection
//
//go:generate gomodifytags -file $GOFILE -struct UpdateTeamRequest -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct UpdateTeamRequest -add-tags json -w -transform camelcase
type UpdateTeamRequest struct {
	Name string `json:"name"`
}

// Update updates team by Team ID.
//
// API doc: https://developers.miro.com/reference#update-team
func (s *TeamsService) Update(ctx context.Context, id string, request *UpdateTeamRequest) (*Team, error) {
	req, err := s.client.NewPatchRequest(fmt.Sprintf("%s/%s", teamsPath, id), request)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if http.StatusBadRequest <= resp.StatusCode && resp.StatusCode <= http.StatusInsufficientStorage {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	t := &Team{}
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

// ListTeamMembersResponse represents list response from Miro
//
//go:generate gomodifytags -file $GOFILE -struct ListTeamMembersResponse -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct ListTeamMembersResponse -add-tags json -w -transform camelcase
type ListTeamMembersResponse struct {
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Size   int     `json:"size"`
	Data   []*Team `json:"data"`
}

// ListTeamMembers gets all team members
//
// API doc: https://developers.miro.com/reference#get-team-user-connections
func (s *TeamsService) ListTeamMembers(ctx context.Context, id string) (*ListTeamMembersResponse, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s", teamsPath, id))
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if http.StatusBadRequest <= resp.StatusCode && resp.StatusCode <= http.StatusInsufficientStorage {
		return nil, fmt.Errorf("status code not expected, got:%d", resp.StatusCode)
	}

	t := &ListTeamMembersResponse{}
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

// GetCurrentUserConnection gets team current user connection by Team ID.
//
// API doc: https://developers.miro.com/reference#get-team-current-user-connection
func (s *TeamsService) GetCurrentUserConnection(ctx context.Context, id string) (*TeamUserConnection, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s/%s/me", teamsPath, id, userConnectionsPath))
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

// Invite invites passed user to specified team.
//
// API doc: https://developers.miro.com/reference#invite-to-team
func (s *TeamsService) Invite(ctx context.Context, id string, email string) ([]*TeamUserConnection, error) {
	req, err := s.client.NewPostRequest(fmt.Sprintf("%s/%s/%s/%s?email=%s", teamsPath, id, userConnectionsPath, teamInvitePath, email), nil)
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

	conns := []*TeamUserConnection{}
	if err := json.NewDecoder(resp.Body).Decode(&conns); err != nil {
		return nil, err
	}

	return conns, nil
}

func (t *Team) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			t.ID = v
		}
		if strings.ToLower(k) == "name" {
			t.Name = v
		}

		if strings.ToLower(k) == "createdAt" {
			at, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
			t.CreatedAt = at
		}

		if strings.ToLower(k) == "modifiedAt" {
			at, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
			t.ModifiedAt = at
		}

		if strings.ToLower(k) == "createdBy" {
			u := &MiniUser{}
			if err := json.Unmarshal([]byte(v), u); err != nil {
				return err
			}

			t.CreatedBy = u
		}

		if strings.ToLower(k) == "modifiedBy" {
			u := &MiniUser{}
			if err := json.Unmarshal([]byte(v), u); err != nil {
				return err
			}

			t.ModifiedBy = u
		}

		if strings.ToLower(k) == "picture" {
			p := &MiniPicture{}
			if err := json.Unmarshal([]byte(v), p); err != nil {
				return err
			}

			t.Picture = p
		}
	}

	return nil
}
