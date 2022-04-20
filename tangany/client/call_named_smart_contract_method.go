package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// CallNamedSmartContractMethodRequest fields that describes calls named smart contract method request.
type CallNamedSmartContractMethodRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Contract               string `json:"contract"`
	Method                 string `json:"method"`
	Types                  []string
}

// CallNamedSmartContractMethodResponse fields that returns from call named smart contract method response.
type CallNamedSmartContractMethodResponse struct {
	List []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"list"`
}

// CallNamedSmartContractMethod calls named smart contract method.
func (client *Client) CallNamedSmartContractMethod(ctx context.Context, callSmartContract CallNamedSmartContractMethodRequest) (response CallNamedSmartContractMethodResponse, err error) {
	route := generateRoute(callSmartContract)

	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+route, nil)
	if err != nil {
		return CallNamedSmartContractMethodResponse{}, err
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
		return CallNamedSmartContractMethodResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return CallNamedSmartContractMethodResponse{}, err
		}

		return CallNamedSmartContractMethodResponse{}, errs.New(errorResp.Message)
	}

	var callNamedSmartContractMethodResponse CallNamedSmartContractMethodResponse
	if err = json.NewDecoder(resp.Body).Decode(&callNamedSmartContractMethodResponse); err != nil {
		return CallNamedSmartContractMethodResponse{}, err
	}

	return callNamedSmartContractMethodResponse, nil
}

func generateRoute(req CallNamedSmartContractMethodRequest) (res string) {
	res = "eth/contract/" + req.Contract + "/call/" + req.Method
	if len(req.Types) == 0 {
		return res
	}

	res = res + "?type="
	for i, v := range req.Types {
		if i != len(req.Types)-1 {
			res = res + v + "&type="
			continue
		}

		res = res + v
	}

	return res
}
