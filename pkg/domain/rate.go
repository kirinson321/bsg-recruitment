package domain

type (
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
