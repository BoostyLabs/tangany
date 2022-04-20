package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// ExecuteEthereumSmartContractRequest fields that describes execute an Ethereum smart contract request.
type ExecuteEthereumSmartContractRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Wallet                 string `json:"wallet"`
	Contract               string `json:"contract"`
	Function               string `json:"function"`
	Amount                 string `json:"amount"`
}

// ExecuteEthereumSmartContractResponse fields that returns from execute an Ethereum smart contract method.
type ExecuteEthereumSmartContractResponse struct {
	StatusURI string `json:"statusUri"`
}

// ExecuteEthereumSmartContract executes an Ethereum smart contract.
func (client *Client) ExecuteEthereumSmartContract(ctx context.Context, execSmartContract ExecuteEthereumSmartContractRequest) (response ExecuteEthereumSmartContractResponse, err error) {
	var requestBody struct {
		Function string `json:"function"`
		Amount   string `json:"amount"`
	}

	requestBody.Function = execSmartContract.Function
	requestBody.Amount = execSmartContract.Amount

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return ExecuteEthereumSmartContractResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"eth/contract/"+execSmartContract.Contract+"/"+execSmartContract.Wallet+"/send-async", bytes.NewReader(jsonBody))
	if err != nil {
		return ExecuteEthereumSmartContractResponse{}, err
	}

	if execSmartContract.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", execSmartContract.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return ExecuteEthereumSmartContractResponse{}, err
	}
	
	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusAccepted {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return ExecuteEthereumSmartContractResponse{}, err
		}

		return ExecuteEthereumSmartContractResponse{}, errs.New(errorResp.Message)
	}

	var executeEthereumSmartContractResponse ExecuteEthereumSmartContractResponse
	if err = json.NewDecoder(resp.Body).Decode(&executeEthereumSmartContractResponse); err != nil {
		return ExecuteEthereumSmartContractResponse{}, err
	}

	return executeEthereumSmartContractResponse, nil
}
