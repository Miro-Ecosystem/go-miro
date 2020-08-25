package miro

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	rateLimitResetHeader     = "X-RateLimit-Reset"
	rateLimitRemainingHeader = "X-RateLimit-Remaining"
	rateLimitLimitHeader     = "X-RateLimit-Limit"
)

type Client struct {
	common service
	client *http.Client

	mu sync.RWMutex

	RateLimit   *RateLimit
	UserAgent   string
	AccessToken string
	BaseURL     *url.URL

	AuditLogs           *AuditLogsService
	AuthzInfo           *AuthzInfoService
	Boards              *BoardsService
	BoardUserConnection *BoardUserConnectionService
	Picture             *PicturesService
	Teams               *TeamsService
	TeamUserConnection  *TeamUserConnectionService
	Users               *UsersService
}

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

func NewClient() *Client {
	baseURL, _ := url.Parse(baseURL)
	c := &Client{
		BaseURL: baseURL,
	}

	if c.UserAgent != "" {
		c.UserAgent = defaultUserAgent
	}

	c.common.client = c
	c.client = http.DefaultClient

	c.AuditLogs = (*AuditLogsService)(&c.common)
	c.AuthzInfo = (*AuthzInfoService)(&c.common)
	c.Boards = (*BoardsService)(&c.common)
	c.BoardUserConnection = (*BoardUserConnectionService)(&c.common)
	c.Picture = (*PicturesService)(&c.common)
	c.Teams = (*TeamsService)(&c.common)
	c.TeamUserConnection = (*TeamUserConnectionService)(&c.common)
	c.Users = (*UsersService)(&c.common)

	return c
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// NewGetRequest creates an API GET request.
func (c *Client) NewGetRequest(urlStr string) (*http.Request, error) {
	return c.NewRequest("GET", urlStr, nil)
}

// NewPOSTRequest creates an API POST request.
func (c *Client) NewPostRequest(urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequest("POST", urlStr, body)
}

// NewPatchRequest creates an API Patch request.
func (c *Client) NewPatchRequest(urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequest("Patch", urlStr, body)
}

// NewDeleteRequest creates an API Delete request.
func (c *Client) NewDeleteRequest(urlStr string) (*http.Request, error) {
	return c.NewRequest("Delete", urlStr, nil)
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if l := resp.Header.Get(rateLimitLimitHeader); l != "" {
		c.mu.Lock()
		limit, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}
		c.RateLimit.Limit = limit
		defer c.mu.Unlock()
	}

	if r := resp.Header.Get(rateLimitRemainingHeader); r != "" {
		c.mu.Lock()
		remaining, err := strconv.Atoi(r)
		if err != nil {
			return nil, err
		}
		c.RateLimit.Remaining = remaining
		defer c.mu.Unlock()
	}

	if r := resp.Header.Get(rateLimitResetHeader); r != "" {
		c.mu.Lock()
		r, err := strconv.Atoi(r)
		if err != nil {
			return nil, err
		}
		c.RateLimit.Reset = time.Unix(int64(r), 0)
		defer c.mu.Unlock()
	}

	return resp, nil
}
