package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"git.sr.ht/~icikowski/ldap-image-proxy/internal/ldap"
)

// New creates a new HTTP handler for the LDAP image proxy service.
func New(
	ldapClient ldap.Client,
	ttl time.Duration,
	cacheDir string,
	log zerolog.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.CleanPath,
		middleware.RequestID,
		zerologMiddleware(log),
		middleware.Recoverer,
	)

	handler := newLDAPHandler(ldapClient, ttl, cacheDir)
	handler.bind(r)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}
