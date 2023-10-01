package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kirinson321/bsg-recruitment/pkg/domain"
)

type downloader struct {
	httpClient *http.Client
}

// NewDownloader returns a new instance of the ExchangeDownloader.
func NewDownloader(c *http.Client) domain.ExchangeDownloader {
	return &downloader{
		httpClient: c,
	}
}

var (
	exchangeAPIURL = "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/?format=json"

	expectedContentType = "application/json; charset=utf-8"
)

func (d *downloader) GetRates(ctx context.Context) (domain.ExchangeRates, domain.RequestMetadata, error) {
	result, metadata, err := d.downloadRates()
	if err != nil {
		return domain.ExchangeRates{}, domain.RequestMetadata{}, fmt.Errorf("error downloading rates: %w", err)
	}

	return *result, *metadata, nil
}

// downloadRates is the function that handles the logic of downloading the rates from the exchange API.
func (d *downloader) downloadRates() (*domain.ExchangeRates, *domain.RequestMetadata, error) {
	// Construct the request.
	req, err := http.NewRequest(http.MethodGet, exchangeAPIURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %w", err)
	}
	// Set the User-Agent header to avoid getting blocked by the API.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/118.0")

	// Start measuring the duration of this request.
	start := time.Now()

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting rates: %w", err)
	}
	defer resp.Body.Close()
	// Finish measuring the duration of the request.
	rDuration := time.Since(start).Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("expected Status OK, got %v", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, expectedContentType) {
		return nil, nil, fmt.Errorf("expected Content-Type %v, got %v", expectedContentType, contentType)
	}

	// Assemble the output of the function.
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
