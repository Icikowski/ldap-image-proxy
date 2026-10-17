package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func zerologMiddleware(l zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := l.With().Str("request_id", middleware.GetReqID(r.Context())).Logger()

			l.Debug().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Str("remote_addr", r.RemoteAddr).
				Msg("incoming request")

			start := time.Now()
			next.ServeHTTP(w, r.WithContext(l.WithContext(r.Context())))
			duration := time.Since(start)

			l.Info().Dur("duration", duration).Msg("request completed")
		})
	}
}
