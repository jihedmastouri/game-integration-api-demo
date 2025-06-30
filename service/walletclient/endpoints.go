package walletclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

func (w *WalletClient) Deposit(req DepositRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := w.makeJSONRequest("POST", "/api/v1/deposit", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (w *WalletClient) Withdraw(req WithdrawRequest) (*OperationResponse, error) {
	var resp OperationResponse
	if err := w.makeJSONRequest("POST", "/api/v1/withdraw", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (w *WalletClient) GetBalance(userID int) (*BalanceResponse, error) {
	var resp BalanceResponse
	url := fmt.Sprintf("/api/v1/withdraw/%d", userID)
	if err := w.makeJSONRequest("POST", url, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (w *WalletClient) makeJSONRequest(method, endpoint string, payload any, response any) error {
	url := fmt.Sprintf("%s%s", w.baseURL, endpoint)

	var body *bytes.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}

	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.token))

	resp, err := w.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		slog.Error(errResp.Msg, "error", errResp.Code)
		return errors.New(errResp.Code)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
