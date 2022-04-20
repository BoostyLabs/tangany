package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// CallWalletBasedNamedSmartContractMethodRequest fields that describes call wallet based named smart contract method request.
type CallWalletBasedNamedSmartContractMethodRequest struct {
	TanganyEthereumNetwork string `json:"tangany-ethereum-network"`
	Wallet                 string `json:"wallet"`
	Contract               string `json:"contract"`
	Method                 string `json:"method"`
	Types                  []string
}

// CallWalletBasedNamedSmartContractMethodResponse fields that returns from call wallet based named smart contract method method.
type CallWalletBasedNamedSmartContractMethodResponse struct {
	List []struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	} `json:"list"`
}

// CallWalletBasedNamedSmartContractMethod calls wallet based named smart contract.
func (client *Client) CallWalletBasedNamedSmartContractMethod(ctx context.Context, callSmartContract CallWalletBasedNamedSmartContractMethodRequest) (response CallWalletBasedNamedSmartContractMethodResponse, err error) {
	route := makeRoute(callSmartContract)

	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+route, nil)
	if err != nil {
		return CallWalletBasedNamedSmartContractMethodResponse{}, err
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
		return CallWalletBasedNamedSmartContractMethodResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return CallWalletBasedNamedSmartContractMethodResponse{}, err
		}

		return CallWalletBasedNamedSmartContractMethodResponse{}, errs.New(errorResp.Message)
	}

	var callWalletBasedNamedSmartContractMethodResponse CallWalletBasedNamedSmartContractMethodResponse
	if err = json.NewDecoder(resp.Body).Decode(&callWalletBasedNamedSmartContractMethodResponse); err != nil {
		return CallWalletBasedNamedSmartContractMethodResponse{}, err
	}

	return callWalletBasedNamedSmartContractMethodResponse, nil
}

func makeRoute(req CallWalletBasedNamedSmartContractMethodRequest) (res string) {
	res = "eth/contract/"+req.Contract+"/"+req.Wallet+"/call/"+req.Method
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
