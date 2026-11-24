package pdfmd

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// convertXML parses fixture XML and turns it into Markdown.
func convertXML(t *testing.T, xml string) (string, error) {
	t.Helper()
	lines, err := ParseStext(strings.NewReader(xml))
	if err != nil {
		t.Fatal(err)
	}
	return ToMarkdown(lines)
}

func mustConvert(t *testing.T, pages ...page) string {
	t.Helper()
	md, err := convertXML(t, stextXML(pages...))
	if err != nil {
		t.Fatal(err)
	}
	return md
}

func TestHeadings(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{l(72, 100, r("Helvetica-Bold", 24, "Title"))},
		{l(72, 140, r("Helvetica", 16, "Subtitle"))},
		{body(72, 180, "This is an exam-"), body(72, 194, "ple text that wraps.")},
		{body(72, 230, "Second paragraph.")},
	}})
	if !strings.HasPrefix(md, "# Title\n\n## Subtitle\n\n") {
		t.Errorf("start:\n%s", md)
	}
	if !strings.Contains(md, "This is an example text that wraps.\n\nSecond paragraph.") {
		t.Errorf("paragraphs:\n%s", md)
	}
}

func TestWrappedHeadingJoined(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{l(72, 100, r("Helvetica-Bold", 24, "A Very Long")), l(72, 128, r("Helvetica-Bold", 24, "Title"))},
		{body(72, 180, "Body one."), body(72, 194, "Body two.")},
	}})
	if !strings.HasPrefix(md, "# A Very Long Title\n\n") {
		t.Errorf("got:\n%s", md)
	}
}

func TestLists(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{body(72, 100, "Fruit:")},
		{body(72, 120, "• Apples"), body(72, 134, "• Pears")},
		{body(72, 160, "2) Second")},
		{body(72, 190, "•")},
		{body(86, 190, "Item")},
	}})
	for _, want := range []string{"- Apples\n- Pears", "2. Second", "- Item"} {
		if !strings.Contains(md, want) {
			t.Errorf("missing %q in:\n%s", want, md)
		}
	}
}

func TestListContinuationLine(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{body(72, 100, "• First point that"), body(86, 114, "continues here")},
	}})
	if !strings.Contains(md, "- First point that continues here") {
		t.Errorf("got:\n%s", md)
	}
}

func TestTable(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{body(72, 80, "Prices below.")},
		{body(72, 100, "Name")}, {body(200, 100, "Qty")}, {body(330, 100, "Price")},
		{body(72, 116, "Apple")}, {body(200, 116, "3")}, {body(330, 116, "1.20")},
		{body(72, 132, "Pear|Big")}, {body(200, 132, "5")}, {body(330, 132, "2.00")},
		{body(72, 170, "After the table.")},
	}})
	want := "| Name | Qty | Price |\n| --- | --- | --- |\n| Apple | 3 | 1.20 |\n| Pear\\|Big | 5 | 2.00 |"
	if !strings.Contains(md, want) {
		t.Errorf("got:\n%s", md)
	}
	if !strings.Contains(md, "Prices below.\n\n| Name") || !strings.Contains(md, "2.00 |\n\nAfter the table.") {
		t.Errorf("table not in place:\n%s", md)
	}
}

// Two columns of long prose side by side must not become a table.
func TestTwoColumnProseIsNotTable(t *testing.T) {
	left := "This is a long line of prose in the left column of a page"
	right := "And this is a long line of prose in the right column too"
	md := mustConvert(t, page{842, []block{
		{body(40, 100, left), body(40, 114, left)},
		{body(310, 100, right), body(310, 114, right)},
	}})
	if strings.Contains(md, "|") {
		t.Errorf("prose became a table:\n%s", md)
	}
}

func TestHeaderFooterRemoved(t *testing.T) {
	var pages []page
	for i := 1; i <= 4; i++ {
		pages = append(pages, page{842, []block{
			{body(72, 30, "ACME Report")},
			{body(72, 400, fmt.Sprintf("Body text of page %d.", i))},
			{body(280, 820, fmt.Sprintf("Page %d of 4", i))},
		}})
	}
	md := mustConvert(t, pages...)
	if strings.Contains(md, "ACME Report") || strings.Contains(md, "Page ") {
		t.Errorf("header/footer kept:\n%s", md)
	}
	for i := 1; i <= 4; i++ {
		if !strings.Contains(md, fmt.Sprintf("Body text of page %d.", i)) {
			t.Errorf("body %d missing:\n%s", i, md)
		}
	}
}

func TestLonePageNumberRemovedEvenInShortDoc(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{body(72, 400, "Only body.")},
		{body(290, 820, "1")},
	}})
	if strings.TrimSpace(md) != "Only body." {
		t.Errorf("got:\n%q", md)
	}
}

func TestStyles(t *testing.T) {
	md := mustConvert(t, page{842, []block{
		{l(72, 100,
			r("Helvetica", 11, "This is "),
			r("Helvetica-Bold", 11, "bold"),
			r("Helvetica", 11, " and "),
			r("Helvetica-Oblique", 11, "italic"),
			r("Helvetica", 11, " and 2*3."),
		)},
	}})
	if !strings.Contains(md, `This is **bold** and *italic* and 2\*3.`) {
		t.Errorf("got:\n%s", md)
	}
}

func TestEscapeParagraphStart(t *testing.T) {
	md := mustConvert(t, page{842, []block{{body(72, 100, "#1 rule: be kind.")}}})
	if !strings.HasPrefix(md, `\#1 rule`) {
		t.Errorf("got:\n%s", md)
	}
}

func TestEmpty(t *testing.T) {
	_, err := convertXML(t, stextXML(page{842, nil}, page{842, nil}))
	if !errors.Is(err, ErrNoText) {
		t.Errorf("err = %v", err)
	}
}

func TestOutputEndsWithOneNewline(t *testing.T) {
	md := mustConvert(t, page{842, []block{{body(72, 100, "Hello.")}}})
	if md != "Hello.\n" {
		t.Errorf("got %q", md)
	}
}

func TestParseStextFontFlags(t *testing.T) {
	lines, err := ParseStext(strings.NewReader(stextXML(page{842, []block{
		{l(72, 100, r("ABCDEF+Arial-BoldItalicMT", 11, "Hi"))},
	}})))
	if err != nil {
		t.Fatal(err)
	}
	s := lines[0].Spans[0]
	if !s.Bold || !s.Italic || s.Text != "Hi" {
		t.Errorf("span = %+v", s)
	}
	if lines[0].Size != 11 || lines[0].PageH != 842 || lines[0].Y != 100 {
		t.Errorf("line = %+v", lines[0])
	}
}
