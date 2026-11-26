package engines

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunUsesDir(t *testing.T) {
	dir := t.TempDir()
	out, err := runExe(context.Background(), dir, "pwd", "pwd")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := filepath.EvalSymlinks(strings.TrimSpace(string(out)))
	want, _ := filepath.EvalSymlinks(dir)
	if got != want {
		t.Errorf("ran in %q, want %q", got, want)
	}
}

func TestEngineErrorTail(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed") // (testutil would import this package: a cycle)
	}
	s := Default()
	_, err := s.Run(context.Background(), t.TempDir(), "ffmpeg", "-hide_banner", "-i", "missing.mov")
	var ee *EngineError
	if !errors.As(err, &ee) {
		t.Fatalf("want *EngineError, got %v", err)
	}
	if ee.Name != "ffmpeg" || ee.Tail == "" {
		t.Errorf("got %+v", ee)
	}
	if !strings.Contains(ee.Stderr, "missing.mov") {
		t.Errorf("full error output not kept: %q", ee.Stderr)
	}
	if len(ee.Tail) > 160 {
		t.Errorf("tail too long: %d", len(ee.Tail))
	}
}

func TestMissingEngine(t *testing.T) {
	t.Setenv("PATH", "")
	s := &Set{Dirs: []string{t.TempDir()}}
	_, err := s.Path("ffmpeg")
	var me *MissingError
	if !errors.As(err, &me) {
		t.Fatalf("want *MissingError, got %v", err)
	}
	if !strings.Contains(me.Error(), "is missing — see README.txt, \"Getting the engines yourself\"") {
		t.Errorf("message: %q", me.Error())
	}
}

func TestPathFindsInDirsFirst(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "mutool")
	if err := writeExecutable(fake); err != nil {
		t.Fatal(err)
	}
	s := &Set{Dirs: []string{dir}}
	got, err := s.Path("mutool")
	if err != nil || got != fake {
		t.Errorf("got %q, %v", got, err)
	}
}
