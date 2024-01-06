package app

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/v7ktory/fullstack/internal/repository"
	"github.com/v7ktory/fullstack/internal/service"
	"github.com/v7ktory/fullstack/internal/transport/rest"
	"github.com/v7ktory/fullstack/pkg/database/postgres"
	"github.com/v7ktory/fullstack/pkg/hasher"
	"github.com/v7ktory/fullstack/pkg/logger"
	"github.com/v7ktory/fullstack/pkg/token"
	"go.uber.org/zap"
)

func Run() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logger := logger.NewLogger()

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return
	}
	hash := hasher.NewHasher(os.Getenv("SALT"))
	jwt := token.NewJWTService(os.Getenv("JWT_SECRET_KEY"))
	repos := repository.NewRepository(db)
	service := service.NewService(repos, hash, *jwt)
	handler := rest.NewHandler(service, jwt)

	if err := handler.InitRoutes().Run(os.Getenv("SERVER_ADDRESS")); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Failed to run server", zap.Error(err))
	}

}
