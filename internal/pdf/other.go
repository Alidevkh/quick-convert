package pdf

// This file converts the other document types MuPDF can read: SVG
// drawings, EPUB e-books and multi-page TIFF scans.

import (
	"context"

	"quickconvert/internal/app"
)

// svgSide is the longest side, in pixels, of a picture made from an SVG.
// An SVG is a drawing, so it can be made any size and stays sharp.
const svgSide = "2048"

// SVG turns a vector drawing (logo, icon) into a PNG, 2048 pixels on its
// longest side, keeping the transparent background.
func SVG(ctx context.Context, j app.Job) (app.Result, error) {
	_, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "draw",
		"-q",
		"-c", "rgba", // keep transparency
		"-w", svgSide, "-h", svgSide, // fit inside 2048 x 2048, keeping the proportions
		"-o", "out.png",
		j.Work.In)
	if err != nil {
		return app.Result{}, &app.UserError{Msg: "Couldn't read this SVG", Err: err}
	}
	return save(j, "out.png", ".png")
}

// ToPDF turns an EPUB e-book (text stays selectable) or a TIFF scan (every
// page, not only the first) into a PDF.
func ToPDF(ctx context.Context, j app.Job) (app.Result, error) {
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "convert", "-o", "out.pdf", j.Work.In); err != nil {
		msg := "Couldn't read this file"
		if j.Ext == ".epub" {
			msg = "Couldn't read this e-book (it may be protected)"
		}
		return app.Result{}, &app.UserError{Msg: msg, Err: err}
	}
	return save(j, "out.pdf", ".pdf")
}

func save(j app.Job, out, ext string) (app.Result, error) {
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, ext)
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}
