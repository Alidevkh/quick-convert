package textcopy

import (
	"os"
	"path/filepath"
	"testing"

	"quickconvert/internal/testutil"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCopyIdenticalBytes(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "notes.txt")
	content := "print('hi')\r\n"
	write(t, in, content)

	r, err := testutil.RunHandler(t, Copy, in, "code-py")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "notes.py") {
		t.Fatalf("output %q", r.Outputs[0])
	}
	got, _ := os.ReadFile(r.Outputs[0])
	if string(got) != content {
		t.Errorf("content %q", got)
	}
	orig, _ := os.ReadFile(in)
	if string(orig) != content {
		t.Error("original changed")
	}
}

func TestCopyCollision(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "notes.txt")
	write(t, in, "x")
	write(t, filepath.Join(dir, "notes.py"), "existing")
	r, err := testutil.RunHandler(t, Copy, in, "code-py")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Outputs[0]) != "notes (1).py" {
		t.Errorf("output %q", r.Outputs[0])
	}
}

func TestCopyCodeToTxt(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "script.m")
	write(t, in, "disp(1)")
	r, err := testutil.RunHandler(t, Copy, in, "txt")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Outputs[0]) != "script.txt" {
		t.Errorf("output %q", r.Outputs[0])
	}
}
