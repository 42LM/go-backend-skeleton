package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// loggingMiddleware logs the body of incoming HTTP requests.
// It reads the request body, logs it, and then replaces the body
// with a new reader so the next handler in the chain can read it.
func loggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			// restore body so next handler can read it
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			defer func(begin time.Time) {
				logger.Info("",
					"request_body", string(bodyBytes),
					"method", r.Method,
					"path", r.URL.Path,
					"took", float64(time.Since(begin))/1e6,
				)
			}(time.Now())

			next.ServeHTTP(w, r)
		})
	}
}
