package handler

import (
	"encoding/json"
	"net/http"
)

type CaptureRequest struct {
	ID       string `json:"id" validate:"required"`
	Amount   int    `json:"amount" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type CaptureResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Remaining int    `json:"remaining,omitempty"`
	Currency  string `json:"currency,omitempty"`
}

func CaptureEndpoint(gateway Gateway, validator Validator, errorResp ErrorResponse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := CaptureRequest{}
		response := CaptureResponse{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = err.Error()
			response.Success = false
			errorResp(w, response, http.StatusBadRequest)
		}

		remain, currency, err := gateway.CaptureAmount(request.ID, request.Amount, request.Currency)

		if err == nil {
			response.Success = true
			response.Remaining = remain
			response.Currency = currency
		} else {
			response.Message = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
