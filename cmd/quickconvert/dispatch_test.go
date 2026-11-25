package main

import (
	"reflect"
	"testing"

	"quickconvert/internal/app"
	"quickconvert/internal/imagepdf"
	"quickconvert/internal/media"
	"quickconvert/internal/pdf"
	"quickconvert/internal/pdfmd"

	"quickconvert/internal/actions"
)

// Every entry the installer puts in the right-click menu must lead to a
// converter.
func TestHandlerForAllMenuItems(t *testing.T) {
	for _, ext := range actions.AllExts() {
		for _, item := range append(actions.Submenu(ext), actions.Extras(ext)...) {
			a, ok := actions.Lookup(item.ActionID)
			if !ok {
				t.Fatalf("%s: unknown action %s", ext, item.ActionID)
			}
			if handlerFor(a, ext) == nil {
				t.Errorf("%s %s: no handler", ext, item.ActionID)
			}
		}
	}
}

func same(a, b app.Handler) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

func TestMp3FromVideoVsAudio(t *testing.T) {
	mp3, _ := actions.Lookup("mp3")
	if !same(handlerFor(mp3, ".mov"), media.Video) || !same(handlerFor(mp3, ".mp4"), media.Video) {
		t.Error("mp3 from a video should use media.Video")
	}
	if !same(handlerFor(mp3, ".wav"), media.Audio) {
		t.Error("mp3 from audio should use media.Audio")
	}
}

func TestHandlerForRoutesPdfImages(t *testing.T) {
	png, _ := actions.Lookup("png")
	if !same(handlerFor(png, ".pdf"), pdf.Render) {
		t.Error("png from a PDF should use pdf.Render")
	}
	if !same(handlerFor(png, ".jpg"), media.Image) {
		t.Error("png from a photo should use media.Image")
	}
	shrink, _ := actions.Lookup("shrink")
	if !same(handlerFor(shrink, ".pdf"), pdf.Shrink) {
		t.Error("shrink should use pdf.Shrink")
	}
}

func TestMarkdownRouting(t *testing.T) {
	md, _ := actions.Lookup("md")
	if !same(handlerFor(md, ".pdf"), pdfmd.Convert) {
		t.Error("md should use pdfmd.Convert")
	}
	for _, ext := range actions.AllExts() {
		for _, item := range append(actions.Submenu(ext), actions.Extras(ext)...) {
			a, _ := actions.Lookup(item.ActionID)
			if handlerFor(a, ext) == nil {
				t.Errorf("%s %s has no converter", ext, item.ActionID)
			}
		}
	}
}

func TestImagesToPDFRouting(t *testing.T) {
	for _, id := range []string{"pdf", "imgpdf"} {
		a, _ := actions.Lookup(id)
		if !same(handlerFor(a, ".jpg"), imagepdf.Convert) {
			t.Errorf("%s should use imagepdf.Convert", id)
		}
	}
}

func TestMergeRouting(t *testing.T) {
	merge, _ := actions.Lookup("merge")
	if !same(handlerFor(merge, ".pdf"), pdf.Merge) {
		t.Error("merge should use pdf.Merge")
	}
}

func TestNewVideoAndAudioRouting(t *testing.T) {
	for _, c := range []struct{ id, ext string }{{"whatsapp", ".mov"}, {"whatsapp", ".mp4"}, {"m4a", ".mov"}} {
		a, _ := actions.Lookup(c.id)
		if !same(handlerFor(a, c.ext), media.Video) {
			t.Errorf("%s from %s should use media.Video", c.id, c.ext)
		}
	}
	for _, id := range []string{"m4a", "flac", "ogg", "wav", "mp3"} {
		a, _ := actions.Lookup(id)
		if !same(handlerFor(a, ".wav"), media.Audio) && !(id == "wav") {
			t.Errorf("%s from .wav should use media.Audio", id)
		}
		if !same(handlerFor(a, ".mp3"), media.Audio) && id != "mp3" {
			t.Errorf("%s from .mp3 should use media.Audio", id)
		}
	}
}
