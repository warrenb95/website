package http

import (
	"net/http"
)

func (s *Server) Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithField("endpoint", r.URL).WithContext(r.Context())

		logger.Info("Request start")
		defer logger.Info("Request end")

		h.ServeHTTP(w, r)
	})
}
