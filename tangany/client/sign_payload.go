package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zeebo/errs"
	"net/http"
)

// SignPayloadRequest fields that describes sign payload request sends to Tangany API.
type SignPayloadRequest struct {
	Payload string `json:"payload"`
	Wallet  string `json:"wallet"`
}

// SignPayloadResponse fields that returns from sign payload method.
type SignPayloadResponse struct {
	Signature string `json:"signature"`
	Encoding  string `json:"encoding"`
}

// SignPayload signs payload, returns signature.
func (client *Client) SignPayload(ctx context.Context, sign SignPayloadRequest) (response SignPayloadResponse, err error) {
	var body struct {
		Payload string `json:"payload"`
	}

	body.Payload = sign.Payload

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return SignPayloadResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, client.config.TanganyDefaultURL+"wallet/"+sign.Wallet+"/sign", bytes.NewReader(jsonBody))
	if err != nil {
		return SignPayloadResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("tangany-client-id", client.config.TanganyClientID)
	req.Header.Add("tangany-client-secret", client.config.TanganyClientSecret)
	req.Header.Add("tangany-vault-url", client.config.TanganyVaultURL)
	req.Header.Add("tangany-subscription", client.config.TanganySubscription)

	resp, err := client.http.Do(req.WithContext(ctx))
	if err != nil {
		return SignPayloadResponse{}, err
	}

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		errorResp := ErrorResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return SignPayloadResponse{}, err
		}

		return SignPayloadResponse{}, errs.New(errorResp.Message)
	}

	var signPayloadResponse SignPayloadResponse
	if err = json.NewDecoder(resp.Body).Decode(&signPayloadResponse); err != nil {
		return SignPayloadResponse{}, err
	}

	return signPayloadResponse, nil
}
