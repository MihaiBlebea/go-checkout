package handler

import (
	"encoding/json"
	"net/http"

	gtway "github.com/MihaiBlebea/go-checkout/gateway"
)

type AuthorizeRequest struct {
	NameOnCard  string `json:"name_on_card"`
	CardNumber  string `json:"card_number"`
	ExpireYear  int    `json:"expire_year"`
	ExpireMonth int    `json:"expire_month"`
	CVV         string `json:"cvv"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
}

type AuthorizeResponse struct {
	ID       string `json:"id,omitempty"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Amount   int    `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func AuthorizeEndpoint(gateway Gateway) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := AuthorizeRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		response := AuthorizeResponse{ID: token}

		if err == nil {
			response.Success = true
			response.Amount = request.Amount
			response.Currency = request.Currency
		} else {
			response.Message = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
