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

func CaptureEndpoint(gateway Gateway, logger Logger, validator Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := CaptureRequest{}
		response := CaptureResponse{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest, logger)
			return
		}

		remain, currency, err := gateway.CaptureAmount(request.ID, request.Amount, request.Currency)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest, logger)
			return
		}

		response.Success = true
		response.Remaining = remain
		response.Currency = currency
		sendResponse(w, response, http.StatusOK, logger)
	})
}
