package rest

import (
	"context"
	"cryptoserver/domain"
	"time"
)

type CryptoService interface {
	GetAll(ctx context.Context) ([]domain.Crypto, error)
	Create(ctx context.Context, symbol, name string, price float64) (*domain.Crypto, error)
	GetBySymbol(ctx context.Context, symbol string) (*domain.Crypto, error)
	UpdateBySymbol(ctx context.Context, symbol, name string, price float64, timestamp time.Time) (*domain.Crypto, error)
	GetHistoryBySymbol(ctx context.Context, symbol string) ([]domain.PriceHistory, error)
	GetStatBySymbol(ctx context.Context, symbol string) (*domain.PriceStats, error)
	Delete(ctx context.Context, symbol string) error
}
