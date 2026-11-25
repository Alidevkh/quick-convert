package pdf

import (
	"bytes"
	"context"
	"errors"
	"os"
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

func names(paths []string) []string {
	var out []string
	for _, p := range paths {
		out = append(out, filepath.Base(p))
	}
	return out
}

func TestFormatSize(t *testing.T) {
	if got := FormatSize(4404019); got != "4.2 MB" {
		t.Errorf("got %q", got)
	}
	if got := FormatSize(870400); got != "850 KB" {
		t.Errorf("got %q", got)
	}
}

func TestRenderSinglePage(t *testing.T) {
	testutil.Need(t, "gs", "mutool", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakePDF(t, dir, "doc.pdf", 1)
	r, err := testutil.RunHandler(t, Render, in, "png")
	if err != nil {
		t.Fatal(err)
	}
	if got := names(r.Outputs); len(got) != 1 || got[0] != "doc.png" {
		t.Fatalf("outputs %v", got)
	}
	if w, _ := testutil.Dims(t, r.Outputs[0]); w < 1238 || w > 1242 {
		t.Errorf("width %d, want about 1240 (A4 at 150 DPI)", w)
	}
}

func TestRenderMultiPage(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	in := testutil.MakePDF(t, dir, "doc.pdf", 3)
	r, err := testutil.RunHandler(t, Render, in, "png")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(names(r.Outputs), ","); got != "doc-1.png,doc-2.png,doc-3.png" {
		t.Errorf("outputs %s", got)
	}
}

func TestRenderManyPagesInOrder(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	in := testutil.MakePDF(t, dir, "big.pdf", 11)
	r, err := testutil.RunHandler(t, Render, in, "png")
	if err != nil {
		t.Fatal(err)
	}
	if got := names(r.Outputs); len(got) != 11 || got[1] != "big-2.png" || got[10] != "big-11.png" {
		t.Errorf("outputs %v", got)
	}
}

func TestRenderJpg(t *testing.T) {
	testutil.Need(t, "gs", "mutool", "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakePDF(t, dir, "doc.pdf", 2)
	r, err := testutil.RunHandler(t, Render, in, "jpg")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(names(r.Outputs), ","); got != "doc-1.jpg,doc-2.jpg" {
		t.Fatalf("outputs %s", got)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["video"] != "mjpeg" {
		t.Errorf("codecs %v", c)
	}
	if left, _ := filepath.Glob(filepath.Join(dir, "*.png")); len(left) != 0 {
		t.Errorf("stray PNGs: %v", left)
	}
}

// runThroughApp runs the whole app flow so we can check cleanup too.
func runThroughApp(t *testing.T, id, input string, h app.Handler) (code int, lastBody string) {
	t.Helper()
	var bodies []string
	code = app.Run(context.Background(), id, input, app.Deps{
		Notify:     func(title, body string) { bodies = append(bodies, body) },
		Acquire:    func() func() { return func() {} },
		Log:        &bytes.Buffer{},
		Eng:        engines.Default(),
		Handlers:   func(actions.Action, string) app.Handler { return h },
		StartDelay: time.Hour,
	})
	return code, bodies[len(bodies)-1]
}

func TestRenderEncrypted(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	in := testutil.MakeEncryptedPDF(t, dir, "secret.pdf")
	code, body := runThroughApp(t, "png", in, Render)
	if code != 1 || body != "This PDF is password-protected" {
		t.Errorf("code %d body %q", code, body)
	}
	if left, _ := filepath.Glob(filepath.Join(dir, ".qc-*.part")); len(left) != 0 {
		t.Errorf("left behind: %v", left)
	}
}

func TestShrinkImagePDF(t *testing.T) {
	testutil.Need(t, "gs", "mutool", "ffmpeg")
	dir := t.TempDir()
	in := testutil.MakeImagePDF(t, dir, "scan.pdf")
	r, err := testutil.RunHandler(t, Shrink, in, "shrink")
	if err != nil {
		t.Fatal(err)
	}
	if got := names(r.Outputs); len(got) != 1 || got[0] != "scan (small).pdf" {
		t.Fatalf("outputs %v", got)
	}
	before, _ := os.Stat(in)
	after, _ := os.Stat(r.Outputs[0])
	if float64(after.Size()) > 0.9*float64(before.Size()) {
		t.Errorf("not smaller: %d → %d", before.Size(), after.Size())
	}
	if !strings.Contains(r.Detail, "→") {
		t.Errorf("detail %q", r.Detail)
	}
}

func TestShrinkTextPDFAlreadyOptimal(t *testing.T) {
	testutil.Need(t, "gs")
	dir := t.TempDir()
	in := testutil.MakePDF(t, dir, "doc.pdf", 1)
	_, err := testutil.RunHandler(t, Shrink, in, "shrink")
	var info *app.Info
	if !errors.As(err, &info) || info.Msg != "Already optimized — no smaller file made" {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "doc (small).pdf")); !os.IsNotExist(err) {
		t.Error("a (small) file was created")
	}
}

func TestShrinkEncrypted(t *testing.T) {
	testutil.Need(t, "gs")
	dir := t.TempDir()
	in := testutil.MakeEncryptedPDF(t, dir, "secret.pdf")
	code, body := runThroughApp(t, "shrink", in, Shrink)
	if code != 1 || body != "This PDF is password-protected" {
		t.Errorf("code %d body %q", code, body)
	}
	if _, err := os.Stat(filepath.Join(dir, "secret (small).pdf")); !os.IsNotExist(err) {
		t.Error("a (small) file was created")
	}
}

// labelledPDF makes a PDF whose every page says "<label> page N".
func labelledPDF(t *testing.T, dir, name, label string, pages int) string {
	t.Helper()
	ps := "%!PS\n<< /PageSize [595 842] >> setpagedevice\n"
	for i := 1; i <= pages; i++ {
		ps += "/Helvetica findfont 14 scalefont setfont 72 700 moveto (" + label + " page " + strconv.Itoa(i) + ") show showpage\n"
	}
	return testutil.MakePDFFromPS(t, dir, name, ps)
}

// mergeThroughApp runs "Send to → Merge PDFs" and returns exit code and
// the last notification (title, body).
func mergeThroughApp(t *testing.T, files ...string) (int, [2]string) {
	t.Helper()
	var last [2]string
	code := app.RunMany(context.Background(), "merge", files, app.Deps{
		Notify:     func(title, body string) { last = [2]string{title, body} },
		Acquire:    func() func() { return func() {} },
		Log:        &bytes.Buffer{},
		Eng:        engines.Default(),
		Handlers:   func(actions.Action, string) app.Handler { return Merge },
		StartDelay: time.Hour,
	})
	return code, last
}

func TestMergeThree(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	a := labelledPDF(t, dir, "part 1.pdf", "First", 1)
	b := labelledPDF(t, dir, "part 2.pdf", "Second", 2)
	c := labelledPDF(t, dir, "part 10.pdf", "Third", 3)
	code, last := mergeThroughApp(t, c, a, b)
	if code != 0 || last[0] != "Done: part 1 (merged).pdf (3 files)" {
		t.Fatalf("code %d notify %v", code, last)
	}
	out := filepath.Join(dir, "part 1 (merged).pdf")
	if boxes := testutil.MediaBoxes(t, out); len(boxes) != 6 {
		t.Fatalf("pages = %d, want 6", len(boxes))
	}
	text := testutil.Run(t, "mutool", "draw", "-q", "-F", "txt", "-o", "-", out)
	first, second, third := strings.Index(text, "First page 1"), strings.Index(text, "Second page 2"), strings.Index(text, "Third page 3")
	if first < 0 || !(first < second && second < third) {
		t.Errorf("pages in wrong order:\n%s", text)
	}
}

func TestMergeOnlyOne(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	a := testutil.MakePDF(t, dir, "a.pdf", 1)
	code, last := mergeThroughApp(t, a)
	if code != 1 || last[1] != "Select at least two PDFs to merge" {
		t.Errorf("code %d notify %v", code, last)
	}
}

func TestMergeEncrypted(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	a := testutil.MakePDF(t, dir, "a.pdf", 1)
	b := testutil.MakeEncryptedPDF(t, dir, "b.pdf")
	code, last := mergeThroughApp(t, a, b)
	if code != 1 || last[1] != "This PDF is password-protected" {
		t.Errorf("code %d notify %v", code, last)
	}
	left, _ := filepath.Glob(filepath.Join(dir, "*merged*"))
	parts, _ := filepath.Glob(filepath.Join(dir, ".qc-*.part"))
	if len(left)+len(parts) != 0 {
		t.Errorf("left behind: %v %v", left, parts)
	}
}

// PDFs that open without a password but have an "owner" password (common
// for bank statements) must still merge.
func TestMergeOwnerPasswordOnly(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	a := labelledPDF(t, dir, "a.pdf", "Alpha", 1)
	plain := labelledPDF(t, t.TempDir(), "plain.pdf", "Bravo", 1)
	b := filepath.Join(dir, "b.pdf")
	testutil.Run(t, "gs", "-q", "-dBATCH", "-dNOPAUSE", "-sDEVICE=pdfwrite", "-sOwnerPassword=owner", "-sOutputFile="+b, plain)
	code, last := mergeThroughApp(t, a, b)
	if code != 0 {
		t.Fatalf("code %d notify %v", code, last)
	}
	text := testutil.Run(t, "mutool", "draw", "-q", "-F", "txt", "-o", "-", filepath.Join(dir, "a (merged).pdf"))
	if !strings.Contains(text, "Bravo page 1") {
		t.Errorf("second PDF's text missing:\n%s", text)
	}
}
