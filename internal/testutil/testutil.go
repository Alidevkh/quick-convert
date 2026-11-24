// Package testutil holds helpers shared by tests: skipping when an engine
// is not installed, and creating small sample files to convert.
package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Need skips the test unless every named program is on PATH.
func Need(t *testing.T, programs ...string) {
	t.Helper()
	for _, p := range programs {
		if _, err := exec.LookPath(p); err != nil {
			t.Skipf("%s not installed", p)
		}
	}
}

// Run runs a program and fails the test if it does not succeed.
func Run(t *testing.T, program string, args ...string) string {
	t.Helper()
	out, err := exec.Command(program, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", program, args, err, out)
	}
	return string(out)
}

// MakeVideo creates a 1-second test video (320x240) with the given video
// codec ("libx264", "libx265", "libvpx-vp9"...), optionally with a tone.
func MakeVideo(t *testing.T, dir, name, vcodec string, withAudio bool) string {
	t.Helper()
	return makeVideo(t, dir, name, vcodec, "320x240", withAudio)
}

// MakeVideoSize is MakeVideo with a chosen size, e.g. "1280x720".
func MakeVideoSize(t *testing.T, dir, name, size string) string {
	t.Helper()
	return makeVideo(t, dir, name, "libx264", size, true)
}

func makeVideo(t *testing.T, dir, name, vcodec, size string, withAudio bool) string {
	t.Helper()
	path := filepath.Join(dir, name)
	args := []string{"-hide_banner", "-loglevel", "error", "-y",
		"-f", "lavfi", "-i", "testsrc=d=1:r=15:s=" + size}
	if withAudio {
		args = append(args, "-f", "lavfi", "-i", "sine=d=1", "-c:a", "aac")
	}
	args = append(args, "-c:v", vcodec, "-pix_fmt", "yuv420p", "-shortest", path)
	Run(t, "ffmpeg", args...)
	return path
}

// MakeVideoWithSubs creates an H.264 + AAC .mov with an extra subtitle
// track, like the extra tracks phones add.
func MakeVideoWithSubs(t *testing.T, dir, name string) string {
	t.Helper()
	srt := filepath.Join(dir, "subs.srt")
	os.WriteFile(srt, []byte("1\n00:00:00,000 --> 00:00:01,000\nhello\n"), 0o644)
	path := filepath.Join(dir, name)
	Run(t, "ffmpeg", "-hide_banner", "-loglevel", "error", "-y",
		"-f", "lavfi", "-i", "testsrc=d=1:r=15:s=320x240",
		"-f", "lavfi", "-i", "sine=d=1",
		"-i", srt,
		"-map", "0", "-map", "1", "-map", "2",
		"-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-c:s", "mov_text",
		"-shortest", path)
	os.Remove(srt)
	return path
}

// MakeAudio creates a 1-second tone in the format given by name's extension.
func MakeAudio(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	Run(t, "ffmpeg", "-hide_banner", "-loglevel", "error", "-y", "-f", "lavfi", "-i", "sine=d=1", path)
	return path
}

// MakeImage creates a w x h image (red, or fully transparent).
func MakeImage(t *testing.T, dir, name string, w, h int, transparent bool) string {
	t.Helper()
	path := filepath.Join(dir, name)
	size := strconv.Itoa(w) + "x" + strconv.Itoa(h)
	source := "color=c=red:s=" + size
	if transparent {
		source = "color=c=black@0.0:s=" + size + ",format=rgba"
	}
	Run(t, "ffmpeg", "-hide_banner", "-loglevel", "error", "-y", "-f", "lavfi", "-i", source, "-frames:v", "1", path)
	return path
}

// SetOrientation writes the EXIF "which way is up" tag (1-8) into a photo.
func SetOrientation(t *testing.T, path string, n int) {
	t.Helper()
	Run(t, "exiftool", "-q", "-overwrite_original", "-Orientation#="+strconv.Itoa(n), path)
}

// Dims returns the width and height of the first video/image stream.
func Dims(t *testing.T, path string) (w, h int) {
	t.Helper()
	out := Run(t, "ffprobe", "-v", "error", "-select_streams", "v:0",
		"-show_entries", "stream=width,height", "-of", "csv=p=0", path)
	parts := strings.Split(strings.TrimSpace(out), ",")
	w, _ = strconv.Atoi(parts[0])
	h, _ = strconv.Atoi(parts[1])
	return w, h
}

// Codecs returns the codec names of a file's streams by type, e.g.
// {"video": "h264", "audio": "aac"}.
func Codecs(t *testing.T, path string) map[string]string {
	t.Helper()
	out := Run(t, "ffprobe", "-v", "error", "-show_entries", "stream=codec_type,codec_name", "-of", "csv=p=0", path)
	codecs := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			codecs[parts[1]] = parts[0]
		}
	}
	return codecs
}
