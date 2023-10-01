package exchange

import (
	"context"
	"fmt"
	"time"

	"github.com/kirinson321/bsg-recruitment/pkg/domain"
)

type service struct {
	exchangeDownloader domain.ExchangeDownloader
	outputter          domain.Outputter
}

// NewService returns a new instance of the ExchangeService.
func NewService(exchangeDownloader domain.ExchangeDownloader, outputter domain.Outputter) domain.ExchangeService {
	return &service{
		exchangeDownloader: exchangeDownloader,
		outputter:          outputter,
	}
}

const (
	interval       = 5 * time.Second
	numberOfChecks = 2
)

// GetRates is a wrapper function for the handleRates, which also schedules it's concurrent executions.
func (s *service) GetRates(ctx context.Context) error {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		for i := 0; i < numberOfChecks; i++ {
			go s.handleRates(context.Background())
		}
	}

	return nil
}

// handleRates is the function that handles the logic of getting the rates from the downloader,
// processing them and sending them to the outputter.
func (s *service) handleRates(ctx context.Context) {
	// Initiate the timestamp to record the time of the request.
	timestamp := time.Now()

	// Get rates and metadata from the downloader.
	rates, metadata, err := s.exchangeDownloader.GetRates(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("error getting rates from the downloader for the %v timestamp: %w", timestamp, err))
		return
	}

	// Find the days in which rates are outside of specified range
	targetDays := findTargetDays(rates)

	// Pack the data into the StructuredOutput.
	o := domain.StructuredOutput{
		OutputTimestamp:     timestamp,
		RequestDuration:     metadata.RequestDuration,
		ResponseHTTPCode:    metadata.ResponseHTTPCode,
		ResponseContentType: metadata.ResponseContentType,
		ResponseValidJSON:   metadata.ResponseValidJSON,
		TargetDays:          targetDays,
	}

	// Send the data to the outputter.
	err = s.outputter.Output(o)
	if err != nil {
		fmt.Println(fmt.Errorf(
			"error sending the structured data with timestamp %v to the outputter: %w",
			timestamp,
			err,
		))
		return
	}
}

var (
	lowerLimit = 4.5
	upperLimit = 4.7
)

// findTargetDays finds the days in which the exchange rates are outside of the specified range.
func findTargetDays(rates domain.ExchangeRates) []string {
	var targetDays []string

	for _, rate := range rates.Rates {
		if rate.MidValue < lowerLimit || rate.MidValue > upperLimit {
			targetDays = append(targetDays, rate.EffectiveDate)
		}
	}

	return targetDays
}
