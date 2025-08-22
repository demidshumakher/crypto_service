package coingecko

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type CoinGeckoClient struct {
	apiKey  string
	baseURL string
	db      map[string]string
}

func NewCoinGeckoClient(apiKey, baseURL string) *CoinGeckoClient {
	return &CoinGeckoClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		db:      make(map[string]string),
	}
}

type coinResponse struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func (cg *CoinGeckoClient) GetIdBySymbol(symbol string) (string, error) {
	if id, ok := cg.db[symbol]; ok {
		return id, nil
	}

	req, err := http.NewRequest("GET", cg.baseURL+"/coins/list", nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("x-cg-demo-api-key", cg.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res := &[]coinResponse{}

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return "", err
	}

	var id string

	for _, el := range *res {
		if strings.EqualFold(el.Symbol, symbol) {
			id = el.Id
			break
		}
	}

	cg.db[symbol] = id

	return id, nil
}

type InfoResponse struct {
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Current_price float64   `json:"current_price"`
	Last_updated  time.Time `json:"last_updated"`
}

func (cg *CoinGeckoClient) GetDataSymbols(symbols ...string) ([]InfoResponse, error) {
	ids := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		if id, ok := cg.db[symbol]; ok {
			ids = append(ids, id)
		} else {
			id, err := cg.GetIdBySymbol(symbol)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}

	url := cg.baseURL + "/coins/markets?vs_currency=rub&ids=" + strings.Join(ids, "%2C")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("x-cg-demo-api-key", cg.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res := &[]InfoResponse{}

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	return *res, nil
}
