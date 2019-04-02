// Package sysinfo provides sysinfo as health details.

// +build linux

package sysinfo

import (
	"fmt"
	"github.com/nelkinda/health-go"
	"net/http"
	"os"
	"syscall"
	"time"
)

type sysinfo struct {
}

func (u *sysinfo) HealthDetails() map[string][]health.Details {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var uptime func() health.Details
	var processes func() health.Details
	var cpuutil func(componentId string, load uint64) health.Details
	var memutil func(componentId string, load uint64) health.Details
	var hostname func() health.Details
	if err != nil {
		cpuutil = func(componentId string, load uint64) health.Details {
			return health.Details{
				ComponentType: "system",
				ComponentID:   componentId,
				Status:        health.Fail,
				Output:        err.Error(),
				Time:          now,
			}
		}
		memutil = cpuutil
		uptime = func() health.Details {
			return health.Details{
				ComponentType: "system",
				Status:        health.Fail,
				Output:        err.Error(),
				Time:          now,
			}
		}
		processes = uptime
	} else {
		memunit := fmt.Sprintf("%d bytes", si.Unit)
		cpuutil = func(componentId string, load uint64) health.Details {
			return health.Details{
				ComponentType: "system",
				ComponentID:   componentId,
				ObservedValue: load / 65536.0,
				ObservedUnit:  "%",
				Status:        health.Pass,
				Time:          now,
			}
		}
		memutil = func(componentId string, memory uint64) health.Details {
			return health.Details{
				ComponentType: "system",
				ComponentID:   componentId,
				ObservedValue: memory,
				ObservedUnit:  memunit,
				Status:        health.Pass,
				Time:          now,
			}
		}
		uptime = func() health.Details {
			return health.Details{
				ComponentType: "system",
				ObservedValue: si.Uptime,
				ObservedUnit:  "s",
				Status:        health.Pass,
				Time:          now,
			}
		}
		processes = func() health.Details {
			return health.Details{
				ComponentID:   "Processes",
				ComponentType: "system",
				ObservedValue: si.Procs,
				Status:        health.Pass,
				Time:          now,
			}
		}
	}

	if hn, err := os.Hostname(); err == nil {
		hostname = func() health.Details {
			return health.Details{
				ComponentID:   "hostname",
				ComponentType: "system",
				ObservedValue: hn,
				Status:        health.Pass,
				Time:          now,
			}
		}
	} else {
		hostname = func() health.Details {
			return health.Details{
				ComponentID:   "hostname",
				ComponentType: "system",
				Status:        health.Fail,
				Time:          now,
				Output:        err.Error(),
			}
		}
	}

	return map[string][]health.Details{
		"uptime": {
			uptime(),
		},
		"hostname": {
			hostname(),
		},
		"cpu:utilization": {
			cpuutil("1 minute", si.Loads[0]),
			cpuutil("5 minutes", si.Loads[1]),
			cpuutil("15 minutes", si.Loads[2]),
			processes(),
		},
		"memory:utilization": {
			memutil("Total Ram", si.Totalram),
			memutil("Free Ram", si.Freeram),
			memutil("Shared Ram", si.Sharedram),
			memutil("Buffer Ram", si.Bufferram),
			memutil("Total Swap", si.Totalswap),
			memutil("Free Swap", si.Freeswap),
			memutil("Total High", si.Totalhigh),
			memutil("Free High", si.Freehigh),
		},
	}
}

func (*sysinfo) AuthorizeHealth(r *http.Request) bool {
	return true
}

// Health returns a DetailsProvider that provides sysinfo statistics.
// On Linux, this will be details from syscall.Sysinfo_t.
// On other platforms, this provider provides no information.
func Health() health.DetailsProvider {
	return &sysinfo{}
}
