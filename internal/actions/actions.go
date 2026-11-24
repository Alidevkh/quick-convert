// Package actions is the single list of everything QuickConvert can do:
// which file types get which right-click menu entries.
//
// Both the program (to check a request is valid) and the installer
// generator (to create the registry entries) read this list, so the menu
// and the program can never disagree. To add a new format, add it here.
package actions

import (
	"sort"
	"strings"
)

// Kind says how heavy a job is. At most two Heavy jobs run at the same
// time, so selecting 30 videos does not freeze the computer.
type Kind int

const (
	Light Kind = iota
	Heavy
)

// Action is one thing the program can do, identified by its ID. The ID is
// what the right-click menu passes on the command line:
// quickconvert.exe <ID> "<file>".
type Action struct {
	ID     string
	Label  string // short name used in notifications, e.g. "MP4"
	Kind   Kind
	OutExt string // extension of the result; "" means "same as the input"
	Tag    string // added to the result's name, e.g. " (clean)"
}

// MenuItem is one entry in a right-click menu.
type MenuItem struct {
	ActionID string
	Label    string // text shown in the menu
}

// codeExts are the programming-language extensions offered for .txt files
// (and which can be turned back into .txt).
var codeExts = []string{
	".py", ".m", ".c", ".cpp", ".h", ".js", ".ts",
	".html", ".css", ".java", ".json", ".md", ".bat", ".sh",
}

// all holds every action, by ID.
var all = map[string]Action{
	"mp4":    {ID: "mp4", Label: "MP4", Kind: Heavy, OutExt: ".mp4"},
	"mp3":    {ID: "mp3", Label: "MP3", Kind: Heavy, OutExt: ".mp3"},
	"gif":    {ID: "gif", Label: "GIF", Kind: Heavy, OutExt: ".gif"},
	"wav":    {ID: "wav", Label: "WAV", Kind: Heavy, OutExt: ".wav"},
	"png":    {ID: "png", Label: "PNG", Kind: Light, OutExt: ".png"},
	"jpg":    {ID: "jpg", Label: "JPG", Kind: Light, OutExt: ".jpg"},
	"webp":   {ID: "webp", Label: "WEBP", Kind: Light, OutExt: ".webp"},
	"md":     {ID: "md", Label: "MD", Kind: Light, OutExt: ".md"},
	"txt":    {ID: "txt", Label: ".txt", Kind: Light, OutExt: ".txt"},
	"clean":  {ID: "clean", Label: "Remove metadata", Kind: Light, Tag: " (clean)"},
	"shrink": {ID: "shrink", Label: "Reduce size", Kind: Heavy, OutExt: ".pdf", Tag: " (small)"},
}

// group is a set of file types that share the same "Convert to" submenu.
type group struct {
	exts  []string
	items []MenuItem
}

var groups = []group{
	{ // video
		exts:  []string{".mov", ".mp4", ".mkv", ".avi", ".webm", ".wmv"},
		items: []MenuItem{{"mp4", "MP4"}, {"mp3", "MP3 (audio only)"}, {"gif", "GIF"}},
	},
	{ // audio
		exts:  []string{".wav", ".m4a", ".flac", ".ogg", ".aac", ".wma"},
		items: []MenuItem{{"mp3", "MP3"}, {"wav", "WAV"}},
	},
	{ // image
		exts:  []string{".png", ".jpg", ".jpeg", ".webp", ".bmp", ".heic"},
		items: []MenuItem{{"png", "PNG"}, {"jpg", "JPG"}, {"webp", "WEBP"}},
	},
	{ // pdf
		exts:  []string{".pdf"},
		items: []MenuItem{{"png", "PNG"}, {"jpg", "JPG"}, {"md", "MD"}},
	},
	{ // text -> code
		exts:  []string{".txt"},
		items: codeMenuItems(),
	},
	{ // code -> text
		exts:  codeExts,
		items: []MenuItem{{"txt", ".txt"}},
	},
}

// extras are the entries shown directly in the right-click menu (not
// inside "Convert to").
var extras = []group{
	{exts: []string{".jpg", ".jpeg", ".png", ".heic", ".webp"}, items: []MenuItem{{"clean", "Remove metadata"}}},
	{exts: []string{".pdf"}, items: []MenuItem{{"shrink", "Reduce size"}}},
}

// codeMenuItems creates the ".py", ".m", ... entries and registers a
// "code-py", "code-m", ... action for each one.
func codeMenuItems() []MenuItem {
	var items []MenuItem
	for _, ext := range codeExts {
		id := "code-" + strings.TrimPrefix(ext, ".")
		all[id] = Action{ID: id, Label: ext, Kind: Light, OutExt: ext}
		items = append(items, MenuItem{id, ext})
	}
	return items
}

// Lookup finds an action by its ID.
func Lookup(id string) (Action, bool) {
	a, ok := all[id]
	return a, ok
}

// Submenu returns the "Convert to" entries for a file extension, in menu
// order. Converting a file to its own format is not offered.
func Submenu(ext string) []MenuItem {
	return itemsFor(groups, ext, true)
}

// Extras returns the top-level entries ("Remove metadata", "Reduce size").
func Extras(ext string) []MenuItem {
	return itemsFor(extras, ext, false)
}

// Supports reports whether action id is offered for files with ext.
func Supports(id, ext string) bool {
	for _, it := range append(Submenu(ext), Extras(ext)...) {
		if it.ActionID == id {
			return true
		}
	}
	return false
}

// AllExts returns every extension that has at least one entry, sorted.
func AllExts() []string {
	seen := map[string]bool{}
	for _, g := range append(append([]group{}, groups...), extras...) {
		for _, e := range g.exts {
			seen[e] = true
		}
	}
	var out []string
	for e := range seen {
		out = append(out, e)
	}
	sort.Strings(out)
	return out
}

// itemsFor collects the entries of every group that contains ext.
func itemsFor(list []group, ext string, skipSameFormat bool) []MenuItem {
	ext = strings.ToLower(ext)
	var out []MenuItem
	for _, g := range list {
		if !contains(g.exts, ext) {
			continue
		}
		for _, it := range g.items {
			if skipSameFormat && sameFormat(all[it.ActionID].OutExt, ext) {
				continue
			}
			out = append(out, it)
		}
	}
	return out
}

// sameFormat treats ".jpeg" and ".jpg" as the same format.
func sameFormat(a, b string) bool {
	normal := func(e string) string {
		if e == ".jpeg" {
			return ".jpg"
		}
		return e
	}
	return normal(a) == normal(b)
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
