package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"git.sr.ht/~icikowski/ldap-image-proxy/internal/ldap"
)

type ldapHandler struct {
	ldapClient ldap.Client
	ttl        time.Duration
	cache      *photoCache
}

func newLDAPHandler(ldapClient ldap.Client, ttl time.Duration, cacheDir string) *ldapHandler {
	effectiveTTL := max(0, ttl)

	return &ldapHandler{
		ldapClient: ldapClient,
		ttl:        effectiveTTL,
		cache:      &photoCache{dir: cacheDir, ttl: effectiveTTL},
	}
}

func (h *ldapHandler) bind(r chi.Router) {
	r.Get("/{username}", h.userPhoto)
}

func (h *ldapHandler) userPhoto(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	l := zerolog.Ctx(r.Context()).With().Str("requested_username", username).Logger()

	size := requestedImageSize(r)
	cacheOpts := requestCacheOptionsFromHeaders(r)

	//nolint:nestif
	if cacheOpts.noStore || h.cache.ttl == 0 {
		l.Debug().Int("size", size).Bool("no_store", cacheOpts.noStore).Msg("photo cache skipped")
	} else {
		if photo, ok, err := h.cache.load(username, size, cacheOpts.maxAge); err != nil {
			l.Error().Err(err).Msg("failed to read photo cache")
		} else if ok {
			l.Debug().Int("size", size).Msg("photo cache hit")
			if err := h.writePhoto(w, photo, cacheOpts.noStore); err != nil {
				l.Error().Err(err).Msg("failed to write response")
			}
			return
		} else {
			l.Debug().Int("size", size).Msg("photo cache miss")
		}
	}

	exists, err := h.ldapClient.UserExists(username)
	if err != nil {
		l.Error().Err(err).Msg("failed to check if user exists")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !exists {
		l.Warn().Msg("user not found")
		http.NotFound(w, r)
		return
	}

	photo, err := h.ldapClient.UserPhoto(username)
	if err != nil {
		l.Error().Err(err).Msg("failed to retrieve user photo from LDAP")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if size > 0 && len(photo) > 0 {
		photo, err = resizePhoto(photo, size)
		if err != nil {
			l.Error().Err(err).Msg("failed to resize user photo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if !cacheOpts.noStore {
		if err := h.cache.store(username, size, photo); err != nil {
			l.Error().Err(err).Msg("failed to write photo cache")
		}
	}

	if err := h.writePhoto(w, photo, cacheOpts.noStore); err != nil {
		l.Error().Err(err).Msg("failed to write response")
	}
}

func (h *ldapHandler) writePhoto(w http.ResponseWriter, photo []byte, noCache bool) error {
	w.Header().Set("Content-Type", "image/jpeg")
	if noCache {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
	} else {
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(int(h.ttl.Seconds())))
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(photo)
	return err
}
