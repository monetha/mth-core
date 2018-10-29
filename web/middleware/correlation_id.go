package middleware

import (
	"net/http"

	"github.com/satori/go.uuid"
	webcontext "gitlab.com/monetha/mth-serva-bazo/web/context"
)

const (
	// HeaderCorrelationID is the name of the correlation id header
	HeaderCorrelationID = "mth-correlation-id"
)

var (
	headerCorrelationIDCanonical = http.CanonicalHeaderKey(HeaderCorrelationID)
)

// CorrelationIDHandler adds or forward mth-correlation-id header
func CorrelationIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := getOrCreateNewCorrelationID(r)
		w.Header().Set(HeaderCorrelationID, correlationID)
		ctx := webcontext.WithCorrelationID(r.Context(), correlationID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getOrCreateNewCorrelationID(r *http.Request) string {
	if val, ok := r.Header[headerCorrelationIDCanonical]; ok {
		return val[0]
	}
	return uuid.NewV4().String()
}
