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
	auditLogsPath = "audit/logs"
)

// AuditLogsService handles communication to Miro Logs API.
//
// API doc: https://developers.miro.com/reference#log-object
type AuditLogsService service

// AuditLog object represents Miro Log.
//
// API doc: https://developers.miro.com/reference#log-object
//go:generate gomodifytags -file $GOFILE -struct AuditLog -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct AuditLog -add-tags json -w -transform camelcase
type AuditLog struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Size     int    `json:"size"`
	NextLink string `json:"nextLink"`
	PrevLink string `json:"prevLink"`
	Data     []Data `json:"data"`
}

//go:generate gomodifytags -file $GOFILE -struct Data -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Data -add-tags json -w -transform camelcase
type Data struct {
	ID        string    `json:"id"`
	Event     string    `json:"event"`
	Details   *Detail   `json:"details"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *MiniUser `json:"createdBy"`
	Context   *Context  `json:"context"`
}

//go:generate gomodifytags -file $GOFILE -struct Context -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Context -add-tags json -w -transform camelcase
type Context struct {
	Organization *Organization `json:"organization"`
	Team         *MiniTeam     `json:"team"`
	IP           string        `json:"ip"`
}

//go:generate gomodifytags -file $GOFILE -struct Organization -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Organization -add-tags json -w -transform camelcase
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//go:generate gomodifytags -file $GOFILE -struct Detail -clear-tags -w
//go:generate gomodifytags --file $GOFILE --struct Detail -add-tags json -w -transform camelcase
type Detail struct {
	Role string `json:"role"`
}

// Get gets logs by condition.
//
// API doc: https://developers.miro.com/reference#get-logs
func (s *AuditLogsService) Get(ctx context.Context) (*AuditLog, error) {
	req, err := s.client.NewGetRequest(auditLogsPath)
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

	l := &AuditLog{}
	if err := json.NewDecoder(resp.Body).Decode(l); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *AuditLog) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "limit" {
			l.Limit = int(v.(float64))
		}

		if strings.ToLower(k) == "offset" {
			l.Offset = int(v.(float64))
		}

		if strings.ToLower(k) == "size" {
			l.Size = int(v.(float64))
		}

		if strings.ToLower(k) == "nextlink" {
			l.NextLink = v.(string)
		}

		if strings.ToLower(k) == "prevlink" {
			l.PrevLink = v.(string)
		}

		if strings.ToLower(k) == "data" {
			data, err := unmarshalData(v.([]interface{}))
			if err != nil {
				return err
			}
			l.Data = data
		}
	}

	return nil
}

func unmarshalData(d []interface{}) ([]Data, error) {
	data := make([]Data, len(d))
	for i, v := range d {
		datum := Data{}
		rawDatum := v.(map[string]interface{})

		for k, v := range rawDatum {
			if strings.ToLower(k) == "id" {
				datum.ID = v.(string)
			}

			if strings.ToLower(k) == "event" {
				datum.Event = v.(string)
			}

			if strings.ToLower(k) == "createdAt" {
				at, err := time.Parse(time.RFC3339, v.(string))
				if err != nil {
					return []Data{}, err
				}
				datum.CreatedAt = at
			}
		}

		data[i] = datum
	}

	return data, nil
}
