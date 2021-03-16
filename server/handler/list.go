package handler

import (
	"encoding/json"
	"net/http"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
)

type ListRequest struct {
}

type ListResponse struct {
	Success      bool                `json:"success"`
	Transactions []gtway.Transaction `json:"transactions"`
}

func ListEndpoint(gateway Gateway, validator Validator, errorResp ErrorResponse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ListResponse{}

		transactions := gateway.ListTransactions()

		response.Success = true
		response.Transactions = transactions

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
