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
	usersPath = "users"
)

// UsersService handles communication to Miro Users API.
//
// API doc: https://developers.miro.com/reference#user-object
type UsersService service

// User object represents Miro User.
//
// API doc: https://developers.miro.com/reference#user-object
//go:generate gomodifytags -file $GOFILE -struct User -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct User -add-tags json -w -transform camelcase
type User struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Company   string       `json:"company"`
	Role      string       `json:"role"`
	Industry  string       `json:"industry"`
	Email     string       `json:"email"`
	State     string       `json:"state"`
	CreatedAt time.Time    `json:"createdAt"`
	Picture   *MiniPicture `json:"picture"`
}

// MiniUser is omitted user.
//go:generate gomodifytags -file $GOFILE -struct MiniUser -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct MiniUser -add-tags json -w -transform camelcase
type MiniUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get gets user by User ID.
//
// API doc: https://developers.miro.com/reference#get-user
func (s *UsersService) Get(ctx context.Context, id string) (*User, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/%s", usersPath, id))
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

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetCurrentUser gets current user
//
// API doc: https://developers.miro.com/reference#get-current-user
func (s *UsersService) GetCurrentUser(ctx context.Context) (*User, error) {
	req, err := s.client.NewGetRequest(fmt.Sprintf("%s/me", usersPath))
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

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			u.ID = v.(string)
		}

		if strings.ToLower(k) == "name" {
			u.Name = v.(string)
		}

		if strings.ToLower(k) == "role" {
			u.Role = v.(string)
		}

		if strings.ToLower(k) == "email" {
			u.Email = v.(string)
		}

		if strings.ToLower(k) == "company" {
			u.Company = v.(string)
		}

		if strings.ToLower(k) == "industry" {
			u.Industry = v.(string)
		}

		if strings.ToLower(k) == "state" {
			u.State = v.(string)
		}

		if strings.ToLower(k) == "createdAt" {
			at, err := time.Parse(time.RFC3339, v.(string))
			if err != nil {
				return err
			}
			u.CreatedAt = at
		}

		if strings.ToLower(k) == "picture" {
			pic := &MiniPicture{}
			p := v.(map[string]interface{})

			for k, v := range p {
				if strings.ToLower(k) == "id" {
					pic.ID = v.(string)
				}

				if strings.ToLower(k) == "imageURL" {
					pic.ImageURL = v.(string)
				}
			}

			u.Picture = pic
		}

	}

	return nil
}
