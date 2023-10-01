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

	expectedContentType = "application/json; charset=utf-8"
)

func (d *downloader) GetRates(ctx context.Context) (domain.ExchangeRates, domain.RequestMetadata, error) {
	result, metadata, err := downloadRates()
	if err != nil {
		return domain.ExchangeRates{}, domain.RequestMetadata{}, fmt.Errorf("error downloading rates: %w", err)
	}

	return *result, *metadata, nil
}

func downloadRates() (*domain.ExchangeRates, *domain.RequestMetadata, error) {
	// measure the duration of this request
	start := time.Now()
	// defer func() {

	req, err := http.NewRequest(http.MethodGet, exchangeAPIURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/118.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting rates: %w", err)
	}
	defer resp.Body.Close()
	rDuration := time.Since(start).Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("expected Status OK, got %v", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != expectedContentType {
		return nil, nil, fmt.Errorf("expected Content-Type %v, got %v", expectedContentType, contentType)
	}

	rates := domain.ExchangeRates{}
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding JSON resp: %w", err)
	}

	metadata := domain.RequestMetadata{
		RequestDuration:     rDuration,
		ResponseHTTPCode:    resp.Status,
		ResponseContentType: contentType,
		ResponseValidJSON:   true,
	}

	return &rates, &metadata, nil
}
