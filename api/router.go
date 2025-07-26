package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/http-swagger"
	"moneyTransfer/api/handler"
	"moneyTransfer/pkg/metrics"
)

func InitRouter(
	transferController *handler.TransferController,
	userController *handler.UserController,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/transfers/{userId}", transferController.GetTransactionsByUserId).Methods("GET")
	router.HandleFunc("/transfers", transferController.CreateTransaction).Methods("POST")
	router.HandleFunc("/balance/{userId}", userController.GetUserBalance).Methods("GET")

	router.Use(metrics.NewPrometheusMiddleware())

	router.Handle("/metrics", promhttp.Handler())
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
