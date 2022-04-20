package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// GetERC20TokenBalanceRequest fields that describes get erc20 token balance request.
type GetERC20TokenBalanceRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Wallet                 string `json:"wallet"`
	Token                  string `json:"token"`
}

// GetERC20TokenBalanceResponse fields that returns from get erc20 token balance method.
type GetERC20TokenBalanceResponse struct {
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

// GetERC20TokenBalance returns erc20 token balance.
func (client *Client) GetERC20TokenBalance(ctx context.Context, erc20Balance GetERC20TokenBalanceRequest) (response GetERC20TokenBalanceResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+"eth/erc20/"+erc20Balance.Token+"/"+erc20Balance.Wallet, nil)
	if err != nil {
		return GetERC20TokenBalanceResponse{}, err
	}

	if erc20Balance.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", erc20Balance.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return GetERC20TokenBalanceResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return GetERC20TokenBalanceResponse{}, err
		}

		return GetERC20TokenBalanceResponse{}, errs.New(errorResp.Message)
	}

	var getERC20TokenBalanceResponse GetERC20TokenBalanceResponse
	if err = json.NewDecoder(resp.Body).Decode(&getERC20TokenBalanceResponse); err != nil {
		return GetERC20TokenBalanceResponse{}, err
	}

	return getERC20TokenBalanceResponse, nil
}
