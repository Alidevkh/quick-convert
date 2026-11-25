// Package imagepdf turns photos into a PDF, one photo per page.
//
// Used by two menu entries:
//   - "Send to → Images to PDF": all selected photos in one PDF;
//   - "Convert to → PDF": one photo, one page.
//
// Each photo is prepared by ffmpeg (upright, not too big, white background,
// JPEG), then package pdfwriter puts the JPEGs into the PDF as they are.
package imagepdf

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	"quickconvert/internal/app"
	"quickconvert/internal/media"
	"quickconvert/internal/pdfwriter"
)

// Convert makes "<first photo name>.pdf" from all photos of the job, in
// the order they were given (RunMany has already sorted them).
func Convert(ctx context.Context, j app.Job) (app.Result, error) {
	var pages []pdfwriter.Page
	for i, in := range j.Work.Ins {
		page, err := preparePage(ctx, j, in, fmt.Sprintf("page-%03d.jpg", i+1))
		if err != nil {
			return app.Result{}, app.ForFile(filepath.Base(j.Inputs[i]), err)
		}
		pages = append(pages, page)
	}

	const out = "out.pdf"
	f, err := os.Create(filepath.Join(j.Work.Dir, out))
	if err != nil {
		return app.Result{}, err
	}
	if err := pdfwriter.Write(f, pages); err != nil {
		f.Close()
		return app.Result{}, err
	}
	if err := f.Close(); err != nil {
		return app.Result{}, err
	}

	saved, err := j.Work.Commit(out, j.Base, ".pdf")
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}, Detail: pageCount(len(pages))}, nil
}

// preparePage runs ffmpeg on one photo and reads the resulting JPEG.
func preparePage(ctx context.Context, j app.Job, in, out string) (pdfwriter.Page, error) {
	orientation := media.Orientation(ctx, j.Eng, j.Work.Dir, in)
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "ffmpeg", media.PageArgs(in, out, orientation)...); err != nil {
		return pdfwriter.Page{}, err
	}
	data, err := os.ReadFile(filepath.Join(j.Work.Dir, out))
	if err != nil {
		return pdfwriter.Page{}, err
	}
	size, err := jpeg.DecodeConfig(bytes.NewReader(data)) // reads only the size, not the picture
	if err != nil {
		return pdfwriter.Page{}, err
	}
	return pdfwriter.Page{JPEG: data, Width: size.Width, Height: size.Height}, nil
}

// pageCount writes "1 page" or "3 pages".
func pageCount(n int) string {
	if n == 1 {
		return "1 page"
	}
	return fmt.Sprintf("%d pages", n)
}
