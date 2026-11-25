package platform

// The small "Which pages?" box (Windows).

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// askTimeout: if nobody answers for this long, give up (so a box that got
// lost behind other windows can't keep QuickConvert waiting forever).
const askTimeout = 15 * time.Minute

// AskPages shows the "Which pages?" box for a PDF with `pages` pages and
// waits for the answer. ok is false if the user cancelled. err is set if
// the box could not be shown at all (e.g. PowerShell blocked by policy).
func AskPages(pages int) (answer string, ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), askTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell.exe",
		"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-EncodedCommand", encodeCommand(AskScript),
	)
	cmd.Env = append(os.Environ(), "QC_PROMPT="+AskPrompt(pages), "QC_TITLE=QuickConvert - Extract pages")
	hide, flags := askProcessFlags()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hide, CreationFlags: flags}
	out, err := cmd.Output()
	if err != nil {
		return "", false, err
	}
	answer = string(out)
	return answer, answer != "", nil
}
