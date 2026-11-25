package main

// This file connects each menu action to the package that does the work.

import (
	"strings"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/imagepdf"
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
	case a.ID == "split":
		return pdf.Split
	case a.ID == "shrink":
		return pdf.Shrink
	case a.ID == "pdf" || a.ID == "imgpdf": // photo(s) to PDF
		return imagepdf.Convert
	case a.ID == "merge": // Send to → Merge PDFs
		return pdf.Merge
	case a.ID == "md":
		return pdfmd.Convert
	case ext == ".pdf": // png, jpg from a PDF
		return pdf.Render
	case isVideo(ext): // mp4, WhatsApp mp4, gif, and the sound as mp3 / m4a
		return media.Video
	case isAudio(a.ID): // an audio file to another audio format
		return media.Audio
	default: // png, jpg, webp from a picture
		return media.Image
	}
}

// isAudio reports whether the action makes an audio file.
func isAudio(id string) bool {
	switch id {
	case "mp3", "m4a", "wav", "flac", "ogg":
		return true
	}
	return false
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
