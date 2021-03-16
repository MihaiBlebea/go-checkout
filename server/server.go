package server

import (
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/MihaiBlebea/go-checkout/server/handler"
	"github.com/MihaiBlebea/go-checkout/server/resp"
	"github.com/MihaiBlebea/go-checkout/server/validate"
	"github.com/gorilla/mux"
)

const prefix = "/api/v1/"

func NewServer(gateway Gateway) {
	r := mux.NewRouter().
		PathPrefix(prefix).
		Headers("Content-Type", "application/json").
		Subrouter()

	r.Handle("/health", loggerMiddleware(handler.HealthEndpoint())).
		Methods("GET")

	r.Handle("/authorize", loggerMiddleware(handler.AuthorizeEndpoint(gateway, validate.Validate, resp.ErrorResponse))).
		Methods("POST")

	r.Handle("/capture", loggerMiddleware(handler.CaptureEndpoint(gateway, validate.Validate, resp.ErrorResponse))).
		Methods("POST")

	r.Handle("/void", loggerMiddleware(handler.VoidEndpoint(gateway, validate.Validate, resp.ErrorResponse))).
		Methods("POST")

	r.Handle("/refund", loggerMiddleware(handler.RefundEndpoint(gateway, validate.Validate, resp.ErrorResponse))).
		Methods("POST")

	r.Handle("/transactions", loggerMiddleware(handler.ListEndpoint(gateway, validate.Validate, resp.ErrorResponse))).
		Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
