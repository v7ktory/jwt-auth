package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	dbConfig := postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := postgres.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return
	}

	hash := hasher.NewHasher(os.Getenv("SALT"))
	jwt := token.NewJWTService(os.Getenv("JWT_SECRET_KEY"))

	repos := repository.NewRepository(db)
	service := service.NewService(repos, hash, *jwt)

	handler := rest.NewHandler(service, jwt)

	serverAddress := os.Getenv("SERVER_ADDRESS")
	server := &http.Server{
		Addr:    serverAddress,
		Handler: handler.InitRoutes(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to run server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed", zap.Error(err))
	}

	logger.Info("Server shutdown gracefully")
}
