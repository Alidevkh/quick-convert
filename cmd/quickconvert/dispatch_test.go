package main

import (
	"reflect"
	"testing"

	"quickconvert/internal/app"
	"quickconvert/internal/media"

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
