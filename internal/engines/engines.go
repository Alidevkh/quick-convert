// Package engines finds and runs the bundled converter programs:
//
//	ffmpeg / ffprobe  video, audio and images
//	mutool            PDF pages to images, PDF text extraction
//	gs                Ghostscript, PDF size reduction
//	exiftool          photo metadata
//
// Engines always run inside the job's temporary folder (see package
// outpath) and are given short relative file names like "in.mov", so they
// never have to deal with unusual characters in real file names.
package engines

import (
	"quickconvert/internal/procjob"

	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Set knows where to look for the engines.
type Set struct {
	Dirs []string // folders searched first, in order; then the system PATH
}

// Default searches the folder of quickconvert.exe first (where the
// installer puts the engines), then, on Windows, a normal Ghostscript
// installation, then the PATH.
func Default() *Set {
	var dirs []string
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Dir(exe))
	}
	if runtime.GOOS == "windows" {
		pattern := filepath.Join(os.Getenv("ProgramFiles"), "gs", "gs*", "bin")
		found, _ := filepath.Glob(pattern)
		dirs = append(dirs, found...)
	}
	return &Set{Dirs: dirs}
}

// MissingError means an engine program could not be found.
type MissingError struct{ Name string }

func (e *MissingError) Error() string {
	return exeName(e.Name) + " is missing — see README.txt, \"Getting the engines yourself\""
}

// EngineError means an engine ran but reported a failure.
type EngineError struct {
	Name   string // engine, e.g. "ffmpeg"
	Tail   string // its last error message line, shortened for a notification
	Stderr string // everything it printed as errors (to look for known problems)
}

func (e *EngineError) Error() string {
	return e.Name + " failed: " + e.Tail
}

// Path returns the full path of an engine. name is one of
// "ffmpeg", "ffprobe", "mutool", "gs", "exiftool".
func (s *Set) Path(name string) (string, error) {
	exe := exeName(name)
	for _, dir := range s.Dirs {
		p := filepath.Join(dir, exe)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, nil
		}
	}
	if p, err := exec.LookPath(exe); err == nil {
		return p, nil
	}
	return "", &MissingError{Name: name}
}

// Run starts engine name with args inside folder dir, waits for it, and
// returns what it printed on standard output.
func (s *Set) Run(ctx context.Context, dir, name string, args ...string) ([]byte, error) {
	exe, err := s.Path(name)
	if err != nil {
		return nil, err
	}
	return runExe(ctx, dir, name, exe, args...)
}

// runExe does the actual work of Run once the program path is known.
func runExe(ctx context.Context, dir, name, exe string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	hide(cmd) // no console window on Windows

	err := cmd.Start()
	if err == nil {
		procjob.Track(cmd) // ends with QuickConvert, even if QuickConvert is killed
		err = cmd.Wait()
	}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return stdout.Bytes(), &EngineError{Name: name, Tail: lastLine(stderr.String()), Stderr: stderr.String()}
		}
		return stdout.Bytes(), err
	}
	return stdout.Bytes(), nil
}

// exeName turns an engine name into its program file name.
// Ghostscript's command-line program on Windows is "gswin64c.exe".
func exeName(name string) string {
	if runtime.GOOS != "windows" {
		return name
	}
	if name == "gs" {
		return "gswin64c.exe"
	}
	return name + ".exe"
}

// lastLine returns the last non-empty line of s, cut to 160 characters.
func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	line := strings.TrimSpace(lines[len(lines)-1])
	if r := []rune(line); len(r) > 160 {
		line = string(r[:157]) + "..."
	}
	if line == "" {
		return "unknown error"
	}
	return line
}
