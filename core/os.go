package core

import "runtime"

var OsLsArgs []string

// OS specific stuff, hopefully very little
func init() {
	switch runtime.GOOS {
	case "darwin":
		OsLsArgs = []string{"-a1", "-G"}
	default:
		OsLsArgs = []string{"-a1", "--color=always"}
	}
}
