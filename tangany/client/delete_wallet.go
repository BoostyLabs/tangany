package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeebo/errs"
)

// DeleteWalletRequest fields that describes delete wallet request.
type DeleteWalletRequest struct {
	Wallet string `json:"wallet"`
}

// DeleteWalletResponse fields that returns from delete wallet method.
type DeleteWalletResponse struct {
	RecoveryID         string    `json:"recoveryId"`
	ScheduledPurgeDate time.Time `json:"scheduledPurgeDate"`
}

// DeleteWallet deletes wallet by wallet name.
func (client *Client) DeleteWallet(ctx context.Context, deleteWallet DeleteWalletRequest) (response DeleteWalletResponse, err error) {
	req, err := http.NewRequest(http.MethodDelete, client.config.TanganyDefaultURL+"wallet/"+deleteWallet.Wallet, nil)
	if err != nil {
		return DeleteWalletResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return DeleteWalletResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return DeleteWalletResponse{}, err
		}

		return DeleteWalletResponse{}, errs.New(errorResp.Message)
	}

	var deleteWalletResponse DeleteWalletResponse
	if err = json.NewDecoder(resp.Body).Decode(&deleteWalletResponse); err != nil {
		return DeleteWalletResponse{}, err
	}

	return deleteWalletResponse, nil
}
