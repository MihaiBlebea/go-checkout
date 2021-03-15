package server

import (
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/MihaiBlebea/go-checkout/server/handler"
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

	r.Handle("/authorize", loggerMiddleware(handler.AuthorizeEndpoint(gateway))).
		Methods("POST")

	r.Handle("/capture", loggerMiddleware(handler.CaptureEndpoint(gateway))).
		Methods("POST")

	r.Handle("/void", loggerMiddleware(handler.VoidEndpoint())).
		Methods("POST")

	r.Handle("/refund", loggerMiddleware(handler.RefundEndpoint())).
		Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
