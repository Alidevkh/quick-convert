package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"quickconvert/internal/actions"
	"quickconvert/internal/engines"
)

// recorder remembers every notification shown.
type recorder struct {
	mu    sync.Mutex
	calls [][2]string
}

func (r *recorder) notify(title, body string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, [2]string{title, body})
}

func (r *recorder) all() [][2]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([][2]string{}, r.calls...)
}

func (r *recorder) last() [2]string {
	c := r.all()
	if len(c) == 0 {
		return [2]string{}
	}
	return c[len(c)-1]
}

// setup creates an input file and Deps whose handler is h.
func setup(t *testing.T, name string, h Handler) (string, *recorder, *bytes.Buffer, Deps) {
	t.Helper()
	dir := t.TempDir()
	input := filepath.Join(dir, name)
	os.WriteFile(input, []byte("original"), 0o644)
	rec := &recorder{}
	log := &bytes.Buffer{}
	d := Deps{
		Notify:     rec.notify,
		Acquire:    func(actions.Kind) func() { return func() {} },
		Log:        log,
		Handlers:   func(actions.Action, string) Handler { return h },
		StartDelay: time.Hour,
	}
	return input, rec, log, d
}

// commitCopy is a fake handler that saves the input as Base+Tag+OutExt.
func commitCopy(ctx context.Context, j Job) (Result, error) {
	data, _ := os.ReadFile(filepath.Join(j.Work.Dir, j.Work.In))
	os.WriteFile(filepath.Join(j.Work.Dir, "out"), data, 0o644)
	p, err := j.Work.Commit("out", j.Base+j.Action.Tag, j.Action.OutExt)
	return Result{Outputs: []string{p}}, err
}

func TestRunSuccessNotifies(t *testing.T) {
	input, rec, _, d := setup(t, "x.txt", commitCopy)
	if code := Run(context.Background(), "code-py", input, d); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got[0] != "Done: x.py" {
		t.Errorf("notify = %v", got)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(input), "x.py")); err != nil {
		t.Error("x.py not created")
	}
}

func TestRunFailureCleansUp(t *testing.T) {
	failing := func(ctx context.Context, j Job) (Result, error) {
		os.WriteFile(filepath.Join(j.Work.Dir, "half.py"), []byte("x"), 0o644)
		return Result{}, &UserError{Msg: "boom"}
	}
	input, rec, _, d := setup(t, "x.txt", failing)
	before, _ := os.Stat(input)

	if code := Run(context.Background(), "code-py", input, d); code != 1 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got != [2]string{"Couldn't convert x.txt", "boom"} {
		t.Errorf("notify = %v", got)
	}
	left, _ := filepath.Glob(filepath.Join(filepath.Dir(input), ".qc-*.part"))
	if len(left) != 0 {
		t.Errorf("left behind: %v", left)
	}
	after, _ := os.Stat(input)
	data, _ := os.ReadFile(input)
	if string(data) != "original" || !after.ModTime().Equal(before.ModTime()) {
		t.Error("original was changed")
	}
}

func TestRunUnsupported(t *testing.T) {
	input, rec, _, d := setup(t, "x.txt", commitCopy)
	if code := Run(context.Background(), "mp4", input, d); code != 1 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got[1] != "MP4 isn't available for .txt files" {
		t.Errorf("notify = %v", got)
	}
}

func TestRunMissingInput(t *testing.T) {
	_, rec, _, d := setup(t, "x.txt", commitCopy)
	if code := Run(context.Background(), "code-py", "/nope/x.txt", d); code != 1 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got[0] != "Couldn't convert x.txt" {
		t.Errorf("notify = %v", got)
	}
}

func TestRunDetail(t *testing.T) {
	h := func(ctx context.Context, j Job) (Result, error) {
		r, err := commitCopy(ctx, j)
		r.Detail = "2 MB → 1 MB"
		return r, err
	}
	input, rec, _, d := setup(t, "x.pdf", h)
	Run(context.Background(), "shrink", input, d)
	if got := rec.last(); got[0] != "Done: x (small).pdf (2 MB → 1 MB)" {
		t.Errorf("notify = %v", got)
	}
}

func TestRunManyOutputs(t *testing.T) {
	h := func(ctx context.Context, j Job) (Result, error) {
		return Result{Outputs: []string{"/a/r-1.png", "/a/r-2.png", "/a/r-3.png", "/a/r-4.png", "/a/r-5.png"}}, nil
	}
	input, rec, _, d := setup(t, "r.pdf", h)
	Run(context.Background(), "png", input, d)
	if got := rec.last(); got[0] != "Done: r-1.png + 4 more" {
		t.Errorf("notify = %v", got)
	}
}

func TestRunInfo(t *testing.T) {
	h := func(ctx context.Context, j Job) (Result, error) {
		return Result{}, &Info{Msg: "Already optimized — no smaller file made"}
	}
	input, rec, _, d := setup(t, "x.pdf", h)
	if code := Run(context.Background(), "shrink", input, d); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got != [2]string{"QuickConvert", "Already optimized — no smaller file made"} {
		t.Errorf("notify = %v", got)
	}
}

func TestStartNotificationOnlyWhenSlow(t *testing.T) {
	slow := func(ctx context.Context, j Job) (Result, error) {
		time.Sleep(120 * time.Millisecond)
		return commitCopy(ctx, j)
	}
	input, rec, _, d := setup(t, "x.mov", slow)
	d.StartDelay = 50 * time.Millisecond
	Run(context.Background(), "mp4", input, d)
	calls := rec.all()
	if len(calls) != 2 || calls[0][0] != "Converting x.mov → MP4…" {
		t.Errorf("slow job calls = %v", calls)
	}

	// A fast job, with the program's real 1-second delay. (A 50 ms delay
	// made this test flaky: on a busy machine even creating the temporary
	// folder can take longer than 50 ms.)
	input2, rec2, _, d2 := setup(t, "y.mov", commitCopy)
	d2.StartDelay = time.Second
	Run(context.Background(), "mp4", input2, d2)
	time.Sleep(1100 * time.Millisecond) // a timer that was not stopped would fire here
	if calls := rec2.all(); len(calls) != 1 {
		t.Errorf("fast job calls = %v", calls)
	}
}

// Every job waits for a slot of its kind (2 heavy, 4 light at once), so
// selecting 30 files doesn't start 30 conversions together. A job that
// first asks the user a question doesn't take a slot while waiting.
func TestJobsAcquireSlotByKind(t *testing.T) {
	var got []actions.Kind
	released := 0
	acquire := func(k actions.Kind) func() {
		got = append(got, k)
		return func() { released++ }
	}
	input, _, _, d := setup(t, "x.mov", commitCopy)
	d.Acquire = acquire
	Run(context.Background(), "mp4", input, d)
	input2, _, _, d2 := setup(t, "x.txt", commitCopy)
	d2.Acquire = acquire
	Run(context.Background(), "code-py", input2, d2)
	if !reflect.DeepEqual(got, []actions.Kind{actions.Heavy, actions.Light}) || released != 2 {
		t.Errorf("acquired %v, released %d", got, released)
	}

	got = nil
	input3, _, _, d3 := setup(t, "x.pdf", commitCopy)
	d3.Acquire = acquire
	Run(context.Background(), "extract", input3, d3)
	if len(got) != 0 {
		t.Errorf("an action that asks first took a slot: %v", got)
	}
}

func TestLogLine(t *testing.T) {
	input, _, log, d := setup(t, "x.txt", commitCopy)
	Run(context.Background(), "code-py", input, d)
	line := log.String()
	if strings.Count(line, "\n") != 1 || !strings.Contains(line, "code-py") ||
		!strings.Contains(line, input) || !strings.Contains(line, "x.py") {
		t.Errorf("log = %q", line)
	}
}

func TestRunEngineFailureIsFriendly(t *testing.T) {
	broken := func(ctx context.Context, j Job) (Result, error) {
		return Result{}, &engines.EngineError{Name: "ffmpeg", Tail: "Invalid data found when processing input"}
	}
	input, rec, log, d := setup(t, "clip.mov", broken)
	Run(context.Background(), "mp4", input, d)
	if got := rec.last(); got[1] != "This file looks damaged or isn't a real MOV file" {
		t.Errorf("notify = %v", got)
	}
	if !strings.Contains(log.String(), "Invalid data found") {
		t.Error("the engine's own message should still be in the log")
	}
}

func TestRunEmptyFile(t *testing.T) {
	called := false
	h := func(ctx context.Context, j Job) (Result, error) { called = true; return Result{}, nil }
	input, rec, _, d := setup(t, "doc.pdf", h)
	os.WriteFile(input, nil, 0o644)
	if code := Run(context.Background(), "shrink", input, d); code != 1 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got[1] != "This file is empty" || called {
		t.Errorf("notify = %v, handler called = %v", got, called)
	}
}

func TestSortNatural(t *testing.T) {
	paths := []string{"b/IMG_10.jpg", "a/IMG_9.jpg", "IMG_2.JPG", "img_1.jpg"}
	SortNatural(paths)
	want := []string{"img_1.jpg", "IMG_2.JPG", "a/IMG_9.jpg", "b/IMG_10.jpg"}
	if fmt.Sprint(paths) != fmt.Sprint(want) {
		t.Errorf("got %v", paths)
	}
}

func TestRunManyFiltersAndSorts(t *testing.T) {
	var seen Job
	h := func(ctx context.Context, j Job) (Result, error) { seen = j; return Result{}, nil }
	input, _, _, d := setup(t, "IMG_10.jpg", h)
	dir := filepath.Dir(input)
	write := func(name, content string) string {
		p := filepath.Join(dir, name)
		os.WriteFile(p, []byte(content), 0o644)
		return p
	}
	notes := write("notes.txt", "x")
	img9 := write("IMG_9.jpg", "x")
	empty := write("IMG_1.jpg", "")
	folder := filepath.Join(dir, "sub.jpg")
	os.Mkdir(folder, 0o755)

	if code := RunMany(context.Background(), "imgpdf", []string{input, notes, img9, folder, empty}, d); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if fmt.Sprint(seen.Inputs) != fmt.Sprint([]string{img9, input}) || seen.Base != "IMG_9" {
		t.Errorf("Inputs = %v Base = %q", seen.Inputs, seen.Base)
	}
	if len(seen.Work.Ins) != 2 {
		t.Errorf("Ins = %v", seen.Work.Ins)
	}
}

func TestRunManyNoneAccepted(t *testing.T) {
	input, rec, _, d := setup(t, "notes.txt", commitCopy)
	if code := RunMany(context.Background(), "imgpdf", []string{input}, d); code != 1 {
		t.Fatalf("exit %d", code)
	}
	if got := rec.last(); got[1] != "No pictures selected" {
		t.Errorf("notify = %v", got)
	}
}

func TestRunManyUnknownAction(t *testing.T) {
	input, _, _, d := setup(t, "a.mov", commitCopy)
	if code := RunMany(context.Background(), "mp4", []string{input}, d); code != 2 {
		t.Errorf("exit %d", code)
	}
}

func TestRunSetsInputs(t *testing.T) {
	var seen Job
	h := func(ctx context.Context, j Job) (Result, error) { seen = j; return Result{}, nil }
	input, _, _, d := setup(t, "x.txt", h)
	Run(context.Background(), "code-py", input, d)
	if fmt.Sprint(seen.Inputs) != fmt.Sprint([]string{input}) || len(seen.Work.Ins) != 1 {
		t.Errorf("Inputs = %v Ins = %v", seen.Inputs, seen.Work.Ins)
	}
}

// The log must say which version ran and what the engine really said, so a
// failure on the owner's PC can be understood from one log line.
func TestLogIncludesVersionAndEngineDetails(t *testing.T) {
	broken := func(ctx context.Context, j Job) (Result, error) {
		return Result{}, &engines.EngineError{Name: "ffmpeg", Tail: "Error opening output files: Invalid argument",
			Stderr: "Input #0, mov\nFiltergraph 'x' was specified for a stream fed from a complex filtergraph.\nError opening output file out.png.\nError opening output files: Invalid argument\n"}
	}
	input, _, log, d := setup(t, "a.heic", broken)
	d.Version = "0.3.0"
	Run(context.Background(), "png", input, d)
	line := log.String()
	if !strings.Contains(line, "v0.3.0") {
		t.Errorf("version missing: %q", line)
	}
	if !strings.Contains(line, "complex filtergraph") {
		t.Errorf("engine details missing: %q", line)
	}
	if strings.Count(line, "\n") != 1 {
		t.Errorf("log entry must stay on one line: %q", line)
	}
}

func TestLogKeepsEngineDetailsThroughForFile(t *testing.T) {
	broken := func(ctx context.Context, j Job) (Result, error) {
		return Result{}, ForFile("a.heic", &engines.EngineError{Name: "ffmpeg", Tail: "Invalid argument",
			Stderr: "complex filtergraph problem\nInvalid argument\n"})
	}
	input, rec, log, d := setup(t, "a.heic", broken)
	RunMany(context.Background(), "imgpdf", []string{input}, d)
	if !strings.Contains(log.String(), "complex filtergraph problem") {
		t.Errorf("engine details lost: %q", log.String())
	}
	if got := rec.last(); !strings.HasPrefix(got[1], "a.heic: ") {
		t.Errorf("notification = %v", got)
	}
}

// Actions that ask the user something (Extract pages...) must not show
// "Converting…" while the question box is open.
func TestNoStartNotificationForInteractiveAction(t *testing.T) {
	slow := func(ctx context.Context, j Job) (Result, error) {
		time.Sleep(120 * time.Millisecond) // the user is typing
		return Result{}, &Info{}
	}
	input, rec, log, d := setup(t, "doc.pdf", slow)
	d.StartDelay = 50 * time.Millisecond
	Run(context.Background(), "extract", input, d)
	if calls := rec.all(); len(calls) != 0 {
		t.Errorf("notifications = %v, want none", calls)
	}
	if !strings.Contains(log.String(), "cancelled") {
		t.Errorf("log should say cancelled: %q", log.String())
	}
}
