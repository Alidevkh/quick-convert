// Package pdf turns PDF pages into images (Render) and makes PDFs smaller
// (Shrink). PDF → Markdown is in its own package, pdfmd.
package pdf

import (
	"context"
	"errors"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"quickconvert/internal/app"
	"quickconvert/internal/engines"
)

// dpi is the image resolution for PDF pages: sharp on screen, with a
// reasonable file size. An A4 page becomes 1240 x 1754 pixels.
const dpi = "150"

// Render saves every page of the PDF as a PNG or JPG next to the original:
// "report.png" for a one-page PDF, "report-1.png", "report-2.png"... otherwise.
func Render(ctx context.Context, j app.Job) (app.Result, error) {
	_, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "draw",
		"-q",      // quiet
		"-r", dpi, // resolution
		"-F", "png", // output format
		"-o", "page-%d.png", // one file per page: page-1.png, page-2.png...
		j.Work.In)
	if err != nil {
		return app.Result{}, Friendly(err)
	}

	pages, err := pageFiles(j.Work.Dir, ".png")
	if err != nil {
		return app.Result{}, err
	}
	ext := j.Action.OutExt
	if ext == ".jpg" {
		if pages, err = toJPG(ctx, j, pages); err != nil {
			return app.Result{}, err
		}
	}

	if len(pages) == 1 {
		saved, err := j.Work.Commit(pages[0], j.Base, ext)
		if err != nil {
			return app.Result{}, err
		}
		return app.Result{Outputs: []string{saved}}, nil
	}
	saved, err := j.Work.CommitPages(pages, j.Base, ext)
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: saved}, nil
}

// toJPG converts each page-N.png into page-N.jpg with ffmpeg.
func toJPG(ctx context.Context, j app.Job, pngs []string) ([]string, error) {
	var jpgs []string
	for _, png := range pngs {
		jpg := strings.TrimSuffix(png, ".png") + ".jpg"
		_, err := j.Eng.Run(ctx, j.Work.Dir, "ffmpeg",
			"-hide_banner", "-nostdin", "-y",
			"-i", png,
			"-q:v", "2", // high JPG quality
			jpg)
		if err != nil {
			return nil, err
		}
		jpgs = append(jpgs, jpg)
	}
	return jpgs, nil
}

// pageFiles lists page-1<ext>, page-2<ext>... in the work folder, in page
// order (so page-10 comes after page-9, not after page-1).
func pageFiles(dir, ext string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "page-*"+ext))
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, &app.UserError{Msg: "This PDF has no pages"}
	}
	number := func(path string) int {
		n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(filepath.Base(path), "page-"), ext))
		return n
	}
	sort.Slice(matches, func(a, b int) bool { return number(matches[a]) < number(matches[b]) })
	var names []string
	for _, m := range matches {
		names = append(names, filepath.Base(m))
	}
	return names, nil
}

// Friendly replaces a known engine failure (a password-protected PDF)
// with a clear message. Other errors are returned unchanged.
func Friendly(err error) error {
	var engineErr *engines.EngineError
	if errors.As(err, &engineErr) && strings.Contains(strings.ToLower(engineErr.Stderr), "password") {
		return &app.UserError{Msg: "This PDF is password-protected"}
	}
	return err
}
