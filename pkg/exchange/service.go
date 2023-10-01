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

func NewService(exchangeDownloader domain.ExchangeDownloader, outputter domain.Outputter) domain.ExchangeService {
	return &service{
		exchangeDownloader: exchangeDownloader,
		outputter:          outputter,
	}
}

const (
	interval       = 5
	numberOfChecks = 10
)

func (s *service) GetRates(ctx context.Context) error {
	tick := time.Tick(interval * time.Second)
	for range tick {
		for i := 0; i < numberOfChecks; i++ {
			go s.getRates(context.Background())
		}
	}

	return nil
}

func (s *service) getRates(ctx context.Context) {
	timestamp := time.Now()

	// get rates and metadata from the downloader
	rates, metadata, err := s.exchangeDownloader.GetRates(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("error getting rates from the downloader for the %v timestamp: %w", timestamp, err))
		return
	}

	// find the days in which rates are outside of specified range
	targetDays := findTargetDays(rates)

	// pack it into the StructuredOutput
	o := domain.StructuredOutput{
		OutputTimestamp:     timestamp,
		RequestDuration:     metadata.RequestDuration,
		ResponseHTTPCode:    metadata.ResponseHTTPCode,
		ResponseContentType: metadata.ResponseContentType,
		ResponseValidJSON:   metadata.ResponseValidJSON,
		TargetDays:          targetDays,
	}

	// send the data to the outputter
	err = s.outputter.Output(o)
	if err != nil {
		fmt.Println(fmt.Errorf("error sending the structured data with timestamp %v to the outputter: %w", timestamp, err))
		return
	}

	return
}

var (
	lowerLimit = 4.5
	upperLimit = 4.7
)

func findTargetDays(rates domain.ExchangeRates) []string {
	var targetDays []string

	for _, rate := range rates.Rates {
		if rate.MidValue < lowerLimit || rate.MidValue > upperLimit {
			targetDays = append(targetDays, rate.EffectiveDate)
		}
	}

	return targetDays
}
