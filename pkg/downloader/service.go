package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kirinson321/bsg-recruitment/pkg/domain"
)

type downloader struct {
}

func NewDownloader() domain.ExchangeDownloader {
	return &downloader{}
}

var (
	exchangeAPIURL = "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/?format=json"
)

func (d *downloader) GetRates(ctx context.Context) (*domain.ExchangeRates, *domain.RequestMetadata, error) {
	result := domain.ExchangeRates{}
	metadata, err := downloadRates(result)
	if err != nil {
		return nil, nil, fmt.Errorf("error downloading rates: %w", err)
	}

	return &result, metadata, nil
}

func downloadRates(rates domain.ExchangeRates) (*domain.RequestMetadata, error) {
	// measure the duration of this request
	start := time.Now()
	// defer func() {
	resp, err := http.Get(exchangeAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error getting rates: %w", err)
	}
	defer resp.Body.Close()
	rDuration := time.Since(start).Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected Status OK, got %v", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type application/json, got %v", contentType)
	}

	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON resp: %w", err)
	}

	metadata := domain.RequestMetadata{
		RequestDuration:     rDuration,
		ResponseHTTPCode:    resp.Status,
		ResponseContentType: contentType,
		ResponseValidJSON:   true,
	}

	return &metadata, nil
}
