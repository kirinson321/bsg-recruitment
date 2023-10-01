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
		GetRates(ctx context.Context)
	}

	Outputter interface {
		Output(o StructuredOutput) error
	}

	StructuredOutput struct {
		OutputTimestamp     time.Time     `json:"outputTimestamp"`
		RequestDuration     int64         `json:"rquestDuration"`
		ResponseHTTPCode    string        `json:"respHTTPCode"`
		ResponseContentType string        `json:"respContentType"`
		ResponseValidJSON   bool          `json:"respValidJSON"`
		RatesData           ExchangeRates `json:"ratesData"`
	}

	ExchangeRates struct {
		Table        string       `json:"table"`
		Currency     string       `json:"currency"`
		CurrencyCode string       `json:"code"`
		Rates        []SingleRate `json:"rates"`
	}

	SingleRate struct {
		RateNumber    string `json:"no"`
		EffectiveDate string `json:"effectiveDate"`
		MidValue      string `json:"mid"`
	}
)
