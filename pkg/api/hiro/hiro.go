package hiro

import (
	"net/http"
)

const DefaultApiBase = "https://stacks.hashhavoc.com"

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}
