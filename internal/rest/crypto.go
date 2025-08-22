package rest

import (
	"context"
	"cryptoserver/domain"
	"encoding/json"
	"net/http"
	"time"
)

type CryptoService interface {
	GetAll(ctx context.Context) ([]domain.Crypto, error)
	Create(ctx context.Context, symbol string) (*domain.Crypto, error)
	GetBySymbol(ctx context.Context, symbol string) (*domain.Crypto, error)
	UpdateBySymbol(ctx context.Context, symbol string) (*domain.Crypto, error)
	GetHistoryBySymbol(ctx context.Context, symbol string) ([]domain.PriceHistory, error)
	GetStatBySymbol(ctx context.Context, symbol string) (*domain.PriceStats, error)
	Delete(ctx context.Context, symbol string) error
}

type CryptoHandler struct {
	cryptoServ CryptoService
}

func NewCryptoHandler(cs CryptoService, mx *Router) {
	ch := &CryptoHandler{
		cryptoServ: cs,
	}
	mx.Handle("GET /crypto", http.HandlerFunc(ch.GetAllHandler))
	mx.Handle("POST /crypto", http.HandlerFunc(ch.AddHandler))
	mx.Handle("GET /crypto/{symbol}", http.HandlerFunc(ch.GetBySymbolHandler))
	mx.Handle("PUT /crypto/{symbol}/refresh", http.HandlerFunc(ch.UpdateHandler))
	mx.Handle("GET /crypto/{symbol}/history", http.HandlerFunc(ch.GetHistoryHandler))
	mx.Handle("GET /crypto/{symbol}/stats", http.HandlerFunc(ch.GetStatHandler))
	mx.Handle("DELETE /crypto/{symbol}", http.HandlerFunc(ch.DeleteHandler))
}

type cryptoResponse struct {
	Crypto *domain.Crypto `json:"crypto"`
}

func (c *CryptoHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	value, err := c.cryptoServ.GetAll(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}

	res := &struct {
		Cryptos []domain.Crypto `json:"cryptos"`
	}{
		Cryptos: value,
	}

	json.NewEncoder(w).Encode(res)
}

func (c *CryptoHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := &struct {
		Symbol string `json:"symbol"`
	}{}

	json.NewDecoder(r.Body).Decode(reqBody)
	defer r.Body.Close()

	if reqBody.Symbol == "" {
		WriteError(w, domain.ErrBadRequest)
		return
	}

	value, err := c.cryptoServ.Create(r.Context(), reqBody.Symbol)
	if err != nil {
		WriteError(w, err)
		return
	}

	res := cryptoResponse{
		Crypto: value,
	}

	json.NewEncoder(w).Encode(res)
}

func (c *CryptoHandler) GetBySymbolHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetBySymbol(r.Context(), symbol)
	if err != nil {
		WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	value, err := c.cryptoServ.UpdateBySymbol(r.Context(), symbol)
	if err != nil {
		WriteError(w, err)
		return
	}

	res := cryptoResponse{
		Crypto: value,
	}

	json.NewEncoder(w).Encode(res)
}

type historyResponse struct {
	Symbol  string           `json:"symbol"`
	History []historyElement `json:"history"`
}

type historyElement struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *CryptoHandler) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	value, err := c.cryptoServ.GetHistoryBySymbol(r.Context(), symbol)
	if err != nil {
		WriteError(w, err)
		return
	}

	res := historyResponse{}
	res.Symbol = value[0].Symbol

	for _, el := range value {
		res.History = append(res.History, historyElement{el.Price, el.Timestamp})
	}

	json.NewEncoder(w).Encode(res)
}

type statResponse struct {
	Symbol        string             `json:"symbol"`
	Current_price float64            `json:"current_price"`
	Stats         *domain.PriceStats `json:"stats"`
}

func (c *CryptoHandler) GetStatHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetStatBySymbol(r.Context(), symbol)
	if err != nil {
		WriteError(w, err)
		return
	}

	res := statResponse{}
	res.Symbol = symbol
	res.Stats = value

	t, err := c.cryptoServ.GetBySymbol(r.Context(), symbol)
	if err != nil {
		WriteError(w, err)
		return
	}

	res.Current_price = t.CurrentPrice

	json.NewEncoder(w).Encode(res)
}

func (c *CryptoHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetHistoryBySymbol(r.Context(), symbol)
	if err != nil {
		WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(value)
}
