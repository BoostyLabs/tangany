package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// EstimateEthereumSmartContractRequest fields that describes estimate an Ethereum smart contract request.
type EstimateEthereumSmartContractRequest struct {
	TanganyEthereumNetwork string   `json:"tangany-ethereum-network"`
	Wallet                 string   `json:"wallet"`
	Contract               string   `json:"contract"`
	Function               string   `json:"function"`
	Inputs                 []string `json:"inputs"`
}

// EstimateEthereumSmartContractResponse fields that returns from estimate an Ethereum smart contract method.
type EstimateEthereumSmartContractResponse struct {
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Fee      string `json:"fee"`
}

// EstimateEthereumSmartContract estimates an Ethereum smart contract.
func (client *Client) EstimateEthereumSmartContract(ctx context.Context, estimateSmartContract EstimateEthereumSmartContractRequest) (response EstimateEthereumSmartContractResponse, err error) {
	var requestBody struct {
		Function string   `json:"function"`
		Inputs   []string `json:"inputs"`
	}

	requestBody.Function = estimateSmartContract.Function
	requestBody.Inputs = estimateSmartContract.Inputs

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return EstimateEthereumSmartContractResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"eth/contract/"+estimateSmartContract.Contract+"/"+estimateSmartContract.Wallet+"/estimate-fee", bytes.NewReader(jsonBody))
	if err != nil {
		return EstimateEthereumSmartContractResponse{}, err
	}

	if estimateSmartContract.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", estimateSmartContract.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return EstimateEthereumSmartContractResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return EstimateEthereumSmartContractResponse{}, err
		}

		return EstimateEthereumSmartContractResponse{}, errs.New(errorResp.Message)
	}

	var estimateEthereumSmartContractResponse EstimateEthereumSmartContractResponse
	if err = json.NewDecoder(resp.Body).Decode(&estimateEthereumSmartContractResponse); err != nil {
		return EstimateEthereumSmartContractResponse{}, err
	}

	return estimateEthereumSmartContractResponse, nil
}
