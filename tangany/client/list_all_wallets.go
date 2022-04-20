package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeebo/errs"
)

// ListAllWalletsRequest fields that describes list all wallets request sends to Tangany API.
type ListAllWalletsRequest struct {
	Index int
	Limit int
	Sort  string
}

// ListAllWalletsResponse fields that returns from list all wallets method.
type ListAllWalletsResponse struct {
	Hits struct {
		Total int `json:"total"`
		Hsm   int `json:"hsm"`
	} `json:"hits"`
	List []struct {
		Wallet string `json:"wallet"`
		Links  []struct {
			Href string `json:"href"`
			Type string `json:"type"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"list"`
	Links struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"links"`
}

// ListAllWallets returns list of wallets by passed parameters.
func (client *Client) ListAllWallets(ctx context.Context, list ListAllWalletsRequest) (response ListAllWalletsResponse, err error) {
	rout := createLink(list)

	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+"wallets"+rout, nil)
	if err != nil {
		return ListAllWalletsResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return ListAllWalletsResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return ListAllWalletsResponse{}, err
		}

		return ListAllWalletsResponse{}, errs.New(errorResp.Message)
	}

	var listAllWalletsResponse ListAllWalletsResponse
	if err = json.NewDecoder(resp.Body).Decode(&listAllWalletsResponse); err != nil {
		return ListAllWalletsResponse{}, err
	}

	return listAllWalletsResponse, nil
}

// forms link by listed params.
func createLink(list ListAllWalletsRequest) (route string) {
	index := strconv.Itoa(list.Index)
	limit := strconv.Itoa(list.Limit)

	route = "?index=" + index + "&limit=" + limit + "&sort=" + list.Sort
	return route
}
