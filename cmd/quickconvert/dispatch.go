package main

import (
	"context"
	"strings"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
)

// handlerFor picks the converter for an action and a file extension.
// "png" and "jpg" mean different work for a PDF (render pages) than for a
// photo (convert the image), so the extension is part of the choice.
func handlerFor(a actions.Action, ext string) app.Handler {
	switch {
	case a.ID == "txt" || strings.HasPrefix(a.ID, "code-"):
		return notBuiltYet
	case a.ID == "clean":
		return notBuiltYet
	case a.ID == "shrink":
		return notBuiltYet
	case a.ID == "md":
		return notBuiltYet
	case ext == ".pdf": // png, jpg from a PDF
		return notBuiltYet
	default: // video, audio and image conversions
		return notBuiltYet
	}
}

// notBuiltYet stands in for converters that are not written yet.
func notBuiltYet(ctx context.Context, j app.Job) (app.Result, error) {
	return app.Result{}, &app.UserError{Msg: "Not built yet"}
}
