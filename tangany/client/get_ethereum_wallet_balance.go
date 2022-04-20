package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// GetEthWalletBalanceRequest fields that describes get eth wallet balance request.
type GetEthWalletBalanceRequest struct {
	Wallet          string `json:"wallet"`
	EthereumNetwork string
}

// GetEthWalletBalanceResponse fields that returns from get eth wallet balance method.
type GetEthWalletBalanceResponse struct {
	Address  string `json:"address"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

// GetEthWalletBalance returns eth wallet balance.
func (client *Client) GetEthWalletBalance(ctx context.Context, getWallet GetEthWalletBalanceRequest) (response GetEthWalletBalanceResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+"eth/wallet/"+getWallet.Wallet, nil)
	if err != nil {
		return GetEthWalletBalanceResponse{}, err
	}

	if getWallet.EthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", getWallet.EthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return GetEthWalletBalanceResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return GetEthWalletBalanceResponse{}, err
		}

		return GetEthWalletBalanceResponse{}, errs.New(errorResp.Message)
	}

	var getEthWalletBalanceResponse GetEthWalletBalanceResponse
	if err = json.NewDecoder(resp.Body).Decode(&getEthWalletBalanceResponse); err != nil {
		return GetEthWalletBalanceResponse{}, err
	}

	return getEthWalletBalanceResponse, nil
}
