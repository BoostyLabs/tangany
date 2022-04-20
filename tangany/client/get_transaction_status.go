package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// GetTransactionStatusRequest fields that describes get transaction status request.
type GetTransactionStatusRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Hash                   string `json:"hash"`
}

// GetTransactionStatusResponse fields that returns from get transaction status method.
type GetTransactionStatusResponse struct {
	BlockNr          int    `json:"blockNr"`
	Data             string `json:"data"`
	IsError          bool   `json:"isError"`
	Status           string `json:"status"`
	Confirmations    int    `json:"confirmations"`
	From             string `json:"from"`
	To               string `json:"to"`
	ContractAddress  string `json:"contractAddress"`
	GasPrice         string `json:"gasPrice"`
	Gas              int    `json:"gas"`
	GasUsed          int    `json:"gasUsed"`
	Nonce            int    `json:"nonce"`
	Value            string `json:"value"`
	Timestamp        int    `json:"timestamp"`
	TransactionIndex int    `json:"transactionIndex"`
}

// GetTransactionStatus returns transaction status.
func (client *Client) GetTransactionStatus(ctx context.Context, transactionStatus GetTransactionStatusRequest) (response GetTransactionStatusResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+"eth/transaction/"+transactionStatus.Hash, nil)
	if err != nil {
		return GetTransactionStatusResponse{}, err
	}

	if transactionStatus.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", transactionStatus.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return GetTransactionStatusResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return GetTransactionStatusResponse{}, err
		}

		return GetTransactionStatusResponse{}, errs.New(errorResp.Message)
	}

	var getTransactionStatusResponse GetTransactionStatusResponse
	if err = json.NewDecoder(resp.Body).Decode(&getTransactionStatusResponse); err != nil {
		return GetTransactionStatusResponse{}, err
	}

	return getTransactionStatusResponse, nil
}
