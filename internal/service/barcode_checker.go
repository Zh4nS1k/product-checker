package service

import (
	"barcode-checker/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BarcodeChecker interface {
	Check(barcode string) (bool, error)
}

type ExternalBarcodeChecker struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

type localBarcodeChecker struct{}

func NewLocalBarcodeChecker() BarcodeChecker {
	return &localBarcodeChecker{}
}

func NewExternalBarcodeChecker(apiKey string) *ExternalBarcodeChecker {
	return &ExternalBarcodeChecker{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: "https://api.barcodelookup.com/v3/products",
	}
}

func (e *ExternalBarcodeChecker) Check(barcode string) (bool, error) {
	if e.apiKey == "" {
		return len(barcode) > 8, nil
	}

	url := fmt.Sprintf("%s?barcode=%s&key=%s", e.baseURL, barcode, e.apiKey)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Products []struct {
			Barcode string `json:"barcode_number"`
		} `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return len(result.Products) > 0 && result.Products[0].Barcode == barcode, nil
}

func (l *localBarcodeChecker) Check(barcode string) (bool, error) {
	return utils.IsBarcodeValid(barcode), nil
}
