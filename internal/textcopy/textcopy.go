// Package textcopy turns a .txt file into a code file (.py, .m, .c...) and
// back. The content is copied byte for byte; only the extension changes.
// (Windows hides extensions by default, which makes this awkward by hand.)
package textcopy

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"quickconvert/internal/app"
)

// Copy saves the file again with the action's extension, e.g. notes.txt → notes.py.
func Copy(ctx context.Context, j app.Job) (app.Result, error) {
	out := "out" + j.Action.OutExt
	if err := copyFile(filepath.Join(j.Work.Dir, j.Work.In), filepath.Join(j.Work.Dir, out)); err != nil {
		return app.Result{}, err
	}
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, j.Action.OutExt)
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
