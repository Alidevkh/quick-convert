package platform

// The small "Which pages?" input box (Windows).

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// AskPages shows the "Which pages?" box for a PDF with `pages` pages and
// waits for the answer. ok is false if the user cancelled.
func AskPages(pages int) (answer string, ok bool) {
	cmd := exec.Command("powershell.exe",
		"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", encodeCommand(AskScript),
	)
	cmd.Env = append(os.Environ(), "QC_PROMPT="+AskPrompt(pages), "QC_TITLE=QuickConvert - Extract pages")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // no console window
	out, err := cmd.Output()
	answer = strings.TrimSpace(string(out))
	return answer, err == nil && answer != ""
}
