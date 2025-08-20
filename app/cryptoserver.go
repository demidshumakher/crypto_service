package main

import (
	"cryptoserver/internal/repository/mapdb"
	"cryptoserver/internal/rest"
	"cryptoserver/internal/service"
	"cryptoserver/pkg/coingecko"
	"cryptoserver/pkg/jwt"
	"net/http"
	"time"
)

var (
	jwtSecret = []byte("my_secret_key")
	apiKey    = "CG-VThWUzH2txUEJ1eSaYFga8QK"
	baseURL   = "https://api.coingecko.com/api/v3"
)

func main() {
	mx := http.NewServeMux()
	jwtConfig := jwt.JWTConfig{
		Secret:         jwtSecret,
		ExpirationTime: time.Now().Add(time.Hour * 24).Unix(),
	}

	rest.NewCryptoHandler(
		service.NewCryptoService(
			mapdb.NewCryptoRepository(),
			coingecko.NewCoinGeckoClient(apiKey, baseURL),
		),
		mx,
	)

	rest.NewUserHandler(
		service.NewUserService(jwtConfig, mapdb.NewUserRepository()),
		mx,
	)

	http.ListenAndServe(":8080", mx)
}
