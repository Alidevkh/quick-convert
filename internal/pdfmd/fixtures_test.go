package pdfmd

import (
	"fmt"
	"html"
	"strings"
)

// These helpers write small documents in the same XML format that
// "mutool draw -F stext" produces, so the tests read like a description of
// the page instead of hundreds of <char> lines.

type run struct {
	font string
	size float64
	text string
}

type line struct {
	x, y float64
	runs []run
}

type block []line

type page struct {
	height float64
	blocks []block
}

// r is a short way to write a run of text in one font.
func r(font string, size float64, text string) run { return run{font, size, text} }

// l is one line at position x, baseline y, made of runs.
func l(x, y float64, runs ...run) line { return line{x, y, runs} }

// body is a line of ordinary 11pt Helvetica text.
func body(x, y float64, text string) line { return l(x, y, r("Helvetica", 11, text)) }

func stextXML(pages ...page) string {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\"?>\n<document name=\"test.pdf\">\n")
	for i, p := range pages {
		fmt.Fprintf(&b, "<page id=\"page%d\" width=\"595\" height=\"%g\">\n", i+1, p.height)
		for _, blk := range p.blocks {
			b.WriteString("<block>\n")
			for _, ln := range blk {
				b.WriteString("<line wmode=\"0\" dir=\"1 0\">\n")
				x := ln.x
				for _, rn := range ln.runs {
					fmt.Fprintf(&b, "<font name=\"%s\" size=\"%g\">\n", rn.font, rn.size)
					for _, c := range rn.text {
						w := rn.size * 0.5
						top := ln.y - rn.size
						fmt.Fprintf(&b, "<char quad=\"%g %g %g %g %g %g %g %g\" x=\"%g\" y=\"%g\" color=\"#000000\" c=\"%s\"/>\n",
							x, top, x+w, top, x, ln.y, x+w, ln.y, x, ln.y, html.EscapeString(string(c)))
						x += w
					}
					b.WriteString("</font>\n")
				}
				b.WriteString("</line>\n")
			}
			b.WriteString("</block>\n")
		}
		b.WriteString("</page>\n")
	}
	b.WriteString("</document>\n")
	return b.String()
}
