package pdf

// This file joins several PDFs into one ("Send to → Merge PDFs").

import (
	"context"
	"fmt"

	"quickconvert/internal/app"
)

// Merge saves all PDFs of the job, in order, as "<first name> (merged).pdf".
// Pages are copied as they are; nothing is re-compressed.
func Merge(ctx context.Context, j app.Job) (app.Result, error) {
	if len(j.Work.Ins) < 2 {
		return app.Result{}, &app.UserError{Msg: "Select at least two PDFs to merge"}
	}
	// mutool merge does not complain about a PDF that needs a password: it
	// silently writes a blank page. So first check every PDF can be opened.
	for _, in := range j.Work.Ins {
		if _, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "info", in); err != nil {
			return app.Result{}, Friendly(err)
		}
	}

	const out = "out.pdf"
	args := append([]string{"merge",
		"-o", out, // the joined file
	}, j.Work.Ins...) // in-001.pdf, in-002.pdf ... in order
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", args...); err != nil {
		return app.Result{}, Friendly(err)
	}
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, ".pdf")
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}, Detail: fmt.Sprintf("%d files", len(j.Work.Ins))}, nil
}
