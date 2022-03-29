package pkg

import (
	"net/http"

	"go.uber.org/ratelimit"
)

type RateLimit map[string]int

type rateLimitedTransport struct {
	http.RoundTripper
	rateLimiters map[string]ratelimit.Limiter
}

func (t rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if limiter, ok := t.rateLimiters[req.URL.Host]; ok {
		limiter.Take()
	}

	return t.RoundTripper.RoundTrip(req)
}

// NewHTTPClient instantiates a new *http.Client.
func NewHTTPClient(rateLimit RateLimit) *http.Client {
	rateLimiters := map[string]ratelimit.Limiter{}
	for host, limit := range rateLimit {
		rateLimiters[host] = ratelimit.New(limit)
	}

	return &http.Client{
		// Don't follow off-site redirects.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Hostname() != via[len(via)-1].URL.Hostname() {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: &rateLimitedTransport{
			http.DefaultTransport,
			rateLimiters,
		},
	}
}
