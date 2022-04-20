package client

import (
	"net/http"
)

// Config is the global configuration for tangany client.
type Config struct {
	TanganyDefaultURL   string `json:"tanganyDefaultURL" default:"https://api.tangany.com/v1/"`
	TanganyClientID     string `json:"tangany-client-id" default:"9f78a896-80e9-4e86-b51e-9c7d1e9ebc0f"`
	TanganyClientSecret string `json:"tangany-client-secret" default:"PQK7Q~n.SFxL.KXipI0h8Y3ruk5USTgFzHUDS"`
	TanganyVaultURL     string `json:"tangany-vault-url" default:"https://cw-keyv-demo-finn.vault.azure.net"`
	TanganySubscription string `json:"tangany-subscription" default:"845cb70213614a9abb95fe3c09d406df"`
}

// Client implementation of the tangany API.
type Client struct {
	http   http.Client
	config Config
}

// NewClient is a constructor for tangany API client.
func NewClient(config Config) *Client {
	return &Client{config: config}
}

// ErrorResponse holds fields that explaining error from Tangany API.
type ErrorResponse struct {
	StatusCode int64  `json:"statusCode"`
	Message    string `json:"message"`
	ActivityId string `json:"activityId"`
}

type Tag map[string]interface{}
