package pdfmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"quickconvert/internal/app"
	"quickconvert/internal/testutil"
)

func TestConvertEndToEnd(t *testing.T) {
	testutil.Need(t, "gs", "mutool")
	dir := t.TempDir()
	in := testutil.MakePDFFromPS(t, dir, "doc.pdf", `%!PS
<< /PageSize [595 842] >> setpagedevice
/Helvetica-Bold findfont 20 scalefont setfont 72 720 moveto (Big Title) show
/Helvetica findfont 11 scalefont setfont
72 680 moveto (The first line of body text in this sample document.) show
72 666 moveto (The second line of body text in this sample document.) show
showpage
`)
	r, err := testutil.RunHandler(t, Convert, in, "md")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "doc.md") {
		t.Fatalf("output %v", r.Outputs)
	}
	data, _ := os.ReadFile(r.Outputs[0])
	md := string(data)
	if !strings.HasPrefix(md, "# Big Title\n\n") || !strings.Contains(md, "body text") {
		t.Errorf("got:\n%s", md)
	}
	if strings.HasPrefix(md, "\xef\xbb\xbf") {
		t.Error("file starts with a BOM")
	}
}

func TestConvertScanned(t *testing.T) {
	testutil.Need(t, "mutool", "ffmpeg")
	dir := t.TempDir()
	in := testutil.MakeImagePDF(t, dir, "scan.pdf")
	_, err := testutil.RunHandler(t, Convert, in, "md")
	var ue *app.UserError
	if !errors.As(err, &ue) || ue.Msg != "This PDF has no selectable text (scanned?)" {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "scan.md")); !os.IsNotExist(err) {
		t.Error("scan.md was created")
	}
}
