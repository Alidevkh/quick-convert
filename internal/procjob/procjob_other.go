//go:build !windows

package procjob

import "os/exec"

// Track does nothing outside Windows.
func Track(cmd *exec.Cmd) {}
