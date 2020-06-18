// Package mongodb provides health checks for a MongoDB connection.
// This works on MongoDB as well as Microsoft Azure Cosmos DB.
package mongodb

import (
	"context"
	"github.com/nelkinda/health-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"time"
)

type mongodb struct {
	componentID string
	client      *mongo.Client
	timeout     time.Duration
	threshold   time.Duration
}

func (m *mongodb) HealthChecks(ctx context.Context) map[string][]health.Checks {
	start := time.Now().UTC()
	startTime := start.Format(time.RFC3339Nano)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	err := m.client.Ping(ctx, readpref.Primary())
	var checks = health.Checks{
		ComponentID: m.componentID,
		Time:        startTime,
	}
	if err != nil {
		checks.Output = err.Error()
		checks.Status = health.Fail
	} else {
		end := time.Now().UTC()
		responseTime := end.Sub(start)
		checks.ObservedValue = responseTime.Nanoseconds()
		checks.ObservedUnit = "ns"
		if responseTime > m.threshold {
			checks.Status = health.Warn
		} else {
			checks.Status = health.Pass
		}
	}
	return map[string][]health.Checks{"mongodb:responseTime": {checks}}
}

func (*mongodb) AuthorizeHealth(r *http.Request) bool {
	return true
}

// Health returns a ChecksProvider for health checks about the process uptime.
// Note that it does not really return the process uptime, but the time since calling this function.
func Health(componentID string, client *mongo.Client, timeout time.Duration, threshold time.Duration) health.ChecksProvider {
	return &mongodb{componentID: componentID, client: client, timeout: timeout, threshold: threshold}
}
