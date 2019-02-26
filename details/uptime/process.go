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

func Process() *process {
	return &process{start: time.Now().UTC()}
}

