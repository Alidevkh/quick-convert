package imagepdf

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/engines"
	"quickconvert/internal/testutil"
)

func need(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "exiftool", "mutool")
}

// sendTo runs "Send to → Images to PDF" on the files and returns the
// notification title and the exit code.
func sendTo(t *testing.T, files ...string) (string, int) {
	t.Helper()
	var title string
	code := app.RunMany(context.Background(), "imgpdf", files, app.Deps{
		Notify:     func(tt, body string) { title = tt },
		Acquire:    func() func() { return func() {} },
		Log:        &bytes.Buffer{},
		Eng:        engines.Default(),
		Handlers:   func(actions.Action, string) app.Handler { return Convert },
		StartDelay: time.Hour,
	})
	return title, code
}

// box parses "0 0 595.2 446.4" into width and height.
func box(t *testing.T, s string) (w, h float64) {
	t.Helper()
	f := strings.Fields(s)
	w, _ = strconv.ParseFloat(f[2], 64)
	h, _ = strconv.ParseFloat(f[3], 64)
	return w, h
}

func near(a, b float64) bool { return a-b < 0.5 && b-a < 0.5 }

func TestImagesToPDFThreePages(t *testing.T) {
	need(t)
	dir := t.TempDir()
	var files []string
	for _, n := range []string{"IMG_2041.jpg", "IMG_2042.jpg", "IMG_2043.jpg"} {
		files = append(files, testutil.MakeImage(t, dir, n, 60, 80, false))
	}
	title, code := sendTo(t, files[2], files[0], files[1])
	if code != 0 || title != "Done: IMG_2041.pdf (3 pages)" {
		t.Fatalf("code %d title %q", code, title)
	}
	if boxes := testutil.MediaBoxes(t, filepath.Join(dir, "IMG_2041.pdf")); len(boxes) != 3 {
		t.Errorf("pages = %v", boxes)
	}
}

func TestImagesToPDFDownscales(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "big.jpg", 4000, 3000, false)
	if _, err := testutil.RunHandler(t, Convert, in, "pdf"); err != nil {
		t.Fatal(err)
	}
	w, h := box(t, testutil.MediaBoxes(t, filepath.Join(dir, "big.pdf"))[0])
	if !near(w, 595.2) || !near(h, 446.4) {
		t.Errorf("page %vx%v, want 595.2x446.4", w, h)
	}
}

func TestImagesToPDFNoEnlarge(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "small.jpg", 600, 400, false)
	if _, err := testutil.RunHandler(t, Convert, in, "pdf"); err != nil {
		t.Fatal(err)
	}
	w, h := box(t, testutil.MediaBoxes(t, filepath.Join(dir, "small.pdf"))[0])
	if !near(w, 144) || !near(h, 96) {
		t.Errorf("page %vx%v, want 144x96", w, h)
	}
}

func TestImagesToPDFRotated(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "phone.jpg", 40, 20, false)
	testutil.SetOrientation(t, in, 6)
	if _, err := testutil.RunHandler(t, Convert, in, "pdf"); err != nil {
		t.Fatal(err)
	}
	w, h := box(t, testutil.MediaBoxes(t, filepath.Join(dir, "phone.pdf"))[0])
	if h <= w {
		t.Errorf("page %vx%v should be portrait, as shown on screen", w, h)
	}
}

func TestImagesToPDFTransparentWhite(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "logo.png", 16, 16, true)
	if _, err := testutil.RunHandler(t, Convert, in, "pdf"); err != nil {
		t.Fatal(err)
	}
	png := filepath.Join(t.TempDir(), "page.png")
	testutil.Run(t, "mutool", "draw", "-q", "-r", "300", "-o", png, filepath.Join(dir, "logo.pdf"))
	pixel, err := exec.Command("ffmpeg", "-v", "error", "-i", png, "-frames:v", "1",
		"-f", "rawvideo", "-pix_fmt", "rgb24", "-").Output()
	if err != nil || len(pixel) < 3 {
		t.Fatalf("decode: %v", err)
	}
	for i, c := range pixel[:3] {
		if c < 245 {
			t.Errorf("channel %d = %d, want white", i, c)
		}
	}
}

func TestConvertToPDFSingle(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "photo.jpg", 30, 30, false)
	r, err := testutil.RunHandler(t, Convert, in, "pdf")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "photo.pdf") || r.Detail != "1 page" {
		t.Errorf("result %+v", r)
	}
	if _, err := os.Stat(in); err != nil {
		t.Error("original is gone")
	}
}

func TestImagesToPDFOddSize(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "odd.png", 601, 401, false)
	if _, err := testutil.RunHandler(t, Convert, in, "pdf"); err != nil {
		t.Fatal(err)
	}
	w, h := box(t, testutil.MediaBoxes(t, filepath.Join(dir, "odd.pdf"))[0])
	if w != 144.24 || h != 96.24 { // exactly 601x401 px; no 1-px stretch
		t.Errorf("page %vx%v, want exactly 144.24x96.24", w, h)
	}
}

// Black must stay black and white must stay white (no washed-out grey).
func TestImagesToPDFFullRange(t *testing.T) {
	need(t)
	dir := t.TempDir()
	in := filepath.Join(dir, "doc.png")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi",
		"-i", "color=c=black:s=40x20,drawbox=x=20:y=0:w=20:h=20:c=white:t=fill", "-frames:v", "1", in)
	if _, err := testutil.RunHandler(t, Convert, in, "pdf"); err != nil {
		t.Fatal(err)
	}
	png := filepath.Join(t.TempDir(), "page.png")
	testutil.Run(t, "mutool", "draw", "-q", "-r", "300", "-o", png, filepath.Join(dir, "doc.pdf"))
	raw, err := exec.Command("ffmpeg", "-v", "error", "-i", png, "-frames:v", "1",
		"-f", "rawvideo", "-pix_fmt", "gray", "-").Output()
	if err != nil {
		t.Fatal(err)
	}
	w, _ := testutil.Dims(t, png)
	row := raw[(len(raw)/w/2)*w:] // a row in the middle
	if black, white := row[2], row[w-3]; black > 5 || white < 250 {
		t.Errorf("black = %d (want ≤5), white = %d (want ≥250)", black, white)
	}
}

// A damaged photo among many: the message must name it.
func TestImagesToPDFBadPhotoNamed(t *testing.T) {
	need(t)
	dir := t.TempDir()
	good := testutil.MakeImage(t, dir, "a.jpg", 20, 20, false)
	bad := filepath.Join(dir, "b.jpg")
	os.WriteFile(bad, []byte("this is not a picture"), 0o644)
	var body string
	app.RunMany(context.Background(), "imgpdf", []string{good, bad}, app.Deps{
		Notify:     func(title, b string) { body = b },
		Acquire:    func() func() { return func() {} },
		Log:        &bytes.Buffer{},
		Eng:        engines.Default(),
		Handlers:   func(actions.Action, string) app.Handler { return Convert },
		StartDelay: time.Hour,
	})
	if !strings.Contains(body, "b.jpg") {
		t.Errorf("message %q does not name the damaged photo", body)
	}
}

func TestImagesToPDFIPhoneHEIC(t *testing.T) {
	testutil.NeedHEIC(t)
	need(t)
	dir := t.TempDir()
	heic := testutil.IPhoneHEIC(t, dir, "IMG_6736.HEIC")
	jpg := testutil.MakeImage(t, dir, "IMG_6737.jpg", 60, 80, false)
	title, code := sendTo(t, jpg, heic)
	if code != 0 || title != "Done: IMG_6736.pdf (2 pages)" {
		t.Fatalf("code %d title %q", code, title)
	}
	w, h := box(t, testutil.MediaBoxes(t, filepath.Join(dir, "IMG_6736.pdf"))[0])
	if !near(w, 446.4) || !near(h, 595.2) { // 1860x2480 px: portrait, longest side 2480
		t.Errorf("page %vx%v, want portrait 446.4x595.2", w, h)
	}
}
