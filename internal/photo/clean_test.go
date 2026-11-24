package photo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"quickconvert/internal/testutil"
)

const srgbProfile = "/usr/share/color/icc/ghostscript/srgb.icc"

// tag reads one metadata value from a file ("" if not present).
func tag(t *testing.T, path, name string) string {
	t.Helper()
	return strings.TrimSpace(testutil.Run(t, "exiftool", "-s3", "-n", "-"+name, path))
}

// pixelsMD5 decodes the picture (without rotating it) and fingerprints it.
func pixelsMD5(t *testing.T, path string) string {
	t.Helper()
	return testutil.Run(t, "ffmpeg", "-v", "error", "-noautorotate", "-i", path, "-f", "md5", "-")
}

func TestCleanRemovesGPSAndCamera(t *testing.T) {
	testutil.Need(t, "exiftool", "ffmpeg")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "photo.jpg", 32, 16, false)
	args := []string{"-q", "-overwrite_original",
		"-GPSLatitude=25.2", "-GPSLatitudeRef=N", "-Make=Apple", "-Model=iPhone 15",
		"-DateTimeOriginal=2026:01:01 10:00:00", "-Orientation#=6"}
	hasICC := false
	if _, err := os.Stat(srgbProfile); err == nil {
		args = append(args, "-ICC_Profile<="+srgbProfile)
		hasICC = true
	}
	testutil.Run(t, "exiftool", append(args, in)...)

	r, err := testutil.RunHandler(t, Clean, in, "clean")
	if err != nil {
		t.Fatal(err)
	}
	out := r.Outputs[0]
	if out != filepath.Join(dir, "photo (clean).jpg") {
		t.Fatalf("output %q", out)
	}
	for _, name := range []string{"GPSLatitude", "Make", "Model", "DateTimeOriginal"} {
		if v := tag(t, out, name); v != "" {
			t.Errorf("%s still present: %q", name, v)
		}
	}
	if v := tag(t, out, "Orientation"); v != "6" {
		t.Errorf("Orientation = %q, want 6", v)
	}
	if hasICC && tag(t, out, "ProfileDescription") == "" {
		t.Error("colour profile was removed")
	}
	if pixelsMD5(t, in) != pixelsMD5(t, out) {
		t.Error("picture pixels changed")
	}
	if tag(t, in, "GPSLatitude") == "" {
		t.Error("original was modified")
	}
}

func TestCleanPNG(t *testing.T) {
	testutil.Need(t, "exiftool", "ffmpeg")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "shot.png", 16, 16, false)
	testutil.Run(t, "exiftool", "-q", "-overwrite_original", "-Make=Test", in)
	r, err := testutil.RunHandler(t, Clean, in, "clean")
	if err != nil {
		t.Fatal(err)
	}
	if v := tag(t, r.Outputs[0], "Make"); v != "" {
		t.Errorf("Make still present: %q", v)
	}
}

func TestCleanUnicodeName(t *testing.T) {
	testutil.Need(t, "exiftool", "ffmpeg")
	dir := filepath.Join(t.TempDir(), "صور")
	os.Mkdir(dir, 0o755)
	in := testutil.MakeImage(t, dir, "عائلة 'العيد'.jpg", 16, 16, false)
	if _, err := exec.Command("exiftool", "-q", "-overwrite_original", "-Make=X", in).CombinedOutput(); err != nil {
		t.Fatal(err)
	}
	r, err := testutil.RunHandler(t, Clean, in, "clean")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(r.Outputs[0]) != "عائلة 'العيد' (clean).jpg" {
		t.Errorf("output %q", r.Outputs[0])
	}
}
