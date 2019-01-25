package middleware

import (
	"encoding/json"
	"net/http"
)

// HealthChecker checks health.
type HealthChecker interface {
	CheckHealth() bool
}

// HealthCheckedDependency is a microservice dependency, which is registered and health checked.
type HealthCheckedDependency struct {
	Name     string
	Critical bool
	Checker  HealthChecker
}

// HealthCheckedDependencies are each of the dependencies which are needed to be checked in order to
// be able to say that service is completely healthy.
var HealthCheckedDependencies = []*HealthCheckedDependency{}

// AddHealthCheckedDependency adds a health checked dependency.
func AddHealthCheckedDependency(name string, critical bool, checker HealthChecker) {
	HealthCheckedDependencies = append(HealthCheckedDependencies, &HealthCheckedDependency{
		Name:     name,
		Critical: critical,
		Checker:  checker,
	})
}

// HealthHandler is simple handler for /health endpoint which does health checking and respond accordingly.
func HealthHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" && (r.Method == "GET" || r.Method == "HEAD") {
			checkHealth(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	var hasCriticalFailure, hasFailure bool
	results := make(map[string]bool)

	for _, dep := range HealthCheckedDependencies {
		healthy := dep.Checker.CheckHealth()
		if !healthy {
			hasFailure = true
			hasCriticalFailure = (hasCriticalFailure || dep.Critical)
		}
		results[dep.Name] = healthy
	}

	w.Header().Set("Content-Type", "application/json")

	if hasFailure {
		if hasCriticalFailure {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	b, _ := json.Marshal(results)
	w.Write(b)
}
