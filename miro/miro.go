package miro

const (
	baseURL          = "https://api.miro.com/v1"
	defaultUserAgent = "go-miro"
)

type service struct {
	client *Client
}
