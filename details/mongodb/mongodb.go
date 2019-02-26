// Provides health details for a MongoDB connection.
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
	componentId string
	client      *mongo.Client
}

func (m *mongodb) HealthDetails() map[string][]health.Details {
	start := time.Now().UTC()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return map[string][]health.Details{
			"mongodb:responseTime": {
				{
					ComponentId: m.componentId,
					Output:      err.Error(),
					Status:      health.Fail,
					Time:        start.Format(time.RFC3339Nano),
				},
			},
		}
	}
	end := time.Now().UTC()
	responseTime := end.Sub(start).Nanoseconds()
	return map[string][]health.Details{
		"mongodb:responseTime": {
			{
				ComponentId:   m.componentId,
				ObservedValue: responseTime,
				ObservedUnit:  "ns",
				Status:        health.Pass,
				Time:          start.Format(time.RFC3339Nano),
			},
		},
	}
}

func (*mongodb) AuthorizeHealth(r *http.Request) bool {
	return true
}

// Process returns a DetailsProvider for health details about the process uptime.
// Note that it does not really return the process uptime, but the time since calling this function.
func Health(componentId string, client *mongo.Client) health.DetailsProvider {
	return &mongodb{componentId: componentId, client: client}
}
