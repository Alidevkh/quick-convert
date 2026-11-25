package pdfwriter

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"quickconvert/internal/testutil"
)

// makeJPEG creates a red JPEG of w x h pixels and returns its bytes.
func makeJPEG(t *testing.T, w, h int) []byte {
	t.Helper()
	path := testutil.MakeImage(t, t.TempDir(), "p.jpg", w, h, false)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// writePDF writes pages to a file and returns its path.
func writePDF(t *testing.T, pages []Page) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "out.pdf")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := Write(f, pages); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return path
}

// mediaBoxes asks mutool for every page's MediaBox, e.g. ["0 0 144 72"].
func mediaBoxes(t *testing.T, pdf string) []string {
	t.Helper()
	out := testutil.Run(t, "mutool", "pages", pdf)
	re := regexp.MustCompile(`<MediaBox l="([\d.]+)" b="([\d.]+)" r="([\d.]+)" t="([\d.]+)"`)
	var boxes []string
	for _, m := range re.FindAllStringSubmatch(out, -1) {
		boxes = append(boxes, strings.Join(m[1:], " "))
	}
	return boxes
}

func TestWriteTwoPages(t *testing.T) {
	testutil.Need(t, "ffmpeg", "mutool")
	pdf := writePDF(t, []Page{
		{JPEG: makeJPEG(t, 600, 300), Width: 600, Height: 300},
		{JPEG: makeJPEG(t, 300, 600), Width: 300, Height: 600},
	})
	boxes := mediaBoxes(t, pdf)
	if len(boxes) != 2 || boxes[0] != "0 0 144 72" || boxes[1] != "0 0 72 144" {
		t.Errorf("boxes = %v", boxes)
	}
}

func TestWriteKeepsJPEGBytes(t *testing.T) {
	testutil.Need(t, "ffmpeg")
	jpg := makeJPEG(t, 40, 20)
	var buf bytes.Buffer
	if err := Write(&buf, []Page{{JPEG: jpg, Width: 40, Height: 20}}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), jpg) {
		t.Error("the JPEG was not stored unchanged")
	}
}

func TestWriteRenders(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "mutool")
	pdf := writePDF(t, []Page{{JPEG: makeJPEG(t, 600, 300), Width: 600, Height: 300}})
	png := filepath.Join(t.TempDir(), "page.png")
	testutil.Run(t, "mutool", "draw", "-q", "-r", "300", "-o", png, pdf)
	if w, _ := testutil.Dims(t, png); w < 599 || w > 601 {
		t.Errorf("rendered width %d, want 600", w)
	}
}

func TestWriteEmpty(t *testing.T) {
	if err := Write(&bytes.Buffer{}, nil); err == nil {
		t.Error("want an error for zero pages")
	}
}
