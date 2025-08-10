package main

import (
	"cryptoserver/internal/repository/mapdb"
	"cryptoserver/internal/rest"
	"cryptoserver/internal/service"
	"net/http"
)

func main() {
	mx := http.NewServeMux()
	rest.RegisterCryptoHandler(
		service.NewService(mapdb.NewCryptoRepository(), mapdb.NewPriceRepository()),
		mx,
	)
	http.ListenAndServe(":8080", mx)
}
