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

	ContentTypePrefix = "Content-Type: "
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
	respValidJSON := true

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

	// Assemble the Content-Type.
	mediaType, _ := parseMediaType(resp.Header.Get("Content-Type"))
	contentType := ContentTypePrefix + mediaType

	// Assemble the output of the function.
	rates := domain.ExchangeRates{}
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		respValidJSON = false
	}

	metadata := domain.RequestMetadata{
		RequestDuration:     rDuration,
		ResponseHTTPCode:    resp.Status,
		ResponseContentType: contentType,
		ResponseValidJSON:   respValidJSON,
	}

	return &rates, &metadata, nil
}

// parseMediaType extracts the media type and it's parameters from the Content-Type header.
func parseMediaType(contentType string) (string, map[string]string) {
	mediaTypeAndParams := strings.Split(contentType, ";")
	mediaType := ""
	params := make(map[string]string)

	// parse the input to extract the media type and it's parameters
	for _, val := range mediaTypeAndParams {
		// if the value contains a slash, it's the media type
		if strings.Contains(val, "/") {
			mediaType = strings.TrimSpace(val)
			continue
		}

		// if the value contains an equal sign, it's a parameter
		if strings.Contains(val, "=") {
			parts := strings.SplitN(val, "=", 2)
			if len(parts) == 2 {
				// trim the key and value of whitespaces and add them to the params map
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				params[key] = value
			}
		}
	}

	// assemble the output
	return mediaType, params
}
