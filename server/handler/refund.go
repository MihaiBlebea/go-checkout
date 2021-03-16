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
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Remaining int    `json:"remaining,omitempty"`
	Currency  string `json:"currency,omitempty"`
}

func RefundEndpoint(gateway Gateway, validator Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := RefundRequest{}
		response := RefundResponse{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = err.Error()
			response.Success = false
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		remain, currency, err := gateway.RefundAmount(request.ID, request.Amount, request.Currency)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, &response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.Remaining = remain
		response.Currency = currency
		sendResponse(w, &response, http.StatusOK)
	})
}
