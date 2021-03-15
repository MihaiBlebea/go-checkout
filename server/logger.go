package server

import (
	"net/http"
	"os"

	logrus "github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"host":   r.Host,
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("Incoming request")

		next.ServeHTTP(w, r)
	})
}
