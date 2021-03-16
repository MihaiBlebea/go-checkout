package handler

import (
	"encoding/json"
	"net/http"
)

type RefundRequest struct {
	ID       string `json:"id" validate:"required"`
	Amount   int    `json:"amount" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type RefundResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Amount   int    `json:"remaining_amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func RefundEndpoint(gateway Gateway, validator Validator, errorResp ErrorResponse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := RefundRequest{}
		response := AuthorizeResponse{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = err.Error()
			response.Success = false
			errorResp(w, response, http.StatusBadRequest)
		}

		remaining, currency, err := gateway.RefundAmount(request.ID, request.Amount, request.Currency)

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
