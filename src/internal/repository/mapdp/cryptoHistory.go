package mapdp

import (
	"cryptoserver/domain"
	"time"
)

func (r *CryptoRepository) AddRecord(symbol string, price float64, timestamp time.Time) error {
	history, ok := r.history[symbol]
	if !ok {
		return domain.ErrNotFound
	}
	history = append(history, domain.PriceHistory{
		Symbol:    symbol,
		Price:     price,
		Timestamp: timestamp,
	})
	return nil
}
