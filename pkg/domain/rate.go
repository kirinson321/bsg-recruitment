package domain

import (
	"context"
	"time"
)

type (
	ExchangeService interface {
		GetRates(ctx context.Context) error
	}

	ExchangeDownloader interface {
		GetRates(ctx context.Context) (ExchangeRates, RequestMetadata, error)
	}

	Outputter interface {
		Output(data StructuredOutput) error
	}

	StructuredOutput struct {
		OutputTimestamp     time.Time `json:"outputTimestamp"`
		RequestDuration     int64     `json:"requestDuration"`
		ResponseHTTPCode    string    `json:"respHTTPCode"`
		ResponseContentType string    `json:"respContentType"`
		ResponseValidJSON   bool      `json:"respValidJSON"`
		TargetDays          []string  `json:"targetDays"`
	}

	ExchangeRates struct {
		Table        string       `json:"table"`
		Currency     string       `json:"currency"`
		CurrencyCode string       `json:"code"`
		Rates        []SingleRate `json:"rates"`
	}

	SingleRate struct {
		RateNumber    string  `json:"no"`
		EffectiveDate string  `json:"effectiveDate"`
		MidValue      float64 `json:"mid"`
	}

	RequestMetadata struct {
		RequestDuration     int64
		ResponseHTTPCode    string
		ResponseContentType string
		ResponseValidJSON   bool
	}
)
