package mapdb

import (
	"cryptoserver/domain"
	"time"
)

type CryptoRepository struct {
	dp      map[string]*domain.Crypto
	history map[string][]domain.PriceHistory
}

func NewCryptoRepository() *CryptoRepository {
	return &CryptoRepository{
		dp:      make(map[string]*domain.Crypto),
		history: make(map[string][]domain.PriceHistory),
	}
}

func (r *CryptoRepository) GetAll() ([]domain.Crypto, error) {
	result := make([]domain.Crypto, len(r.dp))
	for _, crypto := range r.dp {
		result = append(result, *crypto)
	}
	return result, nil
}

func (r *CryptoRepository) GetBySymbol(symbol string) (*domain.Crypto, error) {
	crypto, ok := r.dp[symbol]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return crypto, nil
}

func (r *CryptoRepository) Create(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error) {
	if _, ok := r.dp[symbol]; ok {
		return nil, domain.ErrAlreadyExist
	}
	newCrypto := &domain.Crypto{
		Symbol:       symbol,
		Name:         name,
		CurrentPrice: price,
		LastUpdated:  updatedAt,
	}
	r.dp[symbol] = newCrypto
	r.history[symbol] = []domain.PriceHistory{
		{
			Symbol:    symbol,
			Price:     price,
			Timestamp: newCrypto.LastUpdated,
		},
	}
	return newCrypto, nil
}

func (r *CryptoRepository) Update(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error) {
	crypto, ok := r.dp[symbol]
	if !ok {
		return nil, domain.ErrNotFound
	}
	crypto.CurrentPrice = price
	crypto.LastUpdated = updatedAt
	r.history[symbol] = append(r.history[symbol], domain.PriceHistory{
		Symbol:    symbol,
		Price:     price,
		Timestamp: updatedAt,
	})
	return crypto, nil
}

func (r *CryptoRepository) Delete(symbol string) error {
	if _, ok := r.dp[symbol]; !ok {
		return domain.ErrNotFound
	}
	delete(r.dp, symbol)
	delete(r.history, symbol)
	return nil
}

func (r *CryptoRepository) GetHistory(symbol string) ([]domain.PriceHistory, error) {
	his, ok := r.history[symbol]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return his, nil
}

func (r *CryptoRepository) AddRecord(symbol string, price float64, timestamp time.Time) error {
	if _, ok := r.history[symbol]; !ok {
		return domain.ErrNotFound
	}

	r.history[symbol] = append(r.history[symbol], domain.PriceHistory{
		Symbol:    symbol,
		Price:     price,
		Timestamp: timestamp,
	})
	if len(r.history[symbol]) > 100 {
		r.history[symbol] = r.history[symbol][len(r.history[symbol])-100:]
	}
	return nil
}
