// Package uptime provides uptime-related health Details.
package uptime

import (
	"github.com/nelkinda/health-go"
	"net/http"
	"syscall"
	"time"
)

type system struct {
}

func (u *system) HealthDetails() map[string][]health.Details {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var uptime func() health.Details
	if err != nil {
		uptime = func() health.Details {
			return health.Details{
				ComponentType: "system",
				Status:        health.Fail,
				Output:        err.Error(),
				Time:          now,
			}
		}
	} else {
		uptime = func() health.Details {
			return health.Details{
				ComponentType: "system",
				ObservedValue: si.Uptime,
				ObservedUnit:  "s",
				Status:        health.Pass,
				Time:          now,
			}
		}
	}
	return map[string][]health.Details{
		"uptime": {
			uptime(),
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
