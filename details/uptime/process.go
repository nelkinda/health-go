package uptime

import (
	"github.com/nelkinda/health-go"
	"net/http"
	"time"
)

type process struct {
	start time.Time
}

func (u *process) HealthDetails() map[string][]health.Details {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return map[string][]health.Details{
		"uptime": {
			{
				ComponentType: "process",
				ObservedValue: time.Now().UTC().Sub(u.start).Seconds(),
				ObservedUnit: "s",
				Status:health.Pass,
				Time: now,
			},
		},
	}
}

func (*process) AuthorizeHealth(r *http.Request) bool {
	return true
}

// Process returns a DetailsProvider for health details about the process uptime.
// Note that it does not really return the process uptime, but the time since calling this function.
func Process() health.DetailsProvider {
	return &process{start: time.Now().UTC()}
}

