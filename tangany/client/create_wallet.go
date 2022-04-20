package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeebo/errs"
)

// CreateWalletRequest fields that required for create wallet request.
type CreateWalletRequest struct {
	Wallet string `json:"wallet"`
	UseHsm bool `json:"useHsm"`
	Tags []Tag`json:"tags"`
}

// CreateWalletResponse fields that returns from create wallet.
type CreateWalletResponse struct {
	Wallet   string    `json:"wallet"`
	Version  string    `json:"version"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Security string    `json:"security"`
	Public   struct {
		Secp256K1 string `json:"secp256k1"`
	} `json:"public"`
	Tags []struct {
		HasKYC bool `json:"hasKYC"`
	} `json:"tags"`
}

// CreateWallet creates wallet
func (client *Client) CreateWallet(ctx context.Context, wallet CreateWalletRequest) (response CreateWalletResponse, err error) {
	jsonBody, err := json.Marshal(wallet)
	if err != nil {
		return CreateWalletResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"wallets", bytes.NewReader(jsonBody))
	if err != nil {
		return CreateWalletResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("tangany-client-id", client.config.TanganyClientID)
	req.Header.Set("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Set("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Set("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return CreateWalletResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	// sad, _ := ioutil.ReadAll(resp.Body)
	// print(sad)

	if resp.StatusCode != http.StatusCreated {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return CreateWalletResponse{}, err
		}

		return CreateWalletResponse{}, errs.New(errorResp.Message)
	}

	var createWalletResponse CreateWalletResponse
	if err = json.NewDecoder(resp.Body).Decode(&createWalletResponse); err != nil {
		return CreateWalletResponse{}, err
	}

	return createWalletResponse, nil
}
