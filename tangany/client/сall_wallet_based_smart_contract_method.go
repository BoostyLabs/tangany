package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// CallWalletBasedSmartContractMethodRequest fields that describes call wallet based smart contract method request.
type CallWalletBasedSmartContractMethodRequest struct {
	TanganyEthereumNetwork string   `json:"tangany-ethereum-network"`
	Wallet                 string   `json:"wallet"`
	Contract               string   `json:"contract"`
	Function               string   `json:"function"`
	Inputs                 []string `json:"inputs"`
	Outputs                []string `json:"outputs"`
}

// CallWalletBasedSmartContractMethodResponse fields that returns from call wallet based smart contract method method.
type CallWalletBasedSmartContractMethodResponse struct {
	List []struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	} `json:"list"`
}

// CallWalletBasedSmartContractMethod calls wallet based smart contract.
func (client *Client) CallWalletBasedSmartContractMethod(ctx context.Context, callSmartContract CallWalletBasedSmartContractMethodRequest) (response CallWalletBasedSmartContractMethodResponse, err error) {
	var requestBody struct {
		Function string   `json:"function"`
		Inputs   []string `json:"inputs"`
		Outputs  []string `json:"outputs"`
	}

	requestBody.Function = callSmartContract.Function
	requestBody.Inputs = callSmartContract.Inputs
	requestBody.Outputs = callSmartContract.Outputs

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return CallWalletBasedSmartContractMethodResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"eth/contract/"+callSmartContract.Contract+"/"+callSmartContract.Wallet+"/call", bytes.NewReader(jsonBody))
	if err != nil {
		return CallWalletBasedSmartContractMethodResponse{}, err
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
		return CallWalletBasedSmartContractMethodResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return CallWalletBasedSmartContractMethodResponse{}, err
		}

		return CallWalletBasedSmartContractMethodResponse{}, errs.New(errorResp.Message)
	}

	var callWalletBasedSmartContractMethodResponse CallWalletBasedSmartContractMethodResponse
	if err = json.NewDecoder(resp.Body).Decode(&callWalletBasedSmartContractMethodResponse); err != nil {
		return CallWalletBasedSmartContractMethodResponse{}, err
	}

	return callWalletBasedSmartContractMethodResponse, nil
}
