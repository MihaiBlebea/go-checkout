package server

import (
	"fmt"
	"log"

	"net/http"
	"os"
	"time"

	"github.com/MihaiBlebea/go-checkout/server/handler"
	"github.com/MihaiBlebea/go-checkout/server/validate"
	"github.com/gorilla/mux"
)

const prefix = "/api/v1/"

func NewServer(gateway Gateway, logger Logger) {
	r := mux.NewRouter().
		PathPrefix(prefix).
		Subrouter()

	r.Handle("/health-check", handler.HealthEndpoint(logger)).
		Methods("GET")

	r.Handle("/authorize", handler.AuthorizeEndpoint(gateway, logger, validate.Validate)).
		Methods("POST")

	r.Handle("/capture", handler.CaptureEndpoint(gateway, logger, validate.Validate)).
		Methods("POST")

	r.Handle("/void", handler.VoidEndpoint(gateway, logger, validate.Validate)).
		Methods("POST")

	r.Handle("/refund", handler.RefundEndpoint(gateway, logger, validate.Validate)).
		Methods("POST")

	r.Handle("/transactions", handler.ListEndpoint(gateway, logger, validate.Validate)).
		Methods("GET")

	r.Use(loggerMiddleware(logger))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
