package main

import (
	"cryptoserver/internal/repository/mapdb"
	"cryptoserver/internal/rest"
	"cryptoserver/internal/rest/middleware"
	"cryptoserver/internal/service"
	"cryptoserver/pkg/coingecko"
	"cryptoserver/pkg/jwt"
	"cryptoserver/pkg/trigger"
	"net/http"
	"time"
)

var (
	jwtSecret = []byte("my_secret_key")
	apiKey    = "CG-VThWUzH2txUEJ1eSaYFga8QK"
	baseURL   = "https://api.coingecko.com/api/v3"
)

func main() {
	jwtConfig := jwt.JWTConfig{
		Secret:         jwtSecret,
		ExpirationTime: time.Now().Add(time.Hour * 24).Unix(),
	}

	gecko := coingecko.NewCoinGeckoClient(apiKey, baseURL)

	triggerConfig := trigger.TriggerCfg{
		IntervalSeconds: 30,
	}

	mx := http.NewServeMux()

	publicRouter := rest.NewRouter(mx)
	authRouter := rest.NewRouter(mx)
	authRouter.ApplyMiddleware(middleware.AuthMiddleware(jwtConfig))

	cryptoRepository := mapdb.NewCryptoRepository()

	rest.NewCryptoHandler(
		service.NewCryptoService(
			cryptoRepository,
			gecko,
		),
		authRouter,
	)

	rest.NewUserHandler(
		service.NewUserService(jwtConfig, mapdb.NewUserRepository()),
		publicRouter,
	)

	rest.NewScheduleHandler(
		service.NewScheduleService(cryptoRepository, gecko, triggerConfig),
		authRouter,
	)

	http.ListenAndServe(":8080", mx)
}
