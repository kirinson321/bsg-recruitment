package exchange

import (
	"context"
	"fmt"
	"strconv"
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

func (s *service) GetRates(ctx context.Context) error {
	timestamp := time.Now()

	// get rates and metadata from the downloader
	rates, metadata := s.exchangeDownloader.GetRates(ctx)

	// find the days in which rates are outside of specified range
	targetDays, err := findTargetDays(rates)

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
		return fmt.Errorf("error sending the structured data with timestamp %v to the outputter: %w", timestamp, err)
	}

	return nil
}

var (
	lowerLimit = 4.5
	upperLimit = 4.7
)

func findTargetDays(rates domain.ExchangeRates) ([]string, error) {
	var targetDays []string

	for _, rate := range rates.Rates {

		midValue, err := strconv.ParseFloat(rate.MidValue, 64)
		if err == nil {
			return nil, fmt.Errorf("error converting %v rate's midValue %v to float64: %w", rate.RateNumber, rate.MidValue, err)
		}

		if midValue < lowerLimit || midValue > upperLimit {
			targetDays = append(targetDays, rate.EffectiveDate)
		}
	}

	return targetDays, nil
}
