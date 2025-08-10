package rest

import (
	"context"
	"cryptoserver/domain"
	"encoding/json"
	"net/http"
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

func RegisterCryptoHandler(cs CryptoService, mx *http.ServeMux) {
	ch := &CryptoHandler{
		cryptoServ: cs,
	}
	mx.HandleFunc("GET /crypto", ch.GetAllHandler)
	mx.HandleFunc("POST /crypto", ch.AddHandler)
	mx.HandleFunc("GET /crypto/{symbol}", ch.GetBySymbolHandler)
	mx.HandleFunc("PUT /crypto/{symbol}/refresh", ch.UpdateHandler)
	mx.HandleFunc("GET /crypto/{symbol}/history", ch.GetHistoryHandler)
	mx.HandleFunc("GET /crypto/{symbol}/stats", ch.GetStatHandler)
	mx.HandleFunc("DELETE /crypto/{symbol}", ch.DeleteHandler)
}

func (c *CryptoHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	value, err := c.cryptoServ.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.Create(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) GetBySymbolHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetBySymbol(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.UpdateBySymbol(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetHistoryBySymbol(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) GetStatHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetHistoryBySymbol(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}

func (c *CryptoHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	value, err := c.cryptoServ.GetHistoryBySymbol(r.Context(), symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(value)
}
