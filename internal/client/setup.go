package client

import "net/http"

func New() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}
