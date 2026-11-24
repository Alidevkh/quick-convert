package pdfmd

import (
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// ErrNoText means the PDF has no text layer (usually a scan).
var ErrNoText = errors.New("This PDF has no selectable text (scanned?)")

// Tuning values. All distances are in PDF points (1/72 inch).
const (
	bandFraction  = 0.08 // top/bottom 8% of a page is where headers and footers live
	repeatShare   = 0.5  // a header must repeat on at least half of the pages...
	repeatMinDocs = 3    // ...of a document with at least 3 pages
	headingRatio  = 1.15 // text 15% bigger than normal text is a heading
	sameRowPts    = 2.0  // lines this close vertically are on the same row
	columnPts     = 6.0  // table cells this close horizontally are in the same column
	maxCellRunes  = 40   // longer pieces of text are prose, not table cells
	indentPts     = 3.0  // a line indented more than this continues a list item
)

var (
	// A bullet character followed by a space: "• Apples", "- Pears".
	bulletStart = regexp.MustCompile(`^[•◦▪▫■□●○‣⁃∙·\-–*]\s+`)
	// A number followed by "." or ")": "1. First", "2) Second".
	numberStart = regexp.MustCompile(`^(\d{1,3})[.)]\s+`)
	// A line that is only a page number: "3", "Page 3", "3 of 10", "3/10".
	pageNumber = regexp.MustCompile(`^(page\s*)?#(\s*(of|/)\s*#)?$`)
	digits     = regexp.MustCompile(`\d+`)
	manyBlank  = regexp.MustCompile(`\n{3,}`)
)

// loneBullets are bullet characters that PDFs sometimes place as a
// separate piece of text in front of the item.
const loneBullets = "•◦▪▫■□●○‣⁃∙·"

// ToMarkdown turns the lines of a PDF into Markdown text.
func ToMarkdown(lines []Line) (string, error) {
	if !hasText(lines) {
		return "", ErrNoText
	}
	lines = removeHeadersFooters(lines)
	lines = mergeLoneBullets(lines)

	levels := headingLevels(lines, bodySize(lines))
	tables := findTables(lines)

	w := &writer{}
	for i, ln := range lines {
		if t, inTable := tables[i]; inTable {
			if i == t.first {
				w.flush()
				w.add(element{kind: tableKind, text: t.markdown})
			}
			continue
		}
		if level := levels[roundSize(ln.Size)]; level > 0 {
			w.heading(ln, level)
			continue
		}
		if marker, spans, ok := listStart(ln); ok {
			w.flush()
			w.open(itemKind, ln, marker+renderSpans(spans))
			continue
		}
		if w.continues(ln) {
			w.appendLine(renderSpans(ln.Spans))
			continue
		}
		w.flush()
		w.open(paragraphKind, ln, renderSpans(ln.Spans))
	}
	w.flush()
	return w.result(), nil
}

// ---------------------------------------------------------------------
// Writing the output

type kind int

const (
	paragraphKind kind = iota
	itemKind
	headingKind
	tableKind
)

// element is one finished piece of output: a paragraph, list item,
// heading or table.
type element struct {
	kind  kind
	text  string
	level int // heading level
	block int // block of the first line (used to join wrapped headings)
}

// writer collects elements. A paragraph or list item stays "open" while
// following lines may still be added to it.
type writer struct {
	done    []element
	current *element
	itemX0  float64 // left edge of the open list item
}

func (w *writer) add(e element) { w.done = append(w.done, e) }

func (w *writer) open(k kind, ln Line, text string) {
	w.current = &element{kind: k, text: text, block: ln.Block}
	w.itemX0 = ln.X0
}

// continues reports whether ln belongs to the open paragraph or list item:
// it must be in the same block, and a list item's continuation is indented.
func (w *writer) continues(ln Line) bool {
	if w.current == nil || ln.Block != w.current.block {
		return false
	}
	if w.current.kind == itemKind {
		return ln.X0 > w.itemX0+indentPts
	}
	return true
}

// appendLine adds the next line's text, re-joining words that were split
// with a hyphen at the end of a line ("exam-" + "ple" → "example").
func (w *writer) appendLine(next string) {
	prev := w.current.text
	if endsWithSplitWord(prev) && startsLower(next) {
		w.current.text = strings.TrimSuffix(prev, "-") + next
		return
	}
	w.current.text = prev + " " + next
}

// heading writes a heading. A heading wrapped over two lines (same block,
// same size) becomes one heading.
func (w *writer) heading(ln Line, level int) {
	w.flush()
	text := ln.Text()
	if n := len(w.done); n > 0 {
		last := &w.done[n-1]
		if last.kind == headingKind && last.level == level && last.block == ln.Block {
			last.text += " " + text
			return
		}
	}
	w.add(element{kind: headingKind, text: text, level: level, block: ln.Block})
}

// flush finishes the open paragraph or list item.
func (w *writer) flush() {
	if w.current == nil {
		return
	}
	if w.current.kind == paragraphKind {
		w.current.text = escapeParagraphStart(w.current.text)
	}
	w.add(*w.current)
	w.current = nil
}

// result joins all elements: list items on consecutive lines, everything
// else separated by a blank line.
func (w *writer) result() string {
	var b strings.Builder
	for i, e := range w.done {
		if i > 0 {
			if e.kind == itemKind && w.done[i-1].kind == itemKind {
				b.WriteString("\n")
			} else {
				b.WriteString("\n\n")
			}
		}
		if e.kind == headingKind {
			b.WriteString(strings.Repeat("#", e.level) + " ")
		}
		b.WriteString(e.text)
	}
	out := manyBlank.ReplaceAllString(b.String(), "\n\n")
	return strings.TrimSpace(out) + "\n"
}

// ---------------------------------------------------------------------
// Text styling

// renderSpans writes styled text: **bold**, *italic*, ***both***.
// Spaces stay outside the markers, and a literal * becomes \*.
func renderSpans(spans []Span) string {
	var b strings.Builder
	for _, s := range spans {
		text := strings.ReplaceAll(s.Text, "*", `\*`)
		marker := ""
		switch {
		case s.Bold && s.Italic:
			marker = "***"
		case s.Bold:
			marker = "**"
		case s.Italic:
			marker = "*"
		}
		core := strings.TrimSpace(text)
		if marker == "" || core == "" {
			b.WriteString(text)
			continue
		}
		lead := text[:len(text)-len(strings.TrimLeft(text, " "))]
		trail := text[len(strings.TrimRight(text, " ")):]
		b.WriteString(lead + marker + core + marker + trail)
	}
	return b.String()
}

// escapeParagraphStart stops a paragraph that happens to start with # or >
// from being read as a heading or quote.
func escapeParagraphStart(text string) string {
	if strings.HasPrefix(text, "#") || strings.HasPrefix(text, ">") {
		return `\` + text
	}
	return text
}

func endsWithSplitWord(s string) bool {
	r := []rune(s)
	return len(r) >= 2 && r[len(r)-1] == '-' && unicode.IsLetter(r[len(r)-2])
}

func startsLower(s string) bool {
	for _, r := range s {
		return unicode.IsLower(r)
	}
	return false
}

// ---------------------------------------------------------------------
// Lists

// listStart checks whether a line starts a list item. It returns the
// Markdown marker ("- " or "2. ") and the line's spans without the bullet.
func listStart(ln Line) (marker string, rest []Span, ok bool) {
	text := ln.Text()
	if m := bulletStart.FindString(text); m != "" {
		return "- ", dropPrefix(ln.Spans, len([]rune(m))), true
	}
	if m := numberStart.FindStringSubmatch(text); m != nil {
		return m[1] + ". ", dropPrefix(ln.Spans, len([]rune(m[0]))), true
	}
	return "", nil, false
}

// dropPrefix removes the first n characters from a list of spans.
func dropPrefix(spans []Span, n int) []Span {
	var out []Span
	for _, s := range spans {
		r := []rune(s.Text)
		if n >= len(r) {
			n -= len(r)
			continue
		}
		s.Text = string(r[n:])
		n = 0
		out = append(out, s)
	}
	return out
}

// mergeLoneBullets joins a bullet that stands alone ("•") with the text
// right next to it on the same row ("Item"), giving "• Item".
func mergeLoneBullets(lines []Line) []Line {
	var out []Line
	for i := 0; i < len(lines); i++ {
		ln := lines[i]
		text := ln.Text()
		if len([]rune(text)) == 1 && strings.Contains(loneBullets, text) && i+1 < len(lines) {
			next := lines[i+1]
			if next.Page == ln.Page && math.Abs(next.Y-ln.Y) <= sameRowPts {
				next.Spans = append([]Span{{Text: text + " "}}, next.Spans...)
				next.X0 = ln.X0
				out = append(out, next)
				i++ // the next line is used up
				continue
			}
		}
		out = append(out, ln)
	}
	return out
}

// ---------------------------------------------------------------------
// Headings

// bodySize is the font size used by most of the text.
func bodySize(lines []Line) float64 {
	counts := map[float64]int{}
	for _, ln := range lines {
		counts[roundSize(ln.Size)] += len([]rune(strings.TrimSpace(ln.Text())))
	}
	return mostCommon(counts)
}

// headingLevels maps each heading font size to its level: the biggest
// size is #, the next ##, and anything smaller ###.
func headingLevels(lines []Line, body float64) map[float64]int {
	seen := map[float64]bool{}
	for _, ln := range lines {
		if s := roundSize(ln.Size); s >= headingRatio*body {
			seen[s] = true
		}
	}
	var sizes []float64
	for s := range seen {
		sizes = append(sizes, s)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(sizes)))
	levels := map[float64]int{}
	for i, s := range sizes {
		levels[s] = min(i+1, 3)
	}
	return levels
}

// roundSize rounds a font size to the nearest half point.
func roundSize(s float64) float64 { return math.Round(s*2) / 2 }

// ---------------------------------------------------------------------
// Page headers and footers

// removeHeadersFooters drops lines near the top or bottom of a page that
// are page numbers, or that repeat on many pages (like a report title).
func removeHeadersFooters(lines []Line) []Line {
	pages := 0
	pagesWith := map[string]map[int]bool{} // header text → pages it is on
	for _, ln := range lines {
		pages = max(pages, ln.Page+1)
		if inBand(ln) {
			key := bandKey(ln)
			if pagesWith[key] == nil {
				pagesWith[key] = map[int]bool{}
			}
			pagesWith[key][ln.Page] = true
		}
	}

	var out []Line
	for _, ln := range lines {
		if inBand(ln) {
			key := bandKey(ln)
			repeats := pages >= repeatMinDocs && float64(len(pagesWith[key])) >= repeatShare*float64(pages)
			if repeats || pageNumber.MatchString(key) {
				continue
			}
		}
		out = append(out, ln)
	}
	return out
}

func inBand(ln Line) bool {
	return ln.PageH > 0 && (ln.Y < bandFraction*ln.PageH || ln.Y > (1-bandFraction)*ln.PageH)
}

// bandKey normalises header text so "Page 3 of 9" and "Page 4 of 9" match.
func bandKey(ln Line) string {
	return digits.ReplaceAllString(strings.ToLower(strings.TrimSpace(ln.Text())), "#")
}

// ---------------------------------------------------------------------
// Tables

// table is a run of rows whose cells line up in columns.
type table struct {
	first    int    // index of its first line in reading order
	markdown string // the finished Markdown table
}

// findTables returns, for every line that is part of a table, the table it
// belongs to.
func findTables(lines []Line) map[int]*table {
	result := map[int]*table{}
	for _, rows := range rowsByPage(lines) {
		for i := 0; i < len(rows); {
			end := i + 1
			if isCandidateRow(lines, rows[i]) {
				for end < len(rows) && isCandidateRow(lines, rows[end]) && sameColumns(lines, rows[i], rows[end]) {
					end++
				}
			}
			if end-i >= 2 {
				t := makeTable(lines, rows[i:end])
				for _, row := range rows[i:end] {
					for _, idx := range row {
						result[idx] = t
					}
				}
			}
			i = end
		}
	}
	return result
}

// rowsByPage groups line indices into rows (same height on the page),
// each row sorted left to right, rows sorted top to bottom.
func rowsByPage(lines []Line) [][][]int {
	byPage := map[int][]int{}
	maxPage := 0
	for i, ln := range lines {
		byPage[ln.Page] = append(byPage[ln.Page], i)
		maxPage = max(maxPage, ln.Page)
	}
	var all [][][]int
	for p := 0; p <= maxPage; p++ {
		idx := byPage[p]
		sort.SliceStable(idx, func(a, b int) bool { return lines[idx[a]].Y < lines[idx[b]].Y })
		var rows [][]int
		for _, i := range idx {
			n := len(rows)
			if n > 0 && lines[i].Y-lines[rows[n-1][0]].Y <= sameRowPts {
				rows[n-1] = append(rows[n-1], i)
			} else {
				rows = append(rows, []int{i})
			}
		}
		for _, row := range rows {
			sort.SliceStable(row, func(a, b int) bool { return lines[row[a]].X0 < lines[row[b]].X0 })
		}
		all = append(all, rows)
	}
	return all
}

// isCandidateRow: at least two short pieces of text side by side.
func isCandidateRow(lines []Line, row []int) bool {
	if len(row) < 2 {
		return false
	}
	for _, i := range row {
		if len([]rune(lines[i].Text())) > maxCellRunes {
			return false
		}
	}
	return true
}

// sameColumns: both rows have the same number of cells, starting at the
// same horizontal positions.
func sameColumns(lines []Line, a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if math.Abs(lines[a[k]].X0-lines[b[k]].X0) > columnPts {
			return false
		}
	}
	return true
}

// makeTable writes the rows as a Markdown table; the first row is the header.
func makeTable(lines []Line, rows [][]int) *table {
	first := rows[0][0]
	var b strings.Builder
	for r, row := range rows {
		b.WriteString("|")
		for _, i := range row {
			first = min(first, i)
			b.WriteString(" " + strings.ReplaceAll(lines[i].Text(), "|", `\|`) + " |")
		}
		b.WriteString("\n")
		if r == 0 {
			b.WriteString("|" + strings.Repeat(" --- |", len(row)) + "\n")
		}
	}
	return &table{first: first, markdown: strings.TrimSuffix(b.String(), "\n")}
}

// hasText reports whether there is any visible character at all.
func hasText(lines []Line) bool {
	for _, ln := range lines {
		if strings.TrimSpace(ln.Text()) != "" {
			return true
		}
	}
	return false
}
