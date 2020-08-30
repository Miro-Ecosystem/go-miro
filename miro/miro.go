package miro

const (
	baseURL          = "https://api.miro.com"
	defaultUserAgent = "go-miro"
)

type service struct {
	client *Client
}
