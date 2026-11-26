package outpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestCandidate(t *testing.T) {
	if got := Candidate("trip", 0, ".mp4"); got != "trip.mp4" {
		t.Errorf("got %q", got)
	}
	if got := Candidate("photo (clean)", 1, ".jpg"); got != "photo (clean) (1).jpg" {
		t.Errorf("got %q", got)
	}
}

func TestPageName(t *testing.T) {
	if got := PageName("report", 0, 2, ".png"); got != "report-2.png" {
		t.Errorf("got %q", got)
	}
	if got := PageName("report", 1, 1, ".png"); got != "report (1)-1.png" {
		t.Errorf("got %q", got)
	}
}

// writeTemp creates a file inside the work folder, as an engine would.
func writeTemp(t *testing.T, w *Work, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(w.Dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestCommitSkipsExisting(t *testing.T) {
	parent := t.TempDir()
	os.WriteFile(filepath.Join(parent, "trip.mp4"), []byte("old0"), 0o644)
	os.WriteFile(filepath.Join(parent, "trip (1).mp4"), []byte("old1"), 0o644)

	w, err := NewWork(parent)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Cleanup()
	writeTemp(t, w, "out.mp4", "new")

	got, err := w.Commit("out.mp4", "trip", ".mp4")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(parent, "trip (2).mp4"); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if mustRead(t, got) != "new" {
		t.Error("wrong content in new file")
	}
	if mustRead(t, filepath.Join(parent, "trip.mp4")) != "old0" ||
		mustRead(t, filepath.Join(parent, "trip (1).mp4")) != "old1" {
		t.Error("an existing file was changed")
	}
	if _, err := os.Stat(filepath.Join(w.Dir, "out.mp4")); !os.IsNotExist(err) {
		t.Error("temp file should be gone after commit")
	}
}

func TestCommitConcurrentSameName(t *testing.T) {
	parent := t.TempDir()
	const n = 10
	results := make([]string, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w, err := NewWork(parent)
			if err != nil {
				t.Error(err)
				return
			}
			defer w.Cleanup()
			writeTemp(t, w, "out.mp4", fmt.Sprint(i))
			p, err := w.Commit("out.mp4", "a", ".mp4")
			if err != nil {
				t.Error(err)
			}
			results[i] = p
		}(i)
	}
	wg.Wait()
	seen := map[string]bool{}
	for _, p := range results {
		if seen[p] {
			t.Errorf("duplicate result %q", p)
		}
		seen[p] = true
	}
	files, _ := filepath.Glob(filepath.Join(parent, "a*.mp4"))
	if len(files) != n {
		t.Errorf("want %d files, got %d", n, len(files))
	}
}

func TestCommitLinkUnsupportedFallsBackToCopy(t *testing.T) {
	parent := t.TempDir()
	notSupported := func(oldname, newname string) error { return errors.New("not supported") }

	w, _ := NewWork(parent)
	defer w.Cleanup()
	w.link = notSupported
	writeTemp(t, w, "out.mp4", "first")
	got, err := w.Commit("out.mp4", "a", ".mp4")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "a.mp4" || mustRead(t, got) != "first" {
		t.Fatalf("got %q", got)
	}

	w2, _ := NewWork(parent)
	defer w2.Cleanup()
	w2.link = notSupported
	writeTemp(t, w2, "out.mp4", "second")
	got2, err := w2.Commit("out.mp4", "a", ".mp4")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got2) != "a (1).mp4" || mustRead(t, got2) != "second" {
		t.Fatalf("got %q", got2)
	}
	if mustRead(t, got) != "first" {
		t.Error("first file was overwritten")
	}
}

func TestCommitPagesPicksCommonN(t *testing.T) {
	parent := t.TempDir()
	os.WriteFile(filepath.Join(parent, "report-2.png"), []byte("x"), 0o644)
	w, _ := NewWork(parent)
	defer w.Cleanup()
	for i := 1; i <= 3; i++ {
		writeTemp(t, w, fmt.Sprintf("page-%d.png", i), fmt.Sprint(i))
	}
	got, err := w.CommitPages([]string{"page-1.png", "page-2.png", "page-3.png"}, "report", ".png")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"report (1)-1.png", "report (1)-2.png", "report (1)-3.png"}
	for i := range want {
		if filepath.Base(got[i]) != want[i] {
			t.Errorf("page %d: got %q want %q", i+1, filepath.Base(got[i]), want[i])
		}
	}
}

func TestUnicodeNames(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "مجلد أ")
	os.Mkdir(parent, 0o755)
	w, err := NewWork(parent)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Cleanup()
	writeTemp(t, w, "out.mp4", "x")
	got, err := w.Commit("out.mp4", "رحلة 'صيف'", ".mp4")
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(parent, "رحلة 'صيف'.mp4") {
		t.Errorf("got %q", got)
	}
}

func TestNewWorkNotWritable(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores folder permissions")
	}
	parent := t.TempDir()
	os.Chmod(parent, 0o555)
	defer os.Chmod(parent, 0o755)
	_, err := NewWork(parent)
	if !errors.Is(err, ErrNotWritable) {
		t.Errorf("got %v", err)
	}
}

func TestLinkInput(t *testing.T) {
	parent := t.TempDir()
	src := filepath.Join(parent, "My Trip.MOV")
	os.WriteFile(src, []byte("video"), 0o644)
	w, _ := NewWork(parent)
	defer w.Cleanup()
	if err := w.LinkInput(src); err != nil {
		t.Fatal(err)
	}
	if w.In != "in.mov" {
		t.Errorf("In = %q, want in.mov", w.In)
	}
	if mustRead(t, filepath.Join(w.Dir, w.In)) != "video" {
		t.Error("linked input has wrong content")
	}
}

func TestCleanupRemovesDir(t *testing.T) {
	parent := t.TempDir()
	w, _ := NewWork(parent)
	writeTemp(t, w, "junk.tmp", "x")
	w.Cleanup()
	if _, err := os.Stat(w.Dir); !os.IsNotExist(err) {
		t.Error("work folder still exists")
	}
	left, _ := filepath.Glob(filepath.Join(parent, ".qc-*.part"))
	if len(left) != 0 {
		t.Errorf("left behind: %v", left)
	}
}

// On Windows, deleting a hard link to a read-only file clears the
// read-only flag on the original. So a read-only original must be copied,
// never linked.
func TestLinkInputCopiesReadOnlyOriginal(t *testing.T) {
	parent := t.TempDir()
	src := filepath.Join(parent, "locked.pdf")
	os.WriteFile(src, []byte("pdf"), 0o444)
	w, _ := NewWork(parent)
	defer w.Cleanup()
	if err := w.LinkInput(src); err != nil {
		t.Fatal(err)
	}
	a, _ := os.Stat(src)
	b, _ := os.Stat(filepath.Join(w.Dir, w.In))
	if os.SameFile(a, b) {
		t.Error("read-only original was hard-linked")
	}
}

func TestCommitUsesMoveWhenLinkUnsupported(t *testing.T) {
	parent := t.TempDir()
	os.WriteFile(filepath.Join(parent, "a.mp4"), []byte("old"), 0o644)
	w, _ := NewWork(parent)
	defer w.Cleanup()
	w.link = func(string, string) error { return errors.New("not supported") }
	moved := 0
	w.move = func(oldname, newname string) error {
		if _, err := os.Lstat(newname); err == nil {
			return os.ErrExist
		}
		moved++
		return os.Rename(oldname, newname)
	}
	writeTemp(t, w, "out.mp4", "new")
	got, err := w.Commit("out.mp4", "a", ".mp4")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "a (1).mp4" || mustRead(t, got) != "new" || moved != 1 {
		t.Errorf("got %q moved=%d", got, moved)
	}
	if mustRead(t, filepath.Join(parent, "a.mp4")) != "old" {
		t.Error("existing file overwritten")
	}
}

func TestLinkInputs(t *testing.T) {
	parent := t.TempDir()
	var srcs []string
	for i, name := range []string{"a.JPG", "b.png", "c.jpg"} {
		p := filepath.Join(parent, name)
		os.WriteFile(p, []byte(fmt.Sprint(i)), 0o644)
		srcs = append(srcs, p)
	}
	w, _ := NewWork(parent)
	defer w.Cleanup()
	if err := w.LinkInputs(srcs); err != nil {
		t.Fatal(err)
	}
	want := []string{"in-001.jpg", "in-002.png", "in-003.jpg"}
	if fmt.Sprint(w.Ins) != fmt.Sprint(want) || w.In != "in-001.jpg" {
		t.Fatalf("Ins = %v In = %q", w.Ins, w.In)
	}
	for i, name := range w.Ins {
		if mustRead(t, filepath.Join(w.Dir, name)) != fmt.Sprint(i) {
			t.Errorf("%s has wrong content", name)
		}
	}
}

// Temporary folders left behind by a job that was killed (Task Manager,
// power cut) are removed by the next job in the same folder — but only
// ours, and only old ones (a job may still be running in a recent one).
func TestNewWorkRemovesOldLeftovers(t *testing.T) {
	parent := t.TempDir()
	old := filepath.Join(parent, ".qc-0badf00d.part")
	recent := filepath.Join(parent, ".qc-12345678.part")
	other := filepath.Join(parent, ".qc-notours")
	for _, d := range []string{old, recent, other} {
		if err := os.Mkdir(d, 0o755); err != nil {
			t.Fatal(err)
		}
		os.WriteFile(filepath.Join(d, "out.mp4"), []byte("x"), 0o644)
	}
	long := time.Now().Add(-25 * time.Hour)
	os.Chtimes(old, long, long)
	os.Chtimes(other, long, long)

	w, err := NewWork(parent)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Cleanup()
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("old leftover folder was not removed")
	}
	for _, keep := range []string{recent, other} {
		if _, err := os.Stat(keep); err != nil {
			t.Errorf("%s should be kept: %v", keep, err)
		}
	}
}
