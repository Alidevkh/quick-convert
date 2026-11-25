package app

// Natural sorting: "IMG_9.jpg" comes before "IMG_10.jpg", the way people
// (and phone cameras) count, instead of character-by-character order.

import (
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// SortNatural sorts file paths by their file name, ignoring upper/lower
// case and comparing runs of digits as numbers.
func SortNatural(paths []string) {
	sort.SliceStable(paths, func(i, j int) bool {
		return naturalLess(strings.ToLower(filepath.Base(paths[i])), strings.ToLower(filepath.Base(paths[j])))
	})
}

// naturalLess compares a and b piece by piece: digit runs as numbers,
// everything else character by character.
func naturalLess(a, b string) bool {
	ra, rb := []rune(a), []rune(b)
	i, j := 0, 0
	for i < len(ra) && j < len(rb) {
		if unicode.IsDigit(ra[i]) && unicode.IsDigit(rb[j]) {
			ni, nj := i, j
			for ni < len(ra) && unicode.IsDigit(ra[ni]) {
				ni++
			}
			for nj < len(rb) && unicode.IsDigit(rb[nj]) {
				nj++
			}
			// Compare the numbers without leading zeros: longer means bigger.
			da := strings.TrimLeft(string(ra[i:ni]), "0")
			db := strings.TrimLeft(string(rb[j:nj]), "0")
			if len(da) != len(db) {
				return len(da) < len(db)
			}
			if da != db {
				return da < db
			}
			i, j = ni, nj
			continue
		}
		if ra[i] != rb[j] {
			return ra[i] < rb[j]
		}
		i++
		j++
	}
	return len(ra)-i < len(rb)-j
}
