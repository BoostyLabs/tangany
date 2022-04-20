package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeebo/errs"
)

// VerifySignatureRequest fields that describes verify signature request sends to Tangany API.
type VerifySignatureRequest struct {
	Wallet    string `json:"wallet"`
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

// VerifySignatureResponse fields that returns from verify signature method.
type VerifySignatureResponse struct {
	IsValid bool `json:"isValid"`
}

// VerifySignature verifies signature.
func (client *Client) VerifySignature(ctx context.Context, verify VerifySignatureRequest) (response VerifySignatureResponse, err error) {
	jsonBody, err := json.Marshal(verify)
	if err != nil {
		return VerifySignatureResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"wallet/"+verify.Wallet+"/verify", bytes.NewReader(jsonBody))
	if err != nil {
		return VerifySignatureResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return VerifySignatureResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return VerifySignatureResponse{}, err
		}

		return VerifySignatureResponse{}, errs.New(errorResp.Message)
	}

	var verifySignatureResponse VerifySignatureResponse
	if err = json.NewDecoder(resp.Body).Decode(&verifySignatureResponse); err != nil {
		return VerifySignatureResponse{}, err
	}

	return verifySignatureResponse, nil
}
