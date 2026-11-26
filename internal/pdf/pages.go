package pdf

// This file reads the page list the user types for "Extract pages...",
// e.g. "1-3, 5", and turns a list of pages back into a short label.

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ParsePages reads a page list such as "1-3, 5" or "5, 1-3" for a PDF with
// count pages and returns the pages in the order given. Ranges may go
// backwards ("5-3" = 5, 4, 3). Arabic digits (١٢٣) and the Arabic comma
// (،) are accepted too.
func ParsePages(spec string, count int) ([]int, error) {
	spec = westernDigits(spec)
	spec = spacesAroundDash.ReplaceAllString(spec, "-") // "1 - 3" means "1-3"
	split := func(r rune) bool { return r == ',' || r == '،' || unicode.IsSpace(r) }
	var pages []int
	for _, part := range strings.FieldsFunc(spec, split) {
		from, to, err := parseRange(part)
		if err != nil {
			return nil, fmt.Errorf("Couldn't understand %q — write pages like 1-3, 5", part)
		}
		if from < 1 || to < 1 || from > count || to > count {
			return nil, fmt.Errorf("Pages must be between 1 and %d", count)
		}
		step := 1
		if to < from {
			step = -1
		}
		for p := from; p != to+step; p += step {
			pages = append(pages, p)
		}
	}
	if len(pages) == 0 {
		return nil, errors.New("No pages given — write pages like 1-3, 5")
	}
	return pages, nil
}

// spacesAroundDash matches a dash (or en dash) with spaces around it.
var spacesAroundDash = regexp.MustCompile(`\s*[-–]\s*`)

// parseRange reads "7" or "3-5" (a dash or an en dash).
func parseRange(s string) (from, to int, err error) {
	s = strings.ReplaceAll(s, "–", "-")
	a, b, isRange := strings.Cut(s, "-")
	if from, err = strconv.Atoi(a); err != nil {
		return 0, 0, err
	}
	if !isRange {
		return from, from, nil
	}
	to, err = strconv.Atoi(b)
	return from, to, err
}

// westernDigits turns Arabic-Indic (٠-٩) and Persian (۰-۹) digits into 0-9.
func westernDigits(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '٠' && r <= '٩':
			return '0' + (r - '٠')
		case r >= '۰' && r <= '۹':
			return '0' + (r - '۰')
		}
		return r
	}, s)
}

// maxLabel: longer labels become "selected pages", to keep file names short.
const maxLabel = 30

// PagesLabel describes pages for a file name: "page 7", "pages 1-3, 5",
// "pages 5-3", or "selected pages" when the list is long.
func PagesLabel(pages []int) string {
	if len(pages) == 1 {
		return "page " + strconv.Itoa(pages[0])
	}
	var parts []string
	for i := 0; i < len(pages); {
		j := i
		// Extend the run while pages go up (or down) one by one.
		for j+1 < len(pages) && abs(pages[j+1]-pages[j]) == 1 &&
			(j == i || pages[j+1]-pages[j] == pages[j]-pages[j-1]) {
			j++
		}
		if j > i {
			parts = append(parts, fmt.Sprintf("%d-%d", pages[i], pages[j]))
		} else {
			parts = append(parts, strconv.Itoa(pages[i]))
		}
		i = j + 1
	}
	label := "pages " + strings.Join(parts, ", ")
	if len(label) > maxLabel {
		return "selected pages"
	}
	return label
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
