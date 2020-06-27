// Package uptime provides uptime-related health Checks.
package uptime

import (
	"github.com/nelkinda/health-go"
	"net/http"
	"syscall"
	"time"
)

type system struct {
}

func (u *system) HealthChecks() map[string][]health.Checks {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var uptime func() health.Checks
	if err != nil {
		uptime = func() health.Checks {
			return health.Checks{
				ComponentType: "system",
				Status:        health.Fail,
				Output:        err.Error(),
				Time:          now,
			}
		}
	} else {
		uptime = func() health.Checks {
			return health.Checks{
				ComponentType: "system",
				ObservedValue: si.Uptime,
				ObservedUnit:  "s",
				Status:        health.Pass,
				Time:          now,
			}
		}
	}
	return map[string][]health.Checks{
		"uptime": {
			uptime(),
		},
	}
}

func (*system) AuthorizeHealth(*http.Request) bool {
	return true
}

// System returns a ChecksProvider for health checks about the system uptime.
func System() health.ChecksProvider {
	return &system{}
}
