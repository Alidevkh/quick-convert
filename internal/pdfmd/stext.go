// Package pdfmd turns a text PDF into Markdown, for pasting into AI tools.
//
// It works in three steps:
//
//  1. mutool extracts every character with its position, font and size
//     ("structured text", an XML file)                      — pdfmd.go
//  2. ParseStext reads that XML into Lines of styled text      — this file
//  3. ToMarkdown decides what each line is (heading, list item, table row,
//     paragraph, page header...) and writes Markdown          — markdown.go
package pdfmd

import (
	"encoding/xml"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Span is a piece of a line in one style.
type Span struct {
	Text   string
	Bold   bool
	Italic bool
}

// Line is one line of text as it appears on the page.
type Line struct {
	Page  int     // page number, starting at 0
	Block int     // paragraph-like group the line belongs to (numbered across the document)
	PageH float64 // page height, to know what is near the top or bottom
	X0    float64 // left edge of the line
	X1    float64 // right edge of the line
	Y     float64 // baseline position, measured from the top of the page
	Size  float64 // main font size of the line
	Spans []Span
}

// Text returns the line's plain text, without styles.
func (l Line) Text() string {
	var b strings.Builder
	for _, s := range l.Spans {
		b.WriteString(s.Text)
	}
	return b.String()
}

// The XML layout produced by "mutool draw -F stext":
//
//	<page height="842"> <block> <line> <font name="Arial-Bold" size="11">
//	  <char x="72" y="100" quad="..." c="H"/> ...
type xmlPage struct {
	Height float64    `xml:"height,attr"`
	Blocks []xmlBlock `xml:"block"`
}
type xmlBlock struct {
	Lines []xmlLine `xml:"line"`
}
type xmlLine struct {
	Fonts []xmlFont `xml:"font"`
}
type xmlFont struct {
	Name  string    `xml:"name,attr"`
	Size  float64   `xml:"size,attr"`
	Chars []xmlChar `xml:"char"`
}
type xmlChar struct {
	X    float64 `xml:"x,attr"`
	Y    float64 `xml:"y,attr"`
	Quad string  `xml:"quad,attr"` // corner points; the 3rd number is the right edge
	C    string  `xml:"c,attr"`
}

// ParseStext reads mutool's structured-text XML. Pages are read one at a
// time, so even very long PDFs don't need much memory. Images are ignored.
func ParseStext(r io.Reader) ([]Line, error) {
	decoder := xml.NewDecoder(r)
	var lines []Line
	pageNum, blockNum := 0, 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return lines, nil
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "page" {
			continue
		}
		var p xmlPage
		if err := decoder.DecodeElement(&p, &start); err != nil {
			return nil, err
		}
		for _, b := range p.Blocks {
			for _, xl := range b.Lines {
				if line, ok := makeLine(xl, pageNum, blockNum, p.Height); ok {
					lines = append(lines, line)
				}
			}
			blockNum++
		}
		pageNum++
	}
}

// makeLine turns one XML line into a Line. ok is false for empty lines.
func makeLine(xl xmlLine, page, block int, pageH float64) (Line, bool) {
	line := Line{Page: page, Block: block, PageH: pageH}
	charsPerSize := map[float64]int{}
	first := true
	for _, f := range xl.Fonts {
		bold, italic := fontStyle(f.Name)
		for _, c := range f.Chars {
			if first {
				line.X0, line.Y = c.X, c.Y
				first = false
			}
			line.X1 = maxFloat(line.X1, rightEdge(c))
			if strings.TrimSpace(c.C) != "" {
				charsPerSize[f.Size]++
			}
			addChar(&line, c.C, bold, italic)
		}
	}
	trimSpans(&line)
	if len(line.Spans) == 0 {
		return Line{}, false
	}
	line.Size = mostCommon(charsPerSize)
	return line, true
}

// addChar appends a character, merging it into the last span when the
// style is the same.
func addChar(line *Line, c string, bold, italic bool) {
	n := len(line.Spans)
	if n > 0 && line.Spans[n-1].Bold == bold && line.Spans[n-1].Italic == italic {
		line.Spans[n-1].Text += c
		return
	}
	line.Spans = append(line.Spans, Span{Text: c, Bold: bold, Italic: italic})
}

// trimSpans removes spaces at the start and end of the line.
func trimSpans(line *Line) {
	for len(line.Spans) > 0 {
		first := &line.Spans[0]
		first.Text = strings.TrimLeft(first.Text, " \t")
		if first.Text != "" {
			break
		}
		line.Spans = line.Spans[1:]
	}
	for len(line.Spans) > 0 {
		last := &line.Spans[len(line.Spans)-1]
		last.Text = strings.TrimRight(last.Text, " \t")
		if last.Text != "" {
			break
		}
		line.Spans = line.Spans[:len(line.Spans)-1]
	}
}

// fontStyle guesses bold/italic from the font name, e.g.
// "ABCDEF+Arial-BoldItalicMT" → bold and italic.
func fontStyle(name string) (bold, italic bool) {
	n := strings.ToLower(name)
	for _, w := range []string{"bold", "black", "heavy", "semibold", "demi"} {
		if strings.Contains(n, w) {
			bold = true
		}
	}
	italic = strings.Contains(n, "italic") || strings.Contains(n, "oblique")
	return bold, italic
}

// rightEdge reads the right edge of a character from its quad.
func rightEdge(c xmlChar) float64 {
	parts := strings.Fields(c.Quad)
	if len(parts) >= 3 {
		if v, err := strconv.ParseFloat(parts[2], 64); err == nil {
			return v
		}
	}
	return c.X
}

// mostCommon returns the key with the highest count (smallest key on ties,
// so the result does not change from run to run).
func mostCommon(counts map[float64]int) float64 {
	keys := make([]float64, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	best, bestCount := 0.0, -1
	for _, k := range keys {
		if counts[k] > bestCount {
			best, bestCount = k, counts[k]
		}
	}
	return best
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
