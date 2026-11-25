package pdf

// This file splits a PDF into one PDF per page ("Split into pages").

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"quickconvert/internal/app"
)

// Split saves every page as its own PDF next to the original:
// "report-1.pdf", "report-2.pdf"... Pages are copied as they are.
func Split(ctx context.Context, j app.Job) (app.Result, error) {
	// "mutool show" does not fail on a password-protected PDF (it just
	// says 1 page), so first check the PDF can really be opened.
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "info", j.Work.In); err != nil {
		return app.Result{}, Friendly(err)
	}
	out, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "show", j.Work.In,
		"trailer/Root/Pages/Count") // the number of pages, e.g. "12"
	if err != nil {
		return app.Result{}, Friendly(err)
	}
	pages, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return app.Result{}, &app.UserError{Msg: "Couldn't count the pages of this PDF", Err: err}
	}
	if pages < 2 {
		return app.Result{}, &app.UserError{Msg: "This PDF has only one page"}
	}

	var files []string
	for p := 1; p <= pages; p++ {
		name := fmt.Sprintf("page-%d.pdf", p)
		_, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "merge",
			"-o", name, // the new one-page PDF
			j.Work.In, strconv.Itoa(p)) // page p of the original
		if err != nil {
			return app.Result{}, Friendly(err)
		}
		files = append(files, name)
	}
	saved, err := j.Work.CommitPages(files, j.Base, ".pdf")
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: saved, Detail: fmt.Sprintf("%d pages", pages)}, nil
}
