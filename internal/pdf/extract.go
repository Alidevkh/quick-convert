package pdf

// This file copies chosen pages of a PDF into one new PDF ("Extract pages...").

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"quickconvert/internal/app"
)

// AskPages asks the user which pages to take, for a PDF with `pages` pages.
// It returns what they typed, and false if they cancelled. The program sets
// it to a small Windows input box (see package platform); tests replace it.
var AskPages = func(pages int) (string, bool) { return "", false }

// Extract saves the chosen pages, in the chosen order, as one PDF:
// "report (pages 1-3, 5).pdf". Pages are copied as they are.
func Extract(ctx context.Context, j app.Job) (app.Result, error) {
	count, err := pageCount(ctx, j)
	if err != nil {
		return app.Result{}, err
	}
	answer, ok := AskPages(count)
	if !ok || strings.TrimSpace(answer) == "" {
		return app.Result{}, &app.Info{} // cancelled: nothing to say
	}
	pages, err := ParsePages(answer, count)
	if err != nil {
		return app.Result{}, &app.UserError{Msg: err.Error()}
	}

	list := make([]string, len(pages))
	for i, p := range pages {
		list[i] = strconv.Itoa(p)
	}
	const out = "out.pdf"
	_, err = j.Eng.Run(ctx, j.Work.Dir, "mutool", "merge",
		"-o", out, // the new PDF
		j.Work.In, strings.Join(list, ",")) // the chosen pages, e.g. "5,1,2"
	if err != nil {
		return app.Result{}, Friendly(err)
	}
	saved, err := j.Work.Commit(out, j.Base+" ("+PagesLabel(pages)+")", ".pdf")
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}, Detail: fmt.Sprintf("%d pages", len(pages))}, nil
}

// pageCount checks the PDF opens (a password-protected one gives a clear
// message; "mutool show" alone would just say 1 page) and counts its pages.
func pageCount(ctx context.Context, j app.Job) (int, error) {
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "info", j.Work.In); err != nil {
		return 0, Friendly(err)
	}
	out, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "show", j.Work.In,
		"trailer/Root/Pages/Count") // the number of pages, e.g. "12"
	if err != nil {
		return 0, Friendly(err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, &app.UserError{Msg: "Couldn't count the pages of this PDF", Err: err}
	}
	return n, nil
}
