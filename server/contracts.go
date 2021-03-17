package server

import (
	gtway "github.com/MihaiBlebea/go-checkout/gateway"
)

type Gateway interface {
	AuthorizePayment(options gtway.AuthorizeOptions) (string, error)
	CaptureAmount(id string, amount int, currency string) (int, string, error)
	VoidTransaction(id string) (int, string, error)
	RefundAmount(id string, amount int, currency string) (int, string, error)
	ListTransactions() (transactions []gtway.Transaction)
}

type Logger interface {
	Info(args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}
