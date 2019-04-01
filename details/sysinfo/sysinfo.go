// Package sysinfo provides sysinfo as health details.

// +build !linux

package sysinfo

type sysinfo struct {
}

func (u *sysinfo) HealthDetails() map[string][]health.Details {
	return map[string][]health.Details{}
}

func (*sysinfo) AuthorizeHealth(r *http.Request) bool {
	return true
}

// SysInfo returns a DetailsProvider that provides sysinfo statistics.
func Health() health.DetailsProvider {
	return &sysinfo{}
}
