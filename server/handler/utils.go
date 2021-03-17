package handler

import (
	"encoding/json"
	"net/http"
)

func sendResponse(w http.ResponseWriter, resp interface{}, code int, logger Logger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, _ := json.Marshal(resp)

	log := logger.Info

	if code != 200 {
		log = logger.Error
	}

	log(string(b))
	w.Write(b)
}
