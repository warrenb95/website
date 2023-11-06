package http

import (
	"log"
	"net/http"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request start: %s", r.URL)
		defer log.Printf("Request end: %s", r.URL)

		h.ServeHTTP(w, r)
	})
}
