package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// MakeAsyncWalletTransactionRequest fields that describes async wallet transaction request.
type MakeAsyncWalletTransactionRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Wallet                 string `json:"wallet"`
	Amount                 string `json:"amount"`
	To                     string `json:"to"`
	Data                   string `json:"data"`
}

// MakeAsyncWalletTransactionResponse fields that returns from async wallet transaction method.
type MakeAsyncWalletTransactionResponse struct {
	StatusURI string `json:"statusUri"`
}

// MakeAsyncWalletTransaction makes asynchronous wallet transaction.
func (client *Client) MakeAsyncWalletTransaction(ctx context.Context, asyncWalletTransaction MakeAsyncWalletTransactionRequest) (response MakeAsyncWalletTransactionResponse, err error) {
	var requestBody struct {
		Amount string `json:"amount"`
		To     string `json:"to"`
		Data   string `json:"data"`
	}

	requestBody.To = asyncWalletTransaction.To
	requestBody.Data = asyncWalletTransaction.Data
	requestBody.Amount = asyncWalletTransaction.Amount

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return MakeAsyncWalletTransactionResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"eth/wallet/"+asyncWalletTransaction.Wallet+"/send-async", bytes.NewReader(jsonBody))
	if err != nil {
		return MakeAsyncWalletTransactionResponse{}, err
	}

	if asyncWalletTransaction.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", asyncWalletTransaction.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return MakeAsyncWalletTransactionResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return MakeAsyncWalletTransactionResponse{}, err
		}

		return MakeAsyncWalletTransactionResponse{}, errs.New(errorResp.Message)
	}

	var makeAsyncWalletTransactionResponse MakeAsyncWalletTransactionResponse
	if err = json.NewDecoder(resp.Body).Decode(&makeAsyncWalletTransactionResponse); err != nil {
		return MakeAsyncWalletTransactionResponse{}, err
	}

	return makeAsyncWalletTransactionResponse, nil
}
