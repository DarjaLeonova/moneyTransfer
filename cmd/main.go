package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"log"
	"moneyTransfer/api/handler"
	"moneyTransfer/internal/domain/service"
	"moneyTransfer/internal/repository"
	"moneyTransfer/internal/repository/postgres"
	"net/http"

	_ "moneyTransfer/docs"
)

// @title Money Transfer API
// @version 1.0
// @description REST API for transferring money between users
// @host localhost:8080
// @BasePath /

func main() {
	db, err := postgres.NewPostgresClient()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	var transferRepo repository.TransferRepository = postgres.NewTransferRepository(db)
	var userRepo repository.UserRepository = postgres.NewUserRepository(db)

	transferService := service.NewTransferService(transferRepo, userRepo)
	userService := service.NewUserService(userRepo)

	transferController := handler.NewTransferController(transferService)
	userController := handler.NewUserController(userService)

	router := mux.NewRouter()

	router.HandleFunc("/transfers/{userId}", transferController.GetTransactionsByUserID).Methods("GET")
	router.HandleFunc("/transfers", transferController.CreateTransaction).Methods("POST")
	router.HandleFunc("/balance/{userId}", userController.GetUserBalance).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
