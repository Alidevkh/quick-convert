// Package pdfwriter writes a simple PDF where every page is one JPEG
// picture, sized exactly to that picture.
//
// A PDF is a list of numbered "objects" followed by a table (the "xref")
// that says where each object starts in the file. For our PDFs the objects
// are:
//
//	1          the catalog (the document's starting point)
//	2          the page list
//	then, for every page:
//	  page     the page itself: its size and what it uses
//	  content  the drawing instruction: "draw the picture over the whole page"
//	  image    the JPEG bytes, stored exactly as they are (no re-compression)
package pdfwriter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// PixelsPerInch decides the printed size of a page: a 2480-pixel-wide
// photo becomes 2480/300 = 8.27 inches, the width of an A4 sheet.
const PixelsPerInch = 300

// Page is one picture to put on its own page.
type Page struct {
	JPEG   []byte // the JPEG file's bytes
	Width  int    // in pixels
	Height int    // in pixels
}

// Write writes a PDF containing the pages, in order.
func Write(w io.Writer, pages []Page) error {
	if len(pages) == 0 {
		return errors.New("no pages to write")
	}
	out := &counter{w: bufio.NewWriter(w)}
	var offsets []int // byte position of each object, for the xref table

	object := func(body string, stream []byte) {
		offsets = append(offsets, out.n)
		fmt.Fprintf(out, "%d 0 obj\n%s\n", len(offsets), body)
		if stream != nil {
			out.WriteString("stream\n")
			out.Write(stream)
			out.WriteString("\nendstream\n")
		}
		out.WriteString("endobj\n")
	}

	out.WriteString("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n") // header; the odd bytes mark the file as binary

	// Objects 1 and 2: catalog and page list. Page i uses objects 3+3i, 4+3i, 5+3i.
	var kids []string
	for i := range pages {
		kids = append(kids, fmt.Sprintf("%d 0 R", 3+3*i))
	}
	object("<< /Type /Catalog /Pages 2 0 R >>", nil)
	object(fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), len(pages)), nil)

	for i, p := range pages {
		contentObj, imageObj := 4+3*i, 5+3*i
		w, h := points(p.Width), points(p.Height)

		object(fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %s %s] "+
			"/Resources << /XObject << /Im0 %d 0 R >> >> /Contents %d 0 R >>",
			w, h, imageObj, contentObj), nil)

		// "q ... cm /Im0 Do Q" = scale the picture to the page size and draw it.
		draw := []byte(fmt.Sprintf("q %s 0 0 %s 0 0 cm /Im0 Do Q", w, h))
		object(fmt.Sprintf("<< /Length %d >>", len(draw)), draw)

		object(fmt.Sprintf("<< /Type /XObject /Subtype /Image /Width %d /Height %d "+
			"/ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode /Length %d >>",
			p.Width, p.Height, len(p.JPEG)), p.JPEG)
	}

	// The xref table: where every object starts.
	xref := out.n
	fmt.Fprintf(out, "xref\n0 %d\n0000000000 65535 f \n", len(offsets)+1)
	for _, off := range offsets {
		fmt.Fprintf(out, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(out, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(offsets)+1, xref)
	return out.w.Flush()
}

// points converts pixels to PDF points (1/72 inch) at PixelsPerInch,
// written with at most 2 decimals: 600 px → "144", 2480 px → "595.2".
func points(pixels int) string {
	v := math.Round(float64(pixels)*72/PixelsPerInch*100) / 100
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// counter writes through to w and counts the bytes written, so we know
// where each object starts.
type counter struct {
	w *bufio.Writer
	n int
}

func (c *counter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n += n
	return n, err
}

func (c *counter) WriteString(s string) {
	c.Write([]byte(s))
}
