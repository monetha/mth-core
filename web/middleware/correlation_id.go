package middleware

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	webcontext "gitlab.com/monetha/mth-core/web/context"
	"gitlab.com/monetha/mth-core/web/header"
)

// CorrelationIDHandler adds or forward mth-correlation-id header
func CorrelationIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := getOrCreateNewCorrelationID(r)
		w.Header().Set(header.HeaderKeyCorrelationID, correlationID)
		ctx := webcontext.WithCorrelationID(r.Context(), correlationID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getOrCreateNewCorrelationID(r *http.Request) string {
	if cid := header.CorrelationID(r); cid != nil {
		return *cid
	}
	return uuid.NewV4().String()
}
