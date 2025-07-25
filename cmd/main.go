package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
	"log"
	"moneyTransfer/api/handler"
	"moneyTransfer/internal/domain/contracts"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/internal/repository"
	"moneyTransfer/internal/repository/postgres"
	"moneyTransfer/pkg/metrics"
	"net/http"
	"os"

	_ "moneyTransfer/docs"
)

// @title Money Transfer API
// @version 1.0
// @description REST API for transferring money between users
// @host localhost:8080
// @BasePath /

func main() {
	port := os.Getenv("SERVER_PORT")

	db, err := postgres.NewPostgresClient()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	var transferRepo contracts.TransferRepository = repository.NewTransferRepository(db)
	var userRepo contracts.UserRepository = repository.NewUserRepository(db)

	transferService := service.NewTransferService(transferRepo, userRepo)
	userService := service.NewUserService(userRepo)

	transferController := handler.NewTransferController(transferService)
	userController := handler.NewUserController(userService)

	router := mux.NewRouter()

	router.HandleFunc("/transfers/{userId}", transferController.GetTransactionsByUserID).Methods("GET")
	router.HandleFunc("/transfers", transferController.CreateTransaction).Methods("POST")
	router.HandleFunc("/balance/{userId}", userController.GetUserBalance).Methods("GET")

	router.Use(metrics.NewPrometheusMiddleware())
	router.Handle("/metrics", promhttp.Handler())

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Server listening on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
