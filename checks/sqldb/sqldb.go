package mysqldb

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/nelkinda/health-go"
)

type sqldb struct {
	client      *sql.DB
	componentID string
	timeout     time.Duration
	threshold   time.Duration
}

func (m *sqldb) HealthChecks() map[string][]health.Checks {
	start := time.Now().UTC()
	startTime := start.Format(time.RFC3339Nano)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	checks := health.Checks{
		ComponentID: m.componentID,
		Time:        startTime,
	}

	err := m.client.PingContext(ctx)
	if err != nil {
		checks.Output = err.Error()
		checks.Status = health.Fail

		return map[string][]health.Checks{"mysqldb:responseTime": {checks}}
	}

	end := time.Now().UTC()
	responseTime := end.Sub(start)
	checks.ObservedValue = responseTime.Nanoseconds()
	checks.ObservedUnit = "ns"
	if responseTime > m.threshold {
		checks.Status = health.Warn
	} else {
		checks.Status = health.Pass
	}

	return map[string][]health.Checks{"mongodb:responseTime": {checks}}
}

// AuthorizeHealth return authorize flag status for health checks
func (m *sqldb) AuthorizeHealth(*http.Request) bool {
	return true
}

// Health returns a ChecksProvider for health checks about the process uptime.
// Note that it does not really return the process uptime, but the time since calling this function.
func Health(componentID string, client *sql.DB, timeout time.Duration, threshold time.Duration) health.ChecksProvider {
	return &sqldb{componentID: componentID, client: client, timeout: timeout, threshold: threshold}
}
