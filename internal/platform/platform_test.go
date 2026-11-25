package platform

import (
	"strings"
	"testing"
)

// The notification text must only come from environment variables, never
// be pasted into the script: a file name like  it's "mine".mov  must not be
// able to break (or inject into) the PowerShell script.
func TestToastScriptUsesEnvOnly(t *testing.T) {
	for _, want := range []string{"$env:QC_TITLE", "$env:QC_BODY"} {
		if !strings.Contains(ToastScript, want) {
			t.Errorf("script does not use %s", want)
		}
	}
	for _, bad := range []string{"%s", "{0}", "%v"} {
		if strings.Contains(ToastScript, bad) {
			t.Errorf("script contains placeholder %q", bad)
		}
	}
}

func TestEncodeCommandIsUTF16Base64(t *testing.T) {
	// "A" in UTF-16LE is 0x41 0x00, which is "QQA=" in base64.
	if got := encodeCommand("A"); got != "QQA=" {
		t.Errorf("got %q", got)
	}
}
