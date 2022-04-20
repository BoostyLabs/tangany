package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// CallSmartContractMethodRequest fields that describes calls smart contract method request.
type CallSmartContractMethodRequest struct {
	TanganyEthereumNetwork string   `json:"tangany-ethereum-network"`
	Contract               string   `json:"contract"`
	Function               string   `json:"function"`
	Outputs                []string `json:"outputs"`
}

// CallSmartContractMethodResponse fields that returns from call smart contract method response.
type CallSmartContractMethodResponse struct {
	List []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"list"`
}

// CallSmartContractMethod calls smart contract method.
func (client *Client) CallSmartContractMethod(ctx context.Context, callSmartContract CallSmartContractMethodRequest) (response CallSmartContractMethodResponse, err error) {
	var requestBody struct {
		Function string   `json:"function"`
		Outputs  []string `json:"outputs"`
	}

	requestBody.Function = callSmartContract.Function
	requestBody.Outputs = callSmartContract.Outputs

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return CallSmartContractMethodResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"eth/contract/"+callSmartContract.Contract+"/call", bytes.NewReader(jsonBody))
	if err != nil {
		return CallSmartContractMethodResponse{}, err
	}

	if callSmartContract.TanganyEthereumNetwork != "" {
		req.Header.Add("tangany-ethereum-network", callSmartContract.TanganyEthereumNetwork)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return CallSmartContractMethodResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return CallSmartContractMethodResponse{}, err
		}

		return CallSmartContractMethodResponse{}, errs.New(errorResp.Message)
	}

	var callSmartContractMethodResponse CallSmartContractMethodResponse
	if err = json.NewDecoder(resp.Body).Decode(&callSmartContractMethodResponse); err != nil {
		return CallSmartContractMethodResponse{}, err
	}

	return callSmartContractMethodResponse, nil
}
