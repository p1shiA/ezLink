package zarinpal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/p1shiA/ezLink/internal/core/domain"
)

type Config struct {
	MerchantID string
	Sandbox    bool
}

type client struct {
	merchantID string
	sandbox    bool
	baseURL    string
	startURL   string
	httpClient *http.Client
}

func NewClient(cfg Config) domain.PaymentGateway {
	baseURL := "https://payment.zarinpal.com/pg"
	startURL := "https://payment.zarinpal.com/pg/StartPay"

	if cfg.Sandbox {
		baseURL = "https://sandbox.zarinpal.com/pg"
		startURL = "https://sandbox.zarinpal.com/pg/StartPay"
	}

	return &client{
		merchantID: cfg.MerchantID,
		sandbox:    cfg.Sandbox,
		baseURL:    baseURL,
		startURL:   startURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *client) RequestPayment(ctx context.Context, amountToman int64, description, callbackURL string) (string, string, error) {
	reqBody := PaymentRequest{
		MerchantID:  c.merchantID,
		Amount:      amountToman,
		Currency:    "IRT",
		Description: description,
		CallbackURL: callbackURL,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	url := fmt.Sprintf("%s/v4/payment/request.json", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var res PaymentResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", "", fmt.Errorf("invalid response from zarinpal: %w", err)
	}

	if res.Data == nil || res.Data.Authority == "" {
		return "", "", fmt.Errorf("zarinpal request failed: %s (errors: %v)", string(body), res.Errors)
	}

	authority := res.Data.Authority
	paymentURL := fmt.Sprintf("%s/%s", c.startURL, authority)

	return authority, paymentURL, nil
}

func (c *client) VerifyPayment(ctx context.Context, authority string, amountToman int64) (string, error) {
	reqBody := VerifyRequest{
		MerchantID: c.merchantID,
		Amount:     amountToman,
		Authority:  authority,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v4/payment/verify.json", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res VerifyResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("invalid response from zarinpal verify: %w", err)
	}

	if res.Data == nil || (res.Data.Code != 100 && res.Data.Code != 101) {
		return "", fmt.Errorf("payment verification rejected: %s (errors: %v)", string(body), res.Errors)
	}

	refID := fmt.Sprintf("%v", res.Data.RefID)
	return refID, nil
}
