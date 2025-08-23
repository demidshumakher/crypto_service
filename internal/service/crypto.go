package service

import (
	"context"
	"cryptoserver/domain"
	"cryptoserver/pkg/coingecko"
	"time"
)

type CryptoRepository interface {
	GetAll() ([]domain.Crypto, error)
	GetBySymbol(symbol string) (*domain.Crypto, error)
	Create(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error)
	Update(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error)
	Delete(symbol string) error
	GetHistory(symbol string) ([]domain.PriceHistory, error)
	AddRecord(symbol string, price float64, timestamp time.Time) error
}

type Service struct {
	cryptoRepo CryptoRepository
	gecko      *coingecko.CoinGeckoClient
}

func NewCryptoService(c CryptoRepository, gec *coingecko.CoinGeckoClient) *Service {
	return &Service{
		cryptoRepo: c,
		gecko:      gec,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]domain.Crypto, error) {
	return s.cryptoRepo.GetAll()
}

func (s *Service) Create(ctx context.Context, symbol string) (*domain.Crypto, error) {
	arr, err := s.gecko.GetDataSymbols(symbol)
	if err != nil {
		return nil, err
	}

	crypto := arr[0]
	return s.cryptoRepo.Create(symbol, crypto.Name, crypto.Current_price, crypto.Last_updated)
}

func (s *Service) GetBySymbol(ctx context.Context, symbol string) (*domain.Crypto, error) {
	res, err := s.cryptoRepo.GetBySymbol(symbol)
	return res, err
}

func (s *Service) UpdateBySymbol(ctx context.Context, symbol string) (*domain.Crypto, error) {
	arr, err := s.gecko.GetDataSymbols(symbol)
	if err != nil {
		return nil, err
	}

	crypto := arr[0]

	return s.cryptoRepo.Update(symbol, crypto.Name, crypto.Current_price, crypto.Last_updated)
}

func (s *Service) GetHistoryBySymbol(ctx context.Context, symbol string) ([]domain.PriceHistory, error) {
	return s.cryptoRepo.GetHistory(symbol)
}

func (s *Service) GetStatBySymbol(ctx context.Context, symbol string) (*domain.PriceStats, error) {
	res := &domain.PriceStats{}
	his, err := s.cryptoRepo.GetHistory(symbol)

	if err != nil {
		return nil, err
	}

	res.MaxPrice = his[0].Price
	res.MinPrice = his[0].Price
	res.RecordsCount = len(his)
	sum := 0.0

	for _, value := range his {
		sum += value.Price
		res.MaxPrice = max(res.MaxPrice, value.Price)
		res.MinPrice = min(res.MinPrice, value.Price)
	}

	res.AvgPrice = sum / float64(res.RecordsCount)
	res.PriceChange = res.MaxPrice - res.MinPrice
	res.PriceChangePercent = (100 * res.PriceChange) / res.AvgPrice

	return res, nil
}

func (s *Service) Delete(ctx context.Context, symbol string) error {
	return s.cryptoRepo.Delete(symbol)
}
