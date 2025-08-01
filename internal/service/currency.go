package service

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/sp-va/cryptoCurrencies/internal/dto"
	"github.com/sp-va/cryptoCurrencies/internal/repository"
)

type CurrencyService struct {
	repo repository.CurrencyRepository
}

func NewCurrencyService(r repository.CurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: r}
}

func (s *CurrencyService) InsertCurrency(ctx context.Context, c dto.AddCurrency) error {
	return s.repo.InsertCurrencyToTrack(ctx, c)
}

func (s *CurrencyService) DeleteCurrencyFromTracking(ctx context.Context, coin string) (int64, error) {
	return s.repo.DeleteCurrencyFromTracking(ctx, coin)
}

func (s *CurrencyService) GetCoinValue(ctx context.Context, coin string, timestamp uint32) (decimal.Decimal, error) {
	coinObj, err := s.repo.GetCoinValue(ctx, coin, timestamp)

	if err != nil {
		return coinObj.Price, err
	}

	return coinObj.Price, nil

}
