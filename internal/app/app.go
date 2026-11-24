// Package app runs one right-click request from start to finish:
//
//	check the request → (wait for a free slot if heavy) → make a temporary
//	folder → run the converter → save the result → notify → write the log
//
// The actual converters live in their own packages (media, pdf, photo...)
// and are plugged in as Handlers, so this file contains no conversion
// details at all.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"quickconvert/internal/actions"
	"quickconvert/internal/engines"
	"quickconvert/internal/outpath"
)

// Job is everything a converter needs to know about one request.
type Job struct {
	Input  string         // full path of the original file
	Base   string         // original name without extension, e.g. "trip"
	Ext    string         // original extension, lower case, e.g. ".mov"
	Action actions.Action // what to do
	Work   *outpath.Work  // temporary folder; the input is inside as Work.In
	Eng    *engines.Set   // the converter programs
}

// Result is what a converter returns when it worked.
type Result struct {
	Outputs []string // full paths of the saved files
	Detail  string   // optional extra for the notification, e.g. "4.2 MB → 0.9 MB"
}

// Handler is a converter: it does the work for one Job.
type Handler func(ctx context.Context, j Job) (Result, error)

// Deps are the outside pieces Run uses. Tests replace them with fakes.
type Deps struct {
	Notify     func(title, body string)                   // show a Windows notification
	Acquire    func() (release func())                    // wait for a heavy-job slot
	Log        io.Writer                                  // the log file
	Eng        *engines.Set                               // where the engines are
	Handlers   func(a actions.Action, ext string) Handler // pick the converter
	StartDelay time.Duration                              // show "Converting…" only after this long
}

// Exit codes returned by Run.
const (
	ExitOK     = 0
	ExitFailed = 1
	ExitUsage  = 2
	maxListed  = 3 // file names listed in the "Done" notification
)

// Run carries out action actionID on the file input and returns the
// program's exit code.
func Run(ctx context.Context, actionID, input string, d Deps) (code int) {
	name := filepath.Base(input)
	ext := strings.ToLower(filepath.Ext(name))
	var outputs []string
	var runErr error

	// Whatever happens, write one line to the log at the end.
	defer func() { writeLog(d.Log, actionID, input, outputs, runErr) }()

	fail := func(err error) int {
		runErr = err
		d.Notify("Couldn't convert "+name, reason(err))
		return ExitFailed
	}

	action, ok := actions.Lookup(actionID)
	if !ok {
		runErr = fmt.Errorf("unknown action %q", actionID)
		d.Notify("QuickConvert", runErr.Error())
		return ExitUsage
	}
	if !actions.Supports(actionID, ext) {
		return fail(&UserError{Msg: fmt.Sprintf("%s isn't available for %s files", action.Label, ext)})
	}
	if info, err := os.Stat(input); err != nil || !info.Mode().IsRegular() {
		return fail(&UserError{Msg: "File not found"})
	}

	// Tell the user we started, but only if the job takes a while.
	startTimer := time.AfterFunc(d.StartDelay, func() {
		d.Notify(fmt.Sprintf("Converting %s → %s…", name, action.Label), "")
	})
	defer startTimer.Stop()

	if action.Kind == actions.Heavy {
		release := d.Acquire()
		defer release()
	}

	work, err := outpath.NewWork(filepath.Dir(input))
	if err != nil {
		return fail(err)
	}
	defer work.Cleanup()
	if err := work.LinkInput(input); err != nil {
		return fail(&UserError{Msg: "Couldn't read the file"})
	}

	job := Job{
		Input:  input,
		Base:   strings.TrimSuffix(name, filepath.Ext(name)),
		Ext:    ext,
		Action: action,
		Work:   work,
		Eng:    d.Eng,
	}
	result, err := d.Handlers(action, ext)(ctx, job)
	startTimer.Stop()
	outputs = result.Outputs

	var info *Info
	switch {
	case errors.As(err, &info):
		runErr = err
		d.Notify("QuickConvert", info.Msg)
		return ExitOK
	case err != nil:
		return fail(err)
	}

	d.Notify(doneTitle(result), "Saved in "+filepath.Dir(input))
	return ExitOK
}

// doneTitle builds "Done: a.png, b.png" or "Done: a.png + 4 more",
// followed by the optional detail in brackets.
func doneTitle(r Result) string {
	var names []string
	for _, p := range r.Outputs {
		names = append(names, filepath.Base(p))
	}
	title := "Done: " + strings.Join(names, ", ")
	if len(names) > maxListed {
		title = fmt.Sprintf("Done: %s + %d more", names[0], len(names)-1)
	}
	if r.Detail != "" {
		title += " (" + r.Detail + ")"
	}
	return title
}

// reason turns an error into the short text shown in the notification.
func reason(err error) string {
	var userErr *UserError
	var engineErr *engines.EngineError
	switch {
	case errors.As(err, &userErr):
		return userErr.Msg
	case errors.As(err, &engineErr):
		return engineErr.Tail
	default:
		return err.Error() // includes MissingError and ErrNotWritable
	}
}

// writeLog appends one line: time, action, input, then results or error.
func writeLog(w io.Writer, actionID, input string, outputs []string, err error) {
	if w == nil {
		return
	}
	line := fmt.Sprintf("%s %s %s -> %s", time.Now().Format(time.RFC3339), actionID, input, strings.Join(outputs, ", "))
	if err != nil {
		line += " | " + err.Error()
	}
	fmt.Fprintln(w, line)
}
