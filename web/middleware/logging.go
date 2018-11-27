package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"gitlab.com/monetha/mth-core/log"
	"gitlab.com/monetha/mth-core/web"
	webcontext "gitlab.com/monetha/mth-core/web/context"
	"gitlab.com/monetha/mth-core/web/header"
	"go.uber.org/zap"
)

var sensitiveHeaderKeys = map[string]struct{}{
	http.CanonicalHeaderKey("Authorization"): struct{}{},
}

// LoggingHandler is a middleware that will write the log to 'out' writer.
func LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		path := r.URL.Path
		raw := r.URL.RawQuery

		cacheReader := newLimitCacheReader(r.Body, 4096)
		r.Body = cacheReader

		// call inner handler
		lrw := web.NewLogStatusReponseWriter(w)
		h.ServeHTTP(lrw, r)

		end := time.Now()
		latency := end.Sub(start)

		method := r.Method
		statusCode := lrw.Status()

		correlationID := webcontext.CorrelationID(r.Context())
		clientIP := header.ClientIP(r)
		authClaims := header.AuthClaims(r)

		if raw != "" {
			path = path + "?" + raw
		}

		l := log.With(
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status_code", statusCode),
			zap.String("correlation_id", correlationID),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
		)
		l = l.With(log.FieldsFrom(adjustKeysWithPrefix(authClaims, "c"))...)

		reqHeaderFields := getHeaderFields(r.Header, sensitiveHeaderKeys, "ih")
		reqHeaderLoggingFields := log.FieldsFrom(reqHeaderFields)
		l = l.With(reqHeaderLoggingFields...)

		respHeaderFields := getHeaderFields(lrw.Header(), sensitiveHeaderKeys, "oh")
		respHeaderLoggingFields := log.FieldsFrom(respHeaderFields)
		l = l.With(respHeaderLoggingFields...)

		if statusCode >= 500 {
			var rb bytes.Buffer
			reqBytes, _ := httputil.DumpRequest(r, false)
			rb.Write(reqBytes)
			rb.Write(cacheReader.Bytes())

			l.Error("[HTTP]", zap.String("payload", rb.String()))
		} else {
			l.Info("[HTTP]")
		}
	})
}

func adjustKeysWithPrefix(m map[string]interface{}, prefix string) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range m {
		ret[prefix+"_"+k] = v
	}
	return ret
}

func getHeaderFields(h http.Header, exclude map[string]struct{}, prefix string) (m map[string]interface{}) {
	m = make(map[string]interface{})
	for key := range h {
		if _, ok := exclude[http.CanonicalHeaderKey(key)]; ok {
			continue
		}
		m[prefix+"_"+strings.ToLower(key)] = h.Get(key)
	}
	return
}

func newLimitCacheReader(r io.ReadCloser, limit int64) *limitCacheReader {
	return &limitCacheReader{r: r, n: limit}
}

// limitCacheReader caches first `limit` bytes from underlying reader
type limitCacheReader struct {
	r   io.ReadCloser
	n   int64 // max bytes remaining
	buf bytes.Buffer
}

func (c *limitCacheReader) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	if c.n <= 0 || n <= 0 {
		return
	}
	var r int64
	if int64(n) > c.n {
		r = c.n
	} else {
		r = int64(n)
	}
	c.buf.Write(p[:r])
	c.n -= r
	return
}

func (c *limitCacheReader) Close() error { return c.r.Close() }

func (c *limitCacheReader) Bytes() []byte { return c.buf.Bytes() }
