package handler

import (
	"encoding/json"
	"net/http"
)

func HealthEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			OK bool `json:"ok"`
		}{
			OK: true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
