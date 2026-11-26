package engines

// Hides the black console window when an engine runs on Windows.

import (
	"os/exec"
	"syscall"
)

// createNoWindow is the Windows flag that starts a console program
// without showing a black console window.
const createNoWindow = 0x08000000

// hide makes sure the engine runs invisibly.
func hide(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
}
