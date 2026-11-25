// Package testutil holds helpers shared by tests: skipping when an engine
// is not installed, and creating small sample files to convert.
package testutil

import (
	"archive/zip"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/engines"
	"quickconvert/internal/outpath"
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
	source := "color=c=red:s=" + size + ",format=rgb24" // rgb24 keeps odd sizes exact
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
		parts := strings.Split(line, ",") // rotated videos add an empty extra field
		if len(parts) >= 2 {
			codecs[parts[1]] = parts[0]
		}
	}
	return codecs
}

// RunHandler runs converter h on input for action id, the same way
// app.Run does (temporary folder, linked input), and returns its result.
func RunHandler(t *testing.T, h app.Handler, input, id string) (app.Result, error) {
	t.Helper()
	a, ok := actions.Lookup(id)
	if !ok {
		t.Fatalf("unknown action %s", id)
	}
	w, err := outpath.NewWork(filepath.Dir(input))
	if err != nil {
		t.Fatal(err)
	}
	defer w.Cleanup()
	if err := w.LinkInput(input); err != nil {
		t.Fatal(err)
	}
	name := filepath.Base(input)
	j := app.Job{
		Input:  input,
		Base:   strings.TrimSuffix(name, filepath.Ext(name)),
		Ext:    strings.ToLower(filepath.Ext(name)),
		Action: a,
		Work:   w,
		Eng:    engines.Default(),
	}
	return h(context.Background(), j)
}

// MakePDF creates an A4 text PDF with the given number of pages, using
// Ghostscript. Each page has a title line and its page number.
func MakePDF(t *testing.T, dir, name string, pages int) string {
	t.Helper()
	var ps strings.Builder
	ps.WriteString("%!PS\n<< /PageSize [595 842] >> setpagedevice\n")
	for i := 1; i <= pages; i++ {
		ps.WriteString("/Helvetica findfont 14 scalefont setfont 72 760 moveto (Sample document text) show\n")
		ps.WriteString("72 40 moveto (" + strconv.Itoa(i) + ") show showpage\n")
	}
	return psToPDF(t, dir, name, ps.String())
}

// MakePDFFromPS turns a PostScript program into a PDF named name.
func MakePDFFromPS(t *testing.T, dir, name, postscript string) string {
	t.Helper()
	return psToPDF(t, dir, name, postscript)
}

func psToPDF(t *testing.T, dir, name, postscript string) string {
	t.Helper()
	src := filepath.Join(t.TempDir(), "doc.ps")
	os.WriteFile(src, []byte(postscript), 0o644)
	path := filepath.Join(dir, name)
	Run(t, "gs", "-q", "-dSAFER", "-dBATCH", "-dNOPAUSE", "-sDEVICE=pdfwrite", "-sOutputFile="+path, src)
	return path
}

// MakeImagePDF creates a one-page PDF that is just a big, noisy photo
// (like a scan): it has no text and it can be made much smaller.
func MakeImagePDF(t *testing.T, dir, name string) string {
	t.Helper()
	tmp := t.TempDir()
	jpg := filepath.Join(tmp, "scan.jpg")
	Run(t, "ffmpeg", "-hide_banner", "-loglevel", "error", "-y",
		"-f", "lavfi", "-i", "color=c=gray:s=2000x2800",
		"-vf", "noise=alls=60:allf=t", "-frames:v", "1", "-q:v", "2", jpg)
	page := filepath.Join(tmp, "page.txt")
	os.WriteFile(page, []byte("%%MediaBox 0 0 595 842\n%%Image Im0 "+jpg+"\nq 595 0 0 842 0 0 cm /Im0 Do Q\n"), 0o644)
	path := filepath.Join(dir, name)
	Run(t, "mutool", "create", "-o", path, page)
	return path
}

// MakeEncryptedPDF creates a password-protected PDF.
func MakeEncryptedPDF(t *testing.T, dir, name string) string {
	t.Helper()
	plain := MakePDF(t, t.TempDir(), "plain.pdf", 1)
	path := filepath.Join(dir, name)
	Run(t, "gs", "-q", "-dBATCH", "-dNOPAUSE", "-sDEVICE=pdfwrite",
		"-sOwnerPassword=owner", "-sUserPassword=user", "-sOutputFile="+path, plain)
	return path
}

// MediaBoxes returns every page's size in a PDF as "left bottom right top",
// e.g. ["0 0 144 72"], read with mutool.
func MediaBoxes(t *testing.T, pdf string) []string {
	t.Helper()
	out := Run(t, "mutool", "pages", pdf)
	re := regexp.MustCompile(`<MediaBox l="([\d.-]+)" b="([\d.-]+)" r="([\d.-]+)" t="([\d.-]+)"`)
	var boxes []string
	for _, m := range re.FindAllStringSubmatch(out, -1) {
		boxes = append(boxes, strings.Join(m[1:], " "))
	}
	return boxes
}

// IPhoneHEIC copies the sample iPhone photo (portrait, tiled, rotated) into
// dir under the given name and returns its path.
func IPhoneHEIC(t *testing.T, dir, name string) string {
	t.Helper()
	_, here, _, _ := runtime.Caller(0)
	data, err := os.ReadFile(filepath.Join(filepath.Dir(here), "testdata", "iphone-portrait.heic"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// NeedHEIC skips the test unless the ffmpeg on PATH can read iPhone photos
// (ffmpeg 7.1 or newer). Run the tests with a newer ffmpeg first on PATH to
// include them.
func NeedHEIC(t *testing.T) {
	t.Helper()
	Need(t, "ffmpeg", "ffprobe")
	sample := IPhoneHEIC(t, t.TempDir(), "probe.heic")
	out, err := exec.Command("ffprobe", "-v", "error", "-show_stream_groups", sample).Output()
	if err != nil || !strings.Contains(string(out), "Tile Grid") {
		t.Skip("this ffmpeg cannot read HEIC photos (needs 7.1 or newer)")
	}
}

// MakeEPUB creates a two-chapter e-book (English and Arabic text).
func MakeEPUB(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	z := zip.NewWriter(f)
	files := []struct{ name, body string }{
		{"mimetype", "application/epub+zip"}, // must be first and uncompressed
		{"META-INF/container.xml", `<?xml version="1.0"?><container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container"><rootfiles><rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/></rootfiles></container>`},
		{"OEBPS/content.opf", `<?xml version="1.0"?><package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="id"><metadata xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:identifier id="id">x</dc:identifier><dc:title>Test Book</dc:title><dc:language>en</dc:language></metadata><manifest><item id="c1" href="c1.xhtml" media-type="application/xhtml+xml"/><item id="c2" href="c2.xhtml" media-type="application/xhtml+xml"/></manifest><spine><itemref idref="c1"/><itemref idref="c2"/></spine></package>`},
		{"OEBPS/c1.xhtml", `<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><head><title>c1</title></head><body><h1>Chapter One</h1><p>The first paragraph has <b>bold words</b> in it, and it is long enough to look like a normal paragraph of a book.</p><ul><li>apple</li><li>banana</li></ul></body></html>`},
		{"OEBPS/c2.xhtml", `<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><head><title>c2</title></head><body><h1>Chapter Two</h1><p>مرحبا بالعالم، هذا نص عربي.</p></body></html>`},
	}
	for i, file := range files {
		method := zip.Deflate
		if i == 0 {
			method = zip.Store
		}
		w, err := z.CreateHeader(&zip.FileHeader{Name: file.name, Method: method})
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(file.body))
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}
