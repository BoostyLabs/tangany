package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeebo/errs"
)

// GetWalletRequest fields that describes list get wallet request.
type GetWalletRequest struct {
	Wallet string `json:"wallet"`
}

// GetWalletResponse fields that returns from get wallet method.
type GetWalletResponse struct {
	Wallet   string    `json:"wallet"`
	Version  string    `json:"version"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Security string    `json:"security"`
	Public   struct {
		Secp256K1 string `json:"secp256k1"`
	} `json:"public"`
	Tags []interface{} `json:"tags"`
}

// GetWallet returns all wallet's params by wallet name.
func (client *Client) GetWallet(ctx context.Context, getWallet GetWalletRequest) (response GetWalletResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, client.config.TanganyDefaultURL+"wallet/"+getWallet.Wallet, nil)
	if err != nil {
		return GetWalletResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return GetWalletResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return GetWalletResponse{}, err
		}

		return GetWalletResponse{}, errs.New(errorResp.Message)
	}

	var getWalletResponse GetWalletResponse
	if err = json.NewDecoder(resp.Body).Decode(&getWalletResponse); err != nil {
		return GetWalletResponse{}, err
	}

	return getWalletResponse, nil
}
