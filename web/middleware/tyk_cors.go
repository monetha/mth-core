package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// TykCors is a middleware that handles CORS headers in a fashion that works with Tyk
// Only handles OPTIONS calls, because Tyk handles the rest
// NOTE: Use this middleware when you have enabled CORS handling in Tyk and "Options passthrogh" is enabled
type TykCors struct {
	cors *cors.Cors
}

// NewTykCors creates new TykCors handler
func NewTykCors(c *cors.Cors) *TykCors {
	return &TykCors{
		cors: c,
	}
}

// Handler apply the CORS specification on the request, and add relevant CORS headers
// as necessary
func (tc *TykCors) Handler(h http.Handler) http.Handler {
	corsHandlerFunc := tc.cors.Handler(h)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Handle only OPTIONS, because rest is handled by Tyk
		if r.Method == http.MethodOptions {
			corsHandlerFunc.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}
