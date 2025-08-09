package service

import (
	"context"
	"cryptoserver/domain"
	"time"
)

type CryptoRepository interface {
	GetAll() ([]domain.Crypto, error)
	GetBySymbol(symbol string) (*domain.Crypto, error)
	Create(symbol, name string, price float64) (*domain.Crypto, error)
	Update(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error)
	Delete(symbol string) error
}

type PriceRepository interface {
	GetHistory(symbol string) ([]domain.PriceHistory, error)
	AddRecord(symbol string, price float64, timestamp time.Time) error
	DeleteHistory(symbol string) error
	GetStatBySymbol(symbol string) (*domain.PriceStats, error)
}

type Service struct {
	cryptoRepo  CryptoRepository
	historyRepo PriceRepository
}

func NewService(c CryptoRepository, h PriceRepository) *Service {
	return &Service{
		cryptoRepo:  c,
		historyRepo: h,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]domain.Crypto, error) {
	return s.cryptoRepo.GetAll()
}

func (s *Service) Create(ctx context.Context, symbol, name string, price float64) (*domain.Crypto, error) {
	return s.cryptoRepo.Create(symbol, name, price)
}

func (s *Service) GetBySymbol(ctx context.Context, symbol string) (*domain.Crypto, error) {
	return s.cryptoRepo.GetBySymbol(symbol)
}

func (s *Service) UpdateBySymbol(ctx context.Context, symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error) {
	return s.cryptoRepo.Update(symbol, name, price, updatedAt)
}

func (s *Service) GetHistoryBySymbol(ctx context.Context, symbol string) ([]domain.PriceHistory, error) {
	return s.historyRepo.GetHistory(symbol)
}

func (s *Service) GetStatBySymbol(ctx context.Context, symbol string) (*domain.PriceStats, error) {
	return s.historyRepo.GetStatBySymbol(symbol)
}

func (s *Service) Delete(ctx context.Context, symbol string) error {
	return s.cryptoRepo.Delete(symbol)
}
