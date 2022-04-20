package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeebo/errs"
)

// UpdateWalletRequest fields that describes list update wallet request.
type UpdateWalletRequest struct {
	Wallet string `json:"wallet"`
	Tags   []Tag `json:"tags"`
}

// UpdateWallet update's existing wallet with provided tags.
func (client *Client) UpdateWallet(ctx context.Context, updateWallet UpdateWalletRequest) (response interface{}, err error) {
	jsonBody, err := json.Marshal(updateWallet.Tags)
	if err != nil {
		return CreateWalletResponse{}, err
	}

	var updateWalletResponse struct {
		Wallet   string    `json:"wallet"`
		Version  string    `json:"version"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
		Security string    `json:"security"`
		Public   struct {
			Secp256K1 string `json:"secp256k1"`
		} `json:"public"`
	}

	req, err := http.NewRequest(http.MethodPatch, client.config.TanganyDefaultURL+"wallet/"+updateWallet.Wallet, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json-patch+json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, err
		}

		return nil, errs.New(errorResp.Message)
	}

	if err = json.NewDecoder(resp.Body).Decode(&updateWalletResponse); err != nil {
		return nil, err
	}

	return updateWalletResponse, nil
}
