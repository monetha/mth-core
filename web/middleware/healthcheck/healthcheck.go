package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitlab.com/monetha/mth-core/log"
)

var (
	// Timeout is check timeout.
	Timeout = time.Minute
	// MaxFailureInARow is the number for when a dependency is considered broken/down.
	MaxFailureInARow = 3
	// dependencies are each of the dependencies which are needed to be checked in order to
	// be able to say that service is completely healthy.
	dependencies []*dependency
	// enabledDependencies are assigned to dependencies on start.
	enabledDependencies = []*dependency{}
)

// HealthChecker checks health.
type HealthChecker interface {
	CheckHealth() bool
}

// dependency is a microservice dependency, which is registered and health checked.
type dependency struct {
	Name     string
	Critical bool
	Checker  HealthChecker
	Interval time.Duration

	FailureInARow int

	sync.RWMutex
}

// AddDependency adds a health checked dependency.
func AddDependency(name string, critical bool, checker HealthChecker, interval time.Duration) {
	enabledDependencies = append(enabledDependencies, &dependency{
		Name:     name,
		Critical: critical,
		Checker:  checker,
		Interval: interval,
	})
}

// Handler is simple handler for /health endpoint which reports with health status of dependencies.
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" && (r.Method == "GET" || r.Method == "HEAD") {
			writeHealthStatus(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func writeHealthStatus(w http.ResponseWriter, r *http.Request) {
	var hasCriticalFailure, hasFailure bool
	results := make(map[string]bool)

	for _, dep := range dependencies {
		dep.RLock()
		consideredHealthy := dep.failuresAreNegligible()
		if !consideredHealthy {
			hasFailure = true
			hasCriticalFailure = (hasCriticalFailure || dep.Critical)
		}
		results[dep.Name] = consideredHealthy
		dep.RUnlock()
	}

	w.Header().Set("Content-Type", "application/json")

	switch {
	case hasFailure:
		if hasCriticalFailure {
			w.WriteHeader(http.StatusServiceUnavailable)
			break
		}
		w.WriteHeader(http.StatusInternalServerError)

	default:
		w.WriteHeader(http.StatusOK)
	}

	b, _ := json.Marshal(results)
	w.Write(b)
}

func (dep *dependency) failuresAreNegligible() bool {
	return dep.FailureInARow < MaxFailureInARow
}

func (dep *dependency) applyHealthCheckResult(healthyNow bool) {
	if healthyNow {
		dep.FailureInARow = 0
		return
	}

	if dep.failuresAreNegligible() {
		dep.FailureInARow++ // Increment it so maybe it becomes non-negligible soon
	}
}

// Start starts async health check.
func Start(ctx context.Context) {
	dependencies = enabledDependencies

	for _, dep := range dependencies {
		dep.check()
	}

	logFailingDeps()

	for _, dep := range dependencies {
		go dep.runAsync(ctx)
	}
}

func (dep *dependency) check() {
	healthyNow := dep.Checker.CheckHealth()
	dep.Lock()
	dep.applyHealthCheckResult(healthyNow)
	dep.Unlock()
}

func (dep *dependency) checkAndNotify(c chan struct{}) {
	dep.check()
	close(c)
}

func (dep *dependency) runAsync(ctx context.Context) {
	ticker := time.NewTicker(dep.Interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		c := make(chan struct{})
		go dep.checkAndNotify(c)

		select {
		case <-ctx.Done():
			return
		case <-c:
			logFailingDeps()
		case <-time.After(Timeout):
			continue
		}
	}
}

func logFailingDeps() {
	failingDeps := make([]string, 0)

	for _, dep := range dependencies {
		dep.RLock()
		consideredHealthy := dep.failuresAreNegligible()
		if !consideredHealthy {
			failingDeps = append(failingDeps, dep.Name)
		}
		dep.RUnlock()
	}

	if len(failingDeps) > 0 {
		log.With().Error(fmt.Sprintf("some of service dependencies are failing: %s", strings.Join(failingDeps, ", ")))
	}
}
