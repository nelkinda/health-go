// Package uptime provides uptime-related health Details.
package uptime

import (
	"github.com/capnm/sysinfo"
	"github.com/nelkinda/health-go"
	"net/http"
	"time"
)

type system struct {
}

func (u *system) HealthDetails() map[string][]health.Details {
	si := sysinfo.Get()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return map[string][]health.Details{
		"uptime": {
			{
				ComponentType: "system",
				ObservedValue: si.Uptime.Seconds(),
				ObservedUnit: "s",
				Status:health.Pass,
				Time: now,
			},
		},
	}
}

func (*system) AuthorizeHealth(r *http.Request) bool {
	return true
}

// System returns a DetailsProvider for health details about the system uptime.
func System() health.DetailsProvider {
	return &system{}
}

