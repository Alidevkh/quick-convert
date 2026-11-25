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
	Asks   bool   // asks the user a question first (no "Converting…" notice meanwhile)
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
	"mp4":  {ID: "mp4", Label: "MP4", Kind: Heavy, OutExt: ".mp4"},
	"mp3":  {ID: "mp3", Label: "MP3", Kind: Heavy, OutExt: ".mp3"},
	"gif":  {ID: "gif", Label: "GIF", Kind: Heavy, OutExt: ".gif"},
	"wav":  {ID: "wav", Label: "WAV", Kind: Heavy, OutExt: ".wav"},
	"m4a":  {ID: "m4a", Label: "M4A", Kind: Heavy, OutExt: ".m4a"},
	"flac": {ID: "flac", Label: "FLAC", Kind: Heavy, OutExt: ".flac"},
	"ogg":  {ID: "ogg", Label: "OGG", Kind: Heavy, OutExt: ".ogg"},

	// A video that plays on any phone or computer and is ready for WhatsApp.
	"whatsapp": {ID: "whatsapp", Label: "MP4 (WhatsApp)", Kind: Heavy, OutExt: ".mp4", Tag: " (WhatsApp)"},
	"png":      {ID: "png", Label: "PNG", Kind: Light, OutExt: ".png"},
	"jpg":      {ID: "jpg", Label: "JPG", Kind: Light, OutExt: ".jpg"},
	"webp":     {ID: "webp", Label: "WEBP", Kind: Light, OutExt: ".webp"},
	"md":       {ID: "md", Label: "MD", Kind: Light, OutExt: ".md"},
	"txt":      {ID: "txt", Label: ".txt", Kind: Light, OutExt: ".txt"},
	"clean":    {ID: "clean", Label: "Remove metadata", Kind: Light, Tag: " (clean)"},
	"shrink":   {ID: "shrink", Label: "Reduce size", Kind: Heavy, OutExt: ".pdf", Tag: " (small)"},
	"pdf":      {ID: "pdf", Label: "PDF", Kind: Light, OutExt: ".pdf"},
	"split":    {ID: "split", Label: "Split into pages", Kind: Light, OutExt: ".pdf"},
	"extract":  {ID: "extract", Label: "Extract pages", Kind: Light, OutExt: ".pdf", Asks: true},
	"imgpdf":   {ID: "imgpdf", Label: "Images to PDF", Kind: Light, OutExt: ".pdf"},
	"merge":    {ID: "merge", Label: "Merge PDFs", Kind: Light, OutExt: ".pdf", Tag: " (merged)"},
	"ico":      {ID: "ico", Label: "ICO", Kind: Light, OutExt: ".ico"}, // Windows icon
}

// videoExts and audioExts are the video and audio types QuickConvert understands.
var (
	videoExts = []string{".mov", ".mp4", ".m4v", ".mkv", ".avi", ".webm", ".wmv",
		".3gp", ".mts", ".m2ts", ".mpg", ".flv"}
	audioExts = []string{".wav", ".m4a", ".flac", ".ogg", ".aac", ".wma", ".mp3", ".opus",
		".aiff", ".aif", ".amr"}
)

// imageExts are the picture types QuickConvert understands.
var imageExts = []string{".png", ".jpg", ".jpeg", ".webp", ".bmp", ".heic", ".gif", ".tif", ".tiff", ".avif"}

// group is a set of file types that share the same "Convert to" submenu.
type group struct {
	exts  []string
	items []MenuItem
}

var groups = []group{
	{ // video
		exts: videoExts,
		items: []MenuItem{
			{"mp4", "MP4 (full quality, for editing)"}, {"whatsapp", "MP4 (WhatsApp / any device)"},
			{"mp3", "MP3 (audio only)"}, {"m4a", "M4A (audio only)"}, {"gif", "GIF"},
		},
	},
	{ // audio
		exts:  audioExts,
		items: []MenuItem{{"mp3", "MP3"}, {"m4a", "M4A"}, {"wav", "WAV"}, {"flac", "FLAC"}, {"ogg", "OGG"}},
	},
	{ // image
		exts:  imageExts,
		items: []MenuItem{{"png", "PNG"}, {"jpg", "JPG"}, {"webp", "WEBP"}, {"pdf", "PDF"}, {"ico", "ICO (icon)"}},
	},
	{ // an animated GIF can also become a video (plays in every app)
		exts:  []string{".gif"},
		items: []MenuItem{{"mp4", "MP4 (video)"}},
	},
	{ // a vector drawing (logo, icon) as a sharp picture
		exts:  []string{".svg"},
		items: []MenuItem{{"png", "PNG"}},
	},
	{ // e-book
		exts:  []string{".epub"},
		items: []MenuItem{{"pdf", "PDF"}, {"md", "MD"}},
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
	{exts: []string{".pdf"}, items: []MenuItem{{"shrink", "Reduce size"}, {"split", "Split into pages"}, {"extract", "Extract pages..."}}},
	{exts: append(append([]string{}, cleanVideoExts...), cleanAudioExts...), items: []MenuItem{{"clean", "Remove metadata"}}},
}

// cleanVideoExts and cleanAudioExts are the types offered "Remove metadata".
// Left out: MKV (its attached fonts and track names would be lost), and
// camcorder (.mts .m2ts), old (.mpg .flv) and voice-memo (.amr) files,
// whose tags hold no personal information.
var (
	cleanVideoExts = []string{".mov", ".mp4", ".m4v", ".avi", ".webm", ".wmv", ".3gp"}
	cleanAudioExts = []string{".wav", ".m4a", ".flac", ".ogg", ".aac", ".wma", ".mp3", ".opus", ".aiff", ".aif"}
)

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
			// Converting a file to its own format is pointless, unless the
			// action makes something different (it has a tag, e.g. "(WhatsApp)").
			a := all[it.ActionID]
			if skipSameFormat && a.Tag == "" && sameFormat(a.OutExt, ext) {
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

// SendTo is an entry in Windows' "Send to" menu. Unlike the right-click
// entries, Windows gives a Send to program ALL selected files in one run,
// which is what combining many files into one needs.
type SendTo struct {
	ActionID string
	Label    string   // shortcut name shown in the Send to menu
	Exts     []string // file types it uses; other selected files are skipped
	NoneMsg  string   // shown when none of the selected files can be used
}

var sendTo = []SendTo{
	{ActionID: "imgpdf", Label: "Images to PDF", Exts: imageExts, NoneMsg: "No pictures selected"},
	{ActionID: "merge", Label: "Merge PDFs", Exts: []string{".pdf"}, NoneMsg: "No PDFs selected"},
}

// SendToItems returns the Send to entries, in the order they are installed.
func SendToItems() []SendTo {
	return sendTo
}

// SendToFor finds the Send to entry for an action ID.
func SendToFor(id string) (SendTo, bool) {
	for _, s := range sendTo {
		if s.ActionID == id {
			return s, true
		}
	}
	return SendTo{}, false
}

// IsImage reports whether ext is a picture type (used to pick the right
// "Remove metadata" converter: photos and videos/audio differ).
func IsImage(ext string) bool { return contains(imageExts, strings.ToLower(ext)) }
