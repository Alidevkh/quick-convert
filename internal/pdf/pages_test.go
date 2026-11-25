package pdf

import (
	"reflect"
	"testing"
)

func TestParsePages(t *testing.T) {
	for _, c := range []struct {
		in    string
		pages []int
	}{
		{"1-3, 5", []int{1, 2, 3, 5}},
		{"5, 1-3", []int{5, 1, 2, 3}},
		{"5-3", []int{5, 4, 3}},
		{"2 4  6", []int{2, 4, 6}},
		{"١-٣، ٥", []int{1, 2, 3, 5}}, // Arabic digits and Arabic comma
		{" 7 ", []int{7}},
	} {
		got, err := ParsePages(c.in, 10)
		if err != nil || !reflect.DeepEqual(got, c.pages) {
			t.Errorf("ParsePages(%q) = %v, %v; want %v", c.in, got, err, c.pages)
		}
	}
}

func TestParsePagesErrors(t *testing.T) {
	for in, msg := range map[string]string{
		"0":    "Pages must be between 1 and 10",
		"3-12": "Pages must be between 1 and 10",
		"abc":  `Couldn't understand "abc" — write pages like 1-3, 5`,
		"1-":   `Couldn't understand "1-" — write pages like 1-3, 5`,
	} {
		_, err := ParsePages(in, 10)
		if err == nil || err.Error() != msg {
			t.Errorf("ParsePages(%q) error = %v, want %q", in, err, msg)
		}
	}
}

func TestPagesLabel(t *testing.T) {
	if got := PagesLabel([]int{1, 2, 3, 5}); got != "pages 1-3, 5" {
		t.Errorf("got %q", got)
	}
	if got := PagesLabel([]int{5, 4, 3}); got != "pages 5-3" {
		t.Errorf("got %q", got)
	}
	if got := PagesLabel([]int{7}); got != "page 7" {
		t.Errorf("got %q", got)
	}
	long := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27}
	if got := PagesLabel(long); got != "selected pages" {
		t.Errorf("got %q", got)
	}
}
