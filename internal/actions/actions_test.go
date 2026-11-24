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
	if got := ids(mov); !reflect.DeepEqual(got, []string{"mp4", "mp3", "gif"}) {
		t.Errorf("ids = %v", got)
	}
	if got := labels(mov); !reflect.DeepEqual(got, []string{"MP4", "MP3 (audio only)", "GIF"}) {
		t.Errorf("labels = %v", got)
	}
	if got := ids(Submenu(".mp4")); !reflect.DeepEqual(got, []string{"mp3", "gif"}) {
		t.Errorf(".mp4 ids = %v", got)
	}
}

func TestSubmenuImageExcludesSame(t *testing.T) {
	if got := ids(Submenu(".jpeg")); !reflect.DeepEqual(got, []string{"png", "webp"}) {
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
	if got := Extras(".pdf"); !reflect.DeepEqual(got, []MenuItem{{"shrink", "Reduce size"}}) {
		t.Errorf(".pdf = %v", got)
	}
	if got := Extras(".heic"); !reflect.DeepEqual(got, []MenuItem{{"clean", "Remove metadata"}}) {
		t.Errorf(".heic = %v", got)
	}
	if got := Extras(".mov"); len(got) != 0 {
		t.Errorf(".mov = %v", got)
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
