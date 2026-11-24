package platform

import (
	"os"
	"os/exec"
	"syscall"
)

// Notify shows a Windows notification. It starts a hidden PowerShell and
// does not wait for it, so QuickConvert can exit right away.
func Notify(title, body string) {
	cmd := exec.Command("powershell.exe",
		"-NoProfile",                 // skip the user's PowerShell profile (faster)
		"-NonInteractive",            // never wait for input
		"-ExecutionPolicy", "Bypass", // allow our one script even if scripts are blocked
		"-EncodedCommand", encodeCommand(ToastScript),
	)
	cmd.Env = append(os.Environ(), "QC_TITLE="+title, "QC_BODY="+body)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	_ = cmd.Start()
}
