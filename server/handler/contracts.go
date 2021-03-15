package handler

import gtway "github.com/MihaiBlebea/go-checkout/gateway"

type Gateway interface {
	AuthorizePayment(options gtway.AuthorizeOptions) (string, error)
	CaptureAmount(id string, amount int, currency string) (int, string, error)
	VoidTransaction(id string) (int, string, error)
}
