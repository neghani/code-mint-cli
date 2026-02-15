package api

import "net/http"

type uaRoundTripper struct {
	next      http.RoundTripper
	userAgent string
}

func (u uaRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", u.userAgent)
	return u.next.RoundTrip(r)
}
