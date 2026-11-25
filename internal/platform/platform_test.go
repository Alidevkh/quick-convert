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

// The "Which pages?" box: text comes from environment variables (never
// pasted into the script), and the answer is printed as UTF-8 so Arabic
// digits survive.
func TestAskScript(t *testing.T) {
	for _, want := range []string{"$env:QC_PROMPT", "$env:QC_TITLE", "InputBox", "UTF8"} {
		if !strings.Contains(AskScript, want) {
			t.Errorf("AskScript does not contain %s", want)
		}
	}
	if got := AskPrompt(12); got != "Which pages? (1-12)\n\nExamples:  1-3, 5   or   5, 1-3" {
		t.Errorf("prompt = %q", got)
	}
}
