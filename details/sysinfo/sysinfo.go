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

// Health returns a DetailsProvider that provides sysinfo statistics.
// On Linux, this will be details from syscall.Sysinfo_t.
// On other platforms, this provider provides no information.
func Health() health.DetailsProvider {
	return &sysinfo{}
}
