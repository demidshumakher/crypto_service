package rest

import (
	"context"
	"cryptoserver/domain"
)

type CryptoService interface {
	GetAll(ctx context.Context) ([]domain.Crypto, error)
	Create(ctx context.Context, crypto domain.Crypto) (domain.Crypto, error)
	GetBySymbol(ctx context.Context, symbol string) (domain.Crypto, error)
	UpdateBySymbol(ctx context.Context, crypto domain.Crypto) (domain.Crypto, error)
	GetHistoryBySymbol(ctx context.Context, symbol string) ([]domain.PriceHistory, error)
	GetStatBySymbol(ctx context.Context, symbol string) ([]domain.PriceStats, error)
	Delete(ctx context.Context, symbol string) error
}
