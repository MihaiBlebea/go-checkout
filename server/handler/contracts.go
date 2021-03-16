package handler

import (
	"net/http"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
)

type Gateway interface {
	AuthorizePayment(options gtway.AuthorizeOptions) (string, error)
	CaptureAmount(id string, amount int, currency string) (int, string, error)
	VoidTransaction(id string) (int, string, error)
	RefundAmount(id string, amount int, currency string) (int, string, error)
	ListTransactions() (transactions []gtway.Transaction)
}

type Validator func(t interface{}) error

type ErrorResponse func(w http.ResponseWriter, err interface{}, code int)
