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
// It returns what they typed, false if they cancelled, and an error if the
// question could not be asked. The program sets it to a small Windows box
// (see package platform); tests replace it.
var AskPages = func(pages int) (string, bool, error) { return "", false, nil }

// Extract saves the chosen pages, in the chosen order, as one PDF:
// "report (pages 1-3, 5).pdf". Pages are copied as they are.
func Extract(ctx context.Context, j app.Job) (app.Result, error) {
	count, err := pageCount(ctx, j)
	if err != nil {
		return app.Result{}, err
	}
	answer, ok, err := AskPages(count)
	if err != nil {
		return app.Result{}, &app.UserError{Msg: "Couldn't show the page box (PowerShell may be blocked)", Err: err}
	}
	answer = cleanAnswer(answer)
	if !ok || answer == "" {
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
// message) and reads its number of pages from "mutool info".
func pageCount(ctx context.Context, j app.Job) (int, error) {
	info, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "info", j.Work.In)
	if err != nil {
		return 0, Friendly(err)
	}
	n, ok := pagesFromInfo(string(info))
	if !ok {
		return 0, &app.UserError{Msg: "Couldn't count the pages of this PDF"}
	}
	return n, nil
}

// pagesFromInfo reads the "Pages: 12" line that "mutool info" prints.
func pagesFromInfo(info string) (int, bool) {
	for _, line := range strings.Split(info, "\n") {
		if rest, found := strings.CutPrefix(strings.TrimSpace(line), "Pages:"); found {
			n, err := strconv.Atoi(strings.TrimSpace(rest))
			return n, err == nil
		}
	}
	return 0, false
}

// cleanAnswer removes spaces, line breaks and a stray byte-order mark
// (U+FEFF) around what the user typed.
func cleanAnswer(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "\ufeff"))
}
