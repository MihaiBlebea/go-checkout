package handler

import (
	"encoding/json"
	"net/http"
)

type VoidRequest struct {
	ID string `json:"id" validate:"required"`
}

type VoidResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Balance  int    `json:"balance,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func VoidEndpoint(gateway Gateway, validator Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := VoidRequest{}
		response := VoidResponse{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		balance, currency, err := gateway.VoidTransaction(request.ID)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, &response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.Balance = balance
		response.Currency = currency
		sendResponse(w, &response, http.StatusOK)
	})
}
