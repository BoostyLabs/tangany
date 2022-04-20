package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// EstimateTransactionFeeRequest fields that describes estimate transaction fee request.
type EstimateTransactionFeeRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Wallet                 string `json:"wallet"`
	Amount                 string `json:"amount"`
	To                     string `json:"to"`
	Data                   string `json:"data"`
}

// EstimateTransactionFeeResponse fields that returns from estimate transaction fee method.
type EstimateTransactionFeeResponse struct {
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Fee      string `json:"fee"`
}

// EstimateTransactionFee estimates transaction fee.
func (client *Client) EstimateTransactionFee(ctx context.Context, estimateFee EstimateTransactionFeeRequest) (response EstimateTransactionFeeResponse, err error) {
	var requestBody struct {
		Amount string `json:"amount"`
		To     string `json:"to"`
		Data   string `json:"data"`
	}

	requestBody.To = estimateFee.To
	requestBody.Data = estimateFee.Data
	requestBody.Amount = estimateFee.Amount

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return EstimateTransactionFeeResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"eth/wallet/"+estimateFee.Wallet+"/estimate-fee", bytes.NewReader(jsonBody))
	if err != nil {
		return EstimateTransactionFeeResponse{}, err
	}

	if estimateFee.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", estimateFee.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return EstimateTransactionFeeResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return EstimateTransactionFeeResponse{}, err
		}

		return EstimateTransactionFeeResponse{}, errs.New(errorResp.Message)
	}

	var estimateTransactionFeeResponse EstimateTransactionFeeResponse
	if err = json.NewDecoder(resp.Body).Decode(&estimateTransactionFeeResponse); err != nil {
		return EstimateTransactionFeeResponse{}, err
	}

	return estimateTransactionFeeResponse, nil
}
