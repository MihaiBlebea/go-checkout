package handler

import (
	"encoding/json"
	"net/http"
)

type CaptureRequest struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type CaptureResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Amount   int    `json:"remaining_amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func CaptureEndpoint(gateway Gateway) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := CaptureRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		remaining, _, err := gateway.CaptureAmount(request.ID, request.Amount, request.Currency)

		response := CaptureResponse{}

		if err == nil {
			response.Success = true
			response.Amount = remaining
			response.Currency = request.Currency
		} else {
			response.Message = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
