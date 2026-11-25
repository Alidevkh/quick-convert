package main

// This file connects each menu action to the package that does the work.

import (
	"strings"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/media"
	"quickconvert/internal/pdf"
	"quickconvert/internal/pdfmd"
	"quickconvert/internal/photo"
	"quickconvert/internal/textcopy"
)

// handlerFor picks the converter for an action and a file extension.
// "png" and "jpg" mean different work for a PDF (render pages) than for a
// photo (convert the image), so the extension is part of the choice.
func handlerFor(a actions.Action, ext string) app.Handler {
	switch {
	case a.ID == "txt" || strings.HasPrefix(a.ID, "code-"):
		return textcopy.Copy
	case a.ID == "clean":
		return photo.Clean
	case a.ID == "shrink":
		return pdf.Shrink
	case a.ID == "md":
		return pdfmd.Convert
	case ext == ".pdf": // png, jpg from a PDF
		return pdf.Render
	case a.ID == "mp4" || a.ID == "gif" || (a.ID == "mp3" && isVideo(ext)):
		return media.Video
	case a.ID == "mp3" || a.ID == "wav":
		return media.Audio
	default: // png, jpg, webp from a picture
		return media.Image
	}
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
