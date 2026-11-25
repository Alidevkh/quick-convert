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
	Input  string         // full path of the original file (the first one, if several)
	Inputs []string       // all original files, in order (just Input for one file)
	Base   string         // original name without extension, e.g. "trip"
	Ext    string         // original extension, lower case, e.g. ".mov"
	Action actions.Action // what to do
	Work   *outpath.Work  // temporary folder; inputs are inside as Work.In / Work.Ins
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
	Version    string                                     // program version, written in the log
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
	defer func() { writeLog(d.Log, d.Version, actionID, input, outputs, runErr) }()

	fail := func(err error) int {
		runErr = err
		d.Notify("Couldn't convert "+name, reason(err, ext))
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
	info, err := os.Stat(input)
	if err != nil || !info.Mode().IsRegular() {
		return fail(&UserError{Msg: "File not found"})
	}
	if info.Size() == 0 {
		return fail(&UserError{Msg: "This file is empty"})
	}

	result, err := execute(ctx, action, []string{input}, d)
	outputs, runErr = result.Outputs, err
	return finish(d, result, err, fail, filepath.Dir(input))
}

// RunMany carries out a Send to action (several files in one run, e.g.
// "Images to PDF"). Files the action can't use are skipped, and the rest
// are put in natural file-name order.
func RunMany(ctx context.Context, actionID string, inputs []string, d Deps) (code int) {
	var outputs []string
	var runErr error
	defer func() { writeLog(d.Log, d.Version, actionID, strings.Join(inputs, ", "), outputs, runErr) }()

	entry, ok := actions.SendToFor(actionID)
	if !ok {
		runErr = fmt.Errorf("unknown Send to action %q", actionID)
		d.Notify("QuickConvert", runErr.Error())
		return ExitUsage
	}
	action, _ := actions.Lookup(actionID)

	usable := usableFiles(inputs, entry.Exts)
	if len(usable) == 0 {
		runErr = &UserError{Msg: entry.NoneMsg}
		d.Notify("QuickConvert", entry.NoneMsg)
		return ExitFailed
	}
	SortNatural(usable)

	name := filepath.Base(usable[0])
	fail := func(err error) int {
		runErr = err
		d.Notify("Couldn't make "+action.Label, reason(err, filepath.Ext(name)))
		return ExitFailed
	}
	result, err := execute(ctx, action, usable, d)
	outputs, runErr = result.Outputs, err
	return finish(d, result, err, fail, filepath.Dir(usable[0]))
}

// usableFiles keeps the regular, non-empty files with one of the extensions.
func usableFiles(paths, exts []string) []string {
	var keep []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			continue
		}
		ext := strings.ToLower(filepath.Ext(p))
		for _, e := range exts {
			if ext == e {
				keep = append(keep, p)
				break
			}
		}
	}
	return keep
}

// execute is the part shared by Run and RunMany: show "Converting…" if it
// takes a while, wait for a heavy-job slot, make the temporary folder (next
// to the first file), put the inputs in it, and run the converter.
func execute(ctx context.Context, action actions.Action, inputs []string, d Deps) (Result, error) {
	first := inputs[0]
	name := filepath.Base(first)

	// Tell the user we started, but only if the job takes a while — and
	// not for actions that first ask a question (the user is still typing).
	startTimer := time.AfterFunc(d.StartDelay, func() {
		d.Notify(fmt.Sprintf("Converting %s → %s…", name, action.Label), "")
	})
	if action.Asks {
		startTimer.Stop()
	}
	defer startTimer.Stop()

	if action.Kind == actions.Heavy {
		release := d.Acquire()
		defer release()
	}

	work, err := outpath.NewWork(filepath.Dir(first))
	if err != nil {
		return Result{}, err
	}
	defer work.Cleanup()
	if len(inputs) == 1 {
		err = work.LinkInput(first)
	} else {
		err = work.LinkInputs(inputs)
	}
	if err != nil {
		return Result{}, &UserError{Msg: "Couldn't read the file"}
	}

	ext := strings.ToLower(filepath.Ext(name))
	job := Job{
		Input:  first,
		Inputs: inputs,
		Base:   strings.TrimSuffix(name, filepath.Ext(name)),
		Ext:    ext,
		Action: action,
		Work:   work,
		Eng:    d.Eng,
	}
	return d.Handlers(action, ext)(ctx, job)
}

// finish shows the final notification and returns the exit code.
func finish(d Deps, result Result, err error, fail func(error) int, folder string) int {
	var notice *Info
	switch {
	case errors.As(err, &notice):
		if notice.Msg != "" { // an empty Info means "cancelled": say nothing
			d.Notify("QuickConvert", notice.Msg)
		}
		return ExitOK
	case err != nil:
		return fail(err)
	}
	d.Notify(doneTitle(result), "Saved in "+folder)
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
// Engine messages are technical, so the user sees a plain sentence instead;
// the engine's own message still goes into the log file.
func reason(err error, ext string) string {
	var userErr *UserError
	var engineErr *engines.EngineError
	switch {
	case errors.As(err, &userErr):
		return userErr.Msg
	case errors.As(err, &engineErr):
		return fmt.Sprintf("This file looks damaged or isn't a real %s file", strings.ToUpper(strings.TrimPrefix(ext, ".")))
	default:
		return err.Error() // includes MissingError and ErrNotWritable
	}
}

// writeLog appends one line: time, version, action, input, then the results
// or the error. For an engine failure it also adds the engine's last few
// messages, which usually say exactly why it failed.
func writeLog(w io.Writer, version, actionID, input string, outputs []string, err error) {
	if w == nil {
		return
	}
	line := fmt.Sprintf("%s v%s %s %s -> %s", time.Now().Format(time.RFC3339), version, actionID, input, strings.Join(outputs, ", "))
	var notice *Info
	if errors.As(err, &notice) && notice.Msg == "" {
		line += " | cancelled"
	} else if err != nil {
		line += " | " + err.Error()
		var engineErr *engines.EngineError
		if errors.As(err, &engineErr) {
			line += " | details: " + lastLines(engineErr.Stderr, 5)
		}
	}
	fmt.Fprintln(w, line)
}

// lastLines returns the last n non-empty lines of s, joined with " / ".
func lastLines(s string, n int) string {
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, " / ")
}
