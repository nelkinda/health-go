// Package sysinfo provides sysinfo as health checks.

// +build !linux

package sysinfo

import (
	"github.com/nelkinda/health-go"
	"net/http"
)

type sysinfo struct {
}

func (u *sysinfo) HealthChecks(ctx context.Context) map[string][]health.Checks {
	return map[string][]health.Checks{}
}

func (*sysinfo) AuthorizeHealth(r *http.Request) bool {
	return true
}

// Health returns a ChecksProvider that provides sysinfo statistics.
// On Linux, this will be checks from syscall.Sysinfo_t.
// On other platforms, this provider provides no information.
func Health() health.ChecksProvider {
	return &sysinfo{}
}
