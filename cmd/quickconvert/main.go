// QuickConvert is started by the Windows right-click menu like this:
//
//	quickconvert.exe <action> "<file>"      e.g.  quickconvert.exe mp4 "C:\Videos\trip.mov"
//
// or, from the Send to menu, with every selected file:
//
//	quickconvert.exe imgpdf "a.jpg" "b.jpg" ...
//
// It does the job, shows a notification, and exits. Nothing keeps running
// in the background. See internal/app for the step-by-step flow.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/engines"
	"quickconvert/internal/pdf"
	"quickconvert/internal/platform"
)

// version is set when building (see scripts/build.sh).
var version = "dev"

// maxLogSize: when the log file grows past this, it is started fresh.
const maxLogSize = 1 << 20 // 1 MB

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println("QuickConvert", version)
		return
	}
	_, isSendTo := actions.SendToFor(argAt(1))
	if len(os.Args) < 3 || (!isSendTo && len(os.Args) != 3) {
		platform.Notify("QuickConvert", "Usage: quickconvert <action> <file> (Send to actions: <file> <file>...)")
		os.Exit(app.ExitUsage)
	}

	pdf.AskPages = platform.AskPages // the "Which pages?" box for Extract pages...

	logFile := openLog()
	defer logFile.Close()

	deps := app.Deps{
		Notify:     platform.Notify,
		Acquire:    platform.AcquireHeavy,
		Log:        logFile,
		Eng:        engines.Default(),
		Handlers:   handlerFor,
		StartDelay: time.Second,
		Version:    version,
	}
	var code int
	if isSendTo {
		// "Send to" passes every selected file: quickconvert.exe imgpdf a.jpg b.jpg ...
		code = app.RunMany(context.Background(), os.Args[1], os.Args[2:], deps)
	} else {
		code = app.Run(context.Background(), os.Args[1], os.Args[2], deps)
	}
	logFile.Close()
	os.Exit(code)
}

// openLog opens %LOCALAPPDATA%\QuickConvert\log.txt for adding lines.
// If the log is bigger than 1 MB it is emptied first. If it cannot be
// opened at all, logging is skipped (it must never stop a conversion).
func openLog() io.WriteCloser {
	cache, err := os.UserCacheDir() // %LOCALAPPDATA% on Windows
	if err != nil {
		return nopCloser{io.Discard}
	}
	dir := filepath.Join(cache, "QuickConvert")
	os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, "log.txt")

	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if info, err := os.Stat(path); err == nil && info.Size() > maxLogSize {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		return nopCloser{io.Discard}
	}
	return f
}

// argAt returns command-line argument i, or "" if there is none.
func argAt(i int) string {
	if i < len(os.Args) {
		return os.Args[i]
	}
	return ""
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }
