package handler

import (
	"encoding/json"
	"net/http"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
)

type AuthorizeRequest struct {
	NameOnCard  string `json:"name_on_card" validate:"required"`
	CardNumber  string `json:"card_number" validate:"required"`
	ExpireYear  int    `json:"expire_year" validate:"required"`
	ExpireMonth int    `json:"expire_month" validate:"required"`
	CVV         string `json:"cvv" validate:"required"`
	Amount      int    `json:"amount" validate:"required"`
	Currency    string `json:"currency" validate:"required"`
}

type AuthorizeResponse struct {
	ID       string `json:"id,omitempty"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Amount   int    `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func AuthorizeEndpoint(gateway Gateway, logger Logger, validator Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := AuthorizeRequest{}
		response := AuthorizeResponse{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest, logger)
			return
		}

		err = validator(&request)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest, logger)
			return
		}

		token, err := gateway.AuthorizePayment(gtway.AuthorizeOptions{
			NameOnCard:  request.NameOnCard,
			CardNumber:  request.CardNumber,
			ExpireYear:  request.ExpireYear,
			ExpireMonth: request.ExpireMonth,
			CVV:         request.CVV,
			Amount:      request.Amount,
			Currency:    request.Currency,
		})
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest, logger)
			return
		}

		response.ID = token
		response.Success = true
		response.Amount = request.Amount
		response.Currency = request.Currency
		sendResponse(w, response, 200, logger)
	})
}
