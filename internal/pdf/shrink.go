package pdf

// This file makes a smaller copy of a PDF with Ghostscript.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"quickconvert/internal/app"
)

// minSaving: the smaller PDF is only kept if it is at least 10% smaller.
const minSaving = 0.10

// Shrink saves a smaller copy of the PDF as "<name> (small).pdf".
// Pictures inside are reduced to about 150 DPI and recompressed with
// Ghostscript; text stays sharp. If the result is not clearly smaller (a
// PDF that is mostly text), nothing is saved and the user is told so.
func Shrink(ctx context.Context, j app.Job) (app.Result, error) {
	const out = "out.pdf"
	_, err := j.Eng.Run(ctx, j.Work.Dir, "gs",
		"-q",                   // quiet
		"-dSAFER",              // the PDF may not touch other files
		"-dBATCH", "-dNOPAUSE", // run without asking anything
		"-dPDFSTOPONERROR",  // fail on errors (e.g. password) instead of writing a broken file
		"-sDEVICE=pdfwrite", // write a new PDF
		"-dCompatibilityLevel=1.5",
		"-dPDFSETTINGS=/ebook", // the "ebook" preset: pictures at 150 DPI
		"-sOutputFile="+out,
		j.Work.In)
	if err != nil {
		return app.Result{}, Friendly(err)
	}

	before, err := fileSize(filepath.Join(j.Work.Dir, j.Work.In))
	if err != nil {
		return app.Result{}, err
	}
	after, err := fileSize(filepath.Join(j.Work.Dir, out))
	if err != nil {
		return app.Result{}, err
	}
	if float64(after) > (1-minSaving)*float64(before) {
		return app.Result{}, &app.Info{Msg: "Already optimized — no smaller file made"}
	}

	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, ".pdf")
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{
		Outputs: []string{saved},
		Detail:  FormatSize(before) + " → " + FormatSize(after),
	}, nil
}

// FormatSize writes a file size for people: "4.2 MB" or "850 KB".
func FormatSize(bytes int64) string {
	const kb, mb = 1024, 1024 * 1024
	if bytes >= mb {
		return fmt.Sprintf("%.1f MB", float64(bytes)/mb)
	}
	return fmt.Sprintf("%d KB", (bytes+kb/2)/kb)
}

func fileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
