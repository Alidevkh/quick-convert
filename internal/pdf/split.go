package pdf

// This file splits a PDF into one PDF per page ("Split into pages").

import (
	"context"
	"fmt"
	"strconv"

	"quickconvert/internal/app"
)

// Split saves every page as its own PDF next to the original:
// "report-1.pdf", "report-2.pdf"... Pages are copied as they are.
func Split(ctx context.Context, j app.Job) (app.Result, error) {
	pages, err := pageCount(ctx, j)
	if err != nil {
		return app.Result{}, err
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
