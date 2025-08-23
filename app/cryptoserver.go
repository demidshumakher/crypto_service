package main

import (
	"cryptoserver/internal/repository/postgresql"
	"cryptoserver/internal/rest"
	"cryptoserver/internal/rest/middleware"
	"cryptoserver/internal/service"
	"cryptoserver/pkg/coingecko"
	"cryptoserver/pkg/jwt"
	"cryptoserver/pkg/trigger"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	jwtSecret = []byte("my_secret_key")
	apiKey    = "CG-VThWUzH2txUEJ1eSaYFga8QK"
	baseURL   = "https://api.coingecko.com/api/v3"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Настройка подключения к PostgreSQL
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "cryptodb")
	sslMode := getEnv("DB_SSL_MODE", "disable")

	// Формирование строки подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	// Подключение к базе данных
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось проверить соединение с базой данных: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL")

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

	// Используем PostgreSQL репозитории вместо mapdb
	cryptoRepository := postgresql.NewCryptoRepository(db)
	userRepository := postgresql.NewUserRepository(db)

	rest.NewCryptoHandler(
		service.NewCryptoService(
			cryptoRepository,
			gecko,
		),
		authRouter,
	)

	rest.NewUserHandler(
		service.NewUserService(jwtConfig, userRepository),
		publicRouter,
	)

	rest.NewScheduleHandler(
		service.NewScheduleService(cryptoRepository, gecko, triggerConfig),
		authRouter,
	)

	log.Println("Сервер запущен на порту :8080")
	err = http.ListenAndServe(":8080", mx)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
