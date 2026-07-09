package graph

import (
	"net/http"
	"time"
)

const baseURL = "https://graph.microsoft.com/v1.0"

type Client struct {
	auth       *Auth
	httpClient *http.Client
	userEmail  string

	fileCounter int
}

func NewClient(auth *Auth, userEmail string) *Client {
	return &Client{
		auth: auth,
		userEmail: userEmail,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}