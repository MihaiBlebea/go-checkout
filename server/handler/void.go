package handler

import (
	"encoding/json"
	"net/http"
)

type VoidRequest struct {
	ID string `json:"id"`
}

type VoidResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Amount   int    `json:"remaining_amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func VoidEndpoint(gateway Gateway) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := VoidRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		remaining, currency, err := gateway.VoidTransaction(request.ID)

		response := VoidResponse{}

		if err == nil {
			response.Success = true
			response.Amount = remaining
			response.Currency = currency
		} else {
			response.Message = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
