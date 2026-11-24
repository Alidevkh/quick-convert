package outpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
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
