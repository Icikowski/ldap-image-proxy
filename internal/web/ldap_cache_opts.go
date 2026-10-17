package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	requestedImageSizeMin     = 16
	requestedImageSizeDefault = 256
	requestedImageSizeMax     = 2048
)

type requestCacheOptions struct {
	noStore bool
	maxAge  *time.Duration
}

func requestCacheOptionsFromHeaders(r *http.Request) requestCacheOptions {
	out := requestCacheOptions{}

	if requestHeaderContainsDirective(r.Header.Values("Cache-Control"), "no-cache") {
		out.noStore = true
		return out
	}
	if requestHeaderContainsDirective(r.Header.Values("Cache-Control"), "no-store") {
		out.noStore = true
		return out
	}
	if requestHeaderContainsDirective(r.Header.Values("Pragma"), "no-cache") {
		out.noStore = true
		return out
	}

	if maxAge, ok := requestHeaderMaxAge(r.Header.Values("Cache-Control")); ok {
		out.maxAge = &maxAge
	}

	return out
}

func requestHeaderContainsDirective(values []string, directive string) bool {
	for _, value := range values {
		for part := range strings.SplitSeq(value, ",") {
			if strings.EqualFold(strings.TrimSpace(part), directive) {
				return true
			}
		}
	}

	return false
}

func requestHeaderMaxAge(values []string) (time.Duration, bool) {
	for _, value := range values {
		for part := range strings.SplitSeq(value, ",") {
			part = strings.TrimSpace(part)
			if !strings.HasPrefix(strings.ToLower(part), "max-age=") {
				continue
			}

			raw := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(part), "max-age="))
			seconds, err := strconv.ParseInt(raw, 10, 64)
			if err != nil || seconds < 0 {
				return 0, false
			}

			return time.Duration(seconds) * time.Second, true
		}
	}

	return 0, false
}

func requestedImageSize(r *http.Request) int {
	raw := r.URL.Query().Get("s")
	if raw == "" {
		raw = r.URL.Query().Get("size")
	}
	if raw == "" {
		return requestedImageSizeDefault
	}
	raw = strings.TrimSpace(raw)
	if strings.EqualFold(raw, "orig") || strings.EqualFold(raw, "original") {
		return 0
	}

	size, err := strconv.Atoi(raw)
	if err != nil {
		return requestedImageSizeDefault
	}
	switch {
	case size < requestedImageSizeMin:
		size = requestedImageSizeMin
	case size > requestedImageSizeMax:
		size = requestedImageSizeMax
	}

	return size
}
