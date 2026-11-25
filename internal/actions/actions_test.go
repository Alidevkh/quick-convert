package actions

import (
	"reflect"
	"testing"
)

func ids(items []MenuItem) []string {
	var out []string
	for _, it := range items {
		out = append(out, it.ActionID)
	}
	return out
}

func labels(items []MenuItem) []string {
	var out []string
	for _, it := range items {
		out = append(out, it.Label)
	}
	return out
}

func TestSubmenuVideo(t *testing.T) {
	mov := Submenu(".mov")
	if got := ids(mov); !reflect.DeepEqual(got, []string{"mp4", "whatsapp", "mp3", "m4a", "gif"}) {
		t.Errorf("ids = %v", got)
	}
	if got := labels(mov); !reflect.DeepEqual(got, []string{"MP4 (full quality, for editing)", "MP4 (WhatsApp / any device)", "MP3 (audio only)", "M4A (audio only)", "GIF"}) {
		t.Errorf("labels = %v", got)
	}
	if got := ids(Submenu(".mp4")); !reflect.DeepEqual(got, []string{"whatsapp", "mp3", "m4a", "gif"}) {
		t.Errorf(".mp4 ids = %v", got)
	}
}

func TestSubmenuImageExcludesSame(t *testing.T) {
	if got := ids(Submenu(".jpeg")); !reflect.DeepEqual(got, []string{"png", "webp", "pdf", "ico"}) {
		t.Errorf("got %v", got)
	}
}

func TestSubmenuText(t *testing.T) {
	txt := Submenu(".txt")
	if len(txt) != 14 {
		t.Fatalf("want 14 items, got %d", len(txt))
	}
	if txt[0] != (MenuItem{"code-py", ".py"}) {
		t.Errorf("first = %v", txt[0])
	}
	if got := Submenu(".py"); !reflect.DeepEqual(got, []MenuItem{{"txt", ".txt"}}) {
		t.Errorf(".py = %v", got)
	}
}

func TestSubmenuPDF(t *testing.T) {
	if got := ids(Submenu(".pdf")); !reflect.DeepEqual(got, []string{"png", "jpg", "md"}) {
		t.Errorf("got %v", got)
	}
}

func TestExtras(t *testing.T) {
	if got := Extras(".pdf"); !reflect.DeepEqual(got, []MenuItem{{"shrink", "Reduce size"}, {"split", "Split into pages"}, {"extract", "Extract pages..."}}) {
		t.Errorf(".pdf = %v", got)
	}
	if got := Extras(".heic"); !reflect.DeepEqual(got, []MenuItem{{"clean", "Remove metadata"}}) {
		t.Errorf(".heic = %v", got)
	}
	for _, ext := range []string{".mov", ".mp4", ".mp3", ".m4a", ".flac"} {
		if got := Extras(ext); !reflect.DeepEqual(got, []MenuItem{{"clean", "Remove metadata"}}) {
			t.Errorf("%s = %v", ext, got)
		}
	}
	// MKV files often carry fonts and track names that matter; leave them alone.
	if got := Extras(".mkv"); len(got) != 0 {
		t.Errorf(".mkv = %v", got)
	}
}

func TestSupports(t *testing.T) {
	if Supports("mp4", ".txt") {
		t.Error("mp4 should not apply to .txt")
	}
	if !Supports("clean", ".png") {
		t.Error("clean should apply to .png")
	}
	if Supports("mp4", ".mp4") {
		t.Error("mp4 -> mp4 should not be offered")
	}
}

func TestLookupCode(t *testing.T) {
	a, ok := Lookup("code-py")
	if !ok || a.OutExt != ".py" || a.Kind != Light {
		t.Errorf("got %+v %v", a, ok)
	}
	if _, ok := Lookup("nope"); ok {
		t.Error("unknown id found")
	}
}

func TestLookupKinds(t *testing.T) {
	for id, want := range map[string]Kind{"mp4": Heavy, "shrink": Heavy, "png": Light, "clean": Light} {
		a, _ := Lookup(id)
		if a.Kind != want {
			t.Errorf("%s kind = %v", id, a.Kind)
		}
	}
	if a, _ := Lookup("shrink"); a.Tag != " (small)" || a.OutExt != ".pdf" {
		t.Errorf("shrink = %+v", a)
	}
	if a, _ := Lookup("clean"); a.Tag != " (clean)" {
		t.Errorf("clean = %+v", a)
	}
}

func TestUppercaseExt(t *testing.T) {
	if !reflect.DeepEqual(Submenu(".MOV"), Submenu(".mov")) {
		t.Error("extension matching should ignore case")
	}
}

func TestAllExts(t *testing.T) {
	all := AllExts()
	if len(all) == 0 || all[0] > all[len(all)-1] {
		t.Errorf("not sorted or empty: %v", all)
	}
	want := map[string]bool{".mov": true, ".pdf": true, ".txt": true, ".py": true, ".heic": true}
	for _, e := range all {
		delete(want, e)
	}
	if len(want) != 0 {
		t.Errorf("missing: %v", want)
	}
}

func TestSendToItems(t *testing.T) {
	items := SendToItems()
	if len(items) != 2 || items[0].ActionID != "imgpdf" || items[1].ActionID != "merge" {
		t.Fatalf("items = %+v", items)
	}
	if !contains(items[0].Exts, ".heic") || items[0].NoneMsg != "No pictures selected" {
		t.Errorf("imgpdf = %+v", items[0])
	}
	if items[1].NoneMsg != "No PDFs selected" || !reflect.DeepEqual(items[1].Exts, []string{".pdf"}) {
		t.Errorf("merge = %+v", items[1])
	}
	if _, ok := SendToFor("mp4"); ok {
		t.Error("mp4 is not a Send to action")
	}
	if a, _ := Lookup("merge"); a.Tag != " (merged)" || a.OutExt != ".pdf" {
		t.Errorf("merge action = %+v", a)
	}
}

func TestImageSubmenuHasPDF(t *testing.T) {
	jpg := Submenu(".jpg")
	if jpg[len(jpg)-2] != (MenuItem{"pdf", "PDF"}) {
		t.Errorf(".jpg = %v", jpg)
	}
	for _, it := range Submenu(".pdf") {
		if it.ActionID == "pdf" {
			t.Error(".pdf must not offer PDF")
		}
	}
}

func TestSubmenuAudioFormats(t *testing.T) {
	if got := ids(Submenu(".mp3")); !reflect.DeepEqual(got, []string{"m4a", "wav", "flac", "ogg"}) {
		t.Errorf(".mp3 = %v", got)
	}
	if got := ids(Submenu(".wav")); !reflect.DeepEqual(got, []string{"mp3", "m4a", "flac", "ogg"}) {
		t.Errorf(".wav = %v", got)
	}
	if got := ids(Submenu(".opus")); len(got) != 5 {
		t.Errorf(".opus = %v", got)
	}
}

func TestWhatsAppAction(t *testing.T) {
	a, ok := Lookup("whatsapp")
	if !ok || a.Tag != " (WhatsApp)" || a.OutExt != ".mp4" || a.Kind != Heavy {
		t.Errorf("whatsapp = %+v", a)
	}
}

// Formats added in 0.4: more video, audio and picture types, GIF → MP4,
// icons, SVG and EPUB.
func TestMoreFormats(t *testing.T) {
	for _, ext := range []string{".m4v", ".3gp", ".mts", ".m2ts", ".mpg", ".flv"} {
		if got := ids(Submenu(ext)); !reflect.DeepEqual(got, []string{"mp4", "whatsapp", "mp3", "m4a", "gif"}) {
			t.Errorf("%s = %v", ext, got)
		}
	}
	for _, ext := range []string{".aiff", ".aif", ".amr"} {
		if got := ids(Submenu(ext)); !reflect.DeepEqual(got, []string{"mp3", "m4a", "wav", "flac", "ogg"}) {
			t.Errorf("%s = %v", ext, got)
		}
	}
	for _, ext := range []string{".tif", ".tiff", ".avif", ".jpg"} {
		if got := ids(Submenu(ext)); !reflect.DeepEqual(got, []string{"png", "jpg", "webp", "pdf", "ico"}) && ext != ".jpg" {
			t.Errorf("%s = %v", ext, got)
		}
	}
	if got := ids(Submenu(".png")); !reflect.DeepEqual(got, []string{"jpg", "webp", "pdf", "ico"}) {
		t.Errorf(".png = %v", got)
	}
	if got := ids(Submenu(".gif")); !reflect.DeepEqual(got, []string{"png", "jpg", "webp", "pdf", "ico", "mp4"}) {
		t.Errorf(".gif = %v", got)
	}
	if got := labels(Submenu(".gif")); got[len(got)-1] != "MP4 (video)" || got[4] != "ICO (icon)" {
		t.Errorf(".gif labels = %v", got)
	}
	if got := ids(Submenu(".svg")); !reflect.DeepEqual(got, []string{"png"}) {
		t.Errorf(".svg = %v", got)
	}
	if got := ids(Submenu(".epub")); !reflect.DeepEqual(got, []string{"pdf", "md"}) {
		t.Errorf(".epub = %v", got)
	}
	// Remove metadata: phone videos and Apple audio yes; camcorder, old and
	// voice-memo formats (which hold no personal data in their tags) no.
	for _, ext := range []string{".m4v", ".3gp", ".aiff", ".aif"} {
		if !Supports("clean", ext) {
			t.Errorf("%s should offer Remove metadata", ext)
		}
	}
	for _, ext := range []string{".amr", ".mts", ".m2ts", ".mpg", ".flv", ".svg", ".epub", ".gif", ".tif", ".avif"} {
		if Supports("clean", ext) {
			t.Errorf("%s should not offer Remove metadata", ext)
		}
	}
}
