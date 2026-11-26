//go:build !windows

package engines

// Stand-in for hiding the console window when running outside Windows.

import "os/exec"

// hide does nothing outside Windows (there is no console window to hide).
func hide(cmd *exec.Cmd) {}
