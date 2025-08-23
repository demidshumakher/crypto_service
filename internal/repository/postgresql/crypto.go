package postgresql

import (
	"cryptoserver/domain"
	"database/sql"
	"errors"
	"time"
)

type CryptoRepository struct {
	db *sql.DB
}

func NewCryptoRepository(db *sql.DB) *CryptoRepository {
	return &CryptoRepository{
		db: db,
	}
}

func (r *CryptoRepository) GetAll() ([]domain.Crypto, error) {
	rows, err := r.db.Query(`select c.symbol, c.name, p.price as current_price, p.updated_at
        from cryptos c
        inner join (
            select distinct on (crypto) crypto, price, updated_at from prices
            order by crypto, updated_at desc
            ) p on p.crypto = c.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []domain.Crypto{}
	for rows.Next() {
		crypto := domain.Crypto{}
		err := rows.Scan(&crypto.Symbol, &crypto.Name, &crypto.CurrentPrice, &crypto.LastUpdated)
		if err != nil {
			return nil, err
		}
		res = append(res, crypto)
	}

	// Проверяем ошибки после цикла
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *CryptoRepository) GetBySymbol(symbol string) (*domain.Crypto, error) {
	res := &domain.Crypto{}

	err := r.db.QueryRow(`select c.symbol, c.name, p.price as current_price, p.updated_at
        from cryptos c
        left join (
            select distinct on (crypto) crypto, price, updated_at from prices
            order by crypto, updated_at desc
            ) p on p.crypto = c.id
        where c.symbol = $1`, symbol).Scan(&res.Symbol, &res.Name, &res.CurrentPrice, &res.LastUpdated)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *CryptoRepository) Create(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error) {
	_, err := r.db.Exec("insert into cryptos (symbol, name) values ($1, $2)", symbol, name)
	if err != nil {
		return nil, err
	}

	var id int

	err = r.db.QueryRow("select id from cryptos where symbol = $1", symbol).Scan(&id)
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec("insert into prices (crypto, price, updated_at) values ($1, $2, $3)", id, price, updatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Crypto{
		Symbol:       symbol,
		Name:         name,
		CurrentPrice: price,
		LastUpdated:  updatedAt,
	}, nil
}

func (r *CryptoRepository) Update(symbol, name string, price float64, updatedAt time.Time) (*domain.Crypto, error) {
	_, err := r.db.Exec("update cryptos set name = $1 where symbol = $2", name, symbol)
	if err != nil {
		return nil, err
	}

	var id int
	err = r.db.QueryRow("select id from cryptos where symbol = $1", symbol).Scan(&id)
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec("insert into prices (crypto, price, updated_at) values ($1, $2, $3)", id, price, updatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Crypto{
		Symbol:       symbol,
		Name:         name,
		CurrentPrice: price,
		LastUpdated:  updatedAt,
	}, nil
}

func (r *CryptoRepository) Delete(symbol string) error {
	_, err := r.db.Exec("delete from cryptos where symbol = $1", symbol)
	return err
}

func (r *CryptoRepository) GetHistory(symbol string) ([]domain.PriceHistory, error) {
	rows, err := r.db.Query(`
		select c.symbol, p.price, p.updated_at 
		from prices p
		inner join cryptos c on p.crypto = c.id
		where c.symbol = $1
		order by p.updated_at desc
	`, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []domain.PriceHistory{}
	for rows.Next() {
		record := domain.PriceHistory{}
		err := rows.Scan(&record.Symbol, &record.Price, &record.Timestamp)
		if err != nil {
			return nil, err
		}
		history = append(history, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (r *CryptoRepository) AddRecord(symbol string, price float64, timestamp time.Time) error {
	var id int
	err := r.db.QueryRow("select id from cryptos where symbol = $1", symbol).Scan(&id)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("insert into prices (crypto, price, updated_at) values ($1, $2, $3)", id, price, timestamp)
	return err
}
