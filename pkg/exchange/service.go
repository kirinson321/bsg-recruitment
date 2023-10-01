package exchange

import (
	"context"

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
	return nil
}
