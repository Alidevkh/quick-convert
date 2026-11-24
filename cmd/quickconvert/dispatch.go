package main

import (
	"context"
	"strings"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/media"
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
	case a.ID == "mp4" || a.ID == "gif" || (a.ID == "mp3" && isVideo(ext)):
		return media.Video
	case a.ID == "mp3" || a.ID == "wav":
		return media.Audio
	default: // png, jpg, webp from a picture
		return media.Image
	}
}

// notBuiltYet stands in for converters that are not written yet.
func notBuiltYet(ctx context.Context, j app.Job) (app.Result, error) {
	return app.Result{}, &app.UserError{Msg: "Not built yet"}
}

// isVideo reports whether ext is one of the video types (mp3 from a video
// and mp3 from an audio file are handled by different converters).
func isVideo(ext string) bool {
	for _, item := range actions.Submenu(ext) {
		if item.ActionID == "gif" {
			return true // only videos can become GIFs
		}
	}
	return false
}
