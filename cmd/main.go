package main

import (
	"context"
	"errors"
	"log"
	"moneyTransfer/api"
	"moneyTransfer/api/handler"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/internal/queue"
	"moneyTransfer/internal/repository"
	"moneyTransfer/internal/repository/postgres"
	"moneyTransfer/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "moneyTransfer/docs"
)

// @title Money Transfer API
// @version 1.0
// @description REST API for transferring money between users
// @host localhost:8080
// @BasePath /
func main() {
	logger.Init()
	logger.Log.Info("Logger initialized")

	port := os.Getenv("SERVER_PORT")

	db, err := postgres.NewPostgresClient()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	var transferRepo contracts.TransferRepository = repository.NewTransferRepository(db)
	var userRepo contracts.UserRepository = repository.NewUserRepository(db)

	transferService := service.NewTransferService(transferRepo, userRepo, logger.Log)
	userService := service.NewUserService(userRepo, logger.Log)

	transferController := handler.NewTransferController(transferService, logger.Log)
	userController := handler.NewUserController(userService, logger.Log)

	router := api.InitRouter(transferController, userController)

	queue.StartWorker(userRepo, transferRepo, logger.Log)

	/** Graceful shutdown
	/- syscall.SIGTERM (kill -15, the default signal for docker stop)
	- syscall.SIGINT (kill -2, Ex: ctrl+c for testing on local machine)
	**/
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		logger.Log.Info("HTTP server started on :8080")
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan // wait till we get a signal from the channel (CTRL + C, docker stop, etc.)
	logger.Log.Info("Received shutdown signal, terminating...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	logger.Log.Info("Server shutdown gracefully")
}
