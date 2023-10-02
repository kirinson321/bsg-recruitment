// Package domain contains interfaces and structs that reflect
// the data structures and business logic used in the application.
package domain

import (
	"context"
	"time"
)

type (
	// ExchangeService is the interface that provides an entrypoint
	// to the application's currency rate retrieval, processing and output logic.
	ExchangeService interface {
		GetRates(ctx context.Context) error
	}

	// ExchangeDownloader is the interface that provides functionality
	// for downloading data from the currency exchange API.
	ExchangeDownloader interface {
		GetRates(ctx context.Context) (ExchangeRates, RequestMetadata, error)
	}

	// Outputter is the interface that defines how the already processed data should be outputted.
	Outputter interface {
		Output(data StructuredOutput) error
	}

	// StructuredOutput defines the structure of the data that is outputted by the application.
	StructuredOutput struct {
		OutputTimestamp     time.Time `json:"outputTimestamp"`
		RequestDuration     int64     `json:"requestDuration"`
		ResponseHTTPCode    string    `json:"respHTTPCode"`
		ResponseContentType string    `json:"respContentType"`
		ResponseValidJSON   bool      `json:"respValidJSON"`
		TargetDays          []string  `json:"targetDays,omitempty"`
	}

	// ExchangeRates defines the structure of the data that is downloaded from the currency exchange API.
	ExchangeRates struct {
		Table        string       `json:"table"`
		Currency     string       `json:"currency"`
		CurrencyCode string       `json:"code"`
		Rates        []SingleRate `json:"rates"`
	}

	// SingleRate defines the structure of a single currency rate.
	SingleRate struct {
		RateNumber    string  `json:"no"`
		EffectiveDate string  `json:"effectiveDate"`
		MidValue      float64 `json:"mid"`
	}

	// RequestMetadata defines the structure of the metadata that is returned by the downloader interface.
	RequestMetadata struct {
		RequestDuration     int64
		ResponseHTTPCode    string
		ResponseContentType string
		ResponseValidJSON   bool
	}
)
