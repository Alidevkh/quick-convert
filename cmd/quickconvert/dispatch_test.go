package main

import (
	"testing"

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
