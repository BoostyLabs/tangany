package client

import (
	"net/http"
)

// Config is the global configuration for tangany client.
type Config struct {
	TanganyDefaultURL   string `json:"tanganyDefaultURL" default:"https://api.tangany.com/v1/"`
	TanganyClientID     string `json:"tangany-client-id"`
	TanganyClientSecret string `json:"tangany-client-secret"`
	TanganyVaultURL     string `json:"tangany-vault-url"`
	TanganySubscription string `json:"tangany-subscription"`
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
