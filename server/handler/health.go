package handler

import (
	"net/http"
)

func HealthEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			OK bool `json:"ok"`
		}{
			OK: true,
		}

		sendResponse(w, &response, http.StatusOK)
	})
}
