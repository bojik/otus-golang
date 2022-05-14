package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func newLoggingMiddleware(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := &responseRecorder{
				statusCode:     http.StatusOK,
				ResponseWriter: w,
			}
			t1 := time.Now()
			next.ServeHTTP(response, r)
			timing := time.Since(t1)
			l.Info(
				fmt.Sprintf(
					`%s [%s] %s %s %s %d %.2f "%s"`,
					r.RemoteAddr,
					time.Now().Format(time.RFC3339),
					r.Method,
					r.RequestURI,
					r.Proto,
					response.statusCode,
					float32(timing/time.Second),
					r.UserAgent(),
				),
			)
		})
	}
}
