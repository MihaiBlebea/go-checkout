package handler

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func RefundEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := AuthorizeRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := AuthorizeResponse{
			ID:       uuid.NewV4().String(),
			Success:  true,
			Amount:   200,
			Currency: "GBP",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
