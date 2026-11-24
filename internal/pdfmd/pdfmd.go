package pdfmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"quickconvert/internal/app"
)

// Convert saves the PDF's text as Markdown: "report.pdf" → "report.md".
func Convert(ctx context.Context, j app.Job) (app.Result, error) {
	const xmlFile, mdFile = "out.xml", "out.md"
	_, err := j.Eng.Run(ctx, j.Work.Dir, "mutool", "draw",
		"-q",          // quiet
		"-F", "stext", // "structured text": every character with font, size and position
		"-o", xmlFile,
		j.Work.In)
	if err != nil {
		return app.Result{}, err
	}

	f, err := os.Open(filepath.Join(j.Work.Dir, xmlFile))
	if err != nil {
		return app.Result{}, err
	}
	lines, err := ParseStext(f)
	f.Close()
	if err != nil {
		return app.Result{}, err
	}

	markdown, err := ToMarkdown(lines)
	if errors.Is(err, ErrNoText) {
		return app.Result{}, &app.UserError{Msg: ErrNoText.Error()}
	}
	if err != nil {
		return app.Result{}, err
	}

	// Plain UTF-8, no byte-order mark: what AI tools and editors expect.
	if err := os.WriteFile(filepath.Join(j.Work.Dir, mdFile), []byte(markdown), 0o644); err != nil {
		return app.Result{}, err
	}
	saved, err := j.Work.Commit(mdFile, j.Base+j.Action.Tag, ".md")
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}
