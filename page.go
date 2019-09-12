package main

import (
	"strings"
	"bytes"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Page struct {
	WrappedContent []string
	RawContent string
	Links []string
	Location Url
	ScrollPosition int
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (p *Page) ScrollPositionRange(termHeight int) (int, int) {
	termHeight -= 3
	if len(p.WrappedContent) - p.ScrollPosition < termHeight {
		p.ScrollPosition = len(p.WrappedContent) - termHeight
	}
	if p.ScrollPosition < 0 {
		p.ScrollPosition = 0
	}
	var end int
	if len(p.WrappedContent) < termHeight {
		end = len(p.WrappedContent)
	} else {
		end = p.ScrollPosition + termHeight
	}

	return p.ScrollPosition, end
}

func (p *Page) WrapContent(width int) {
	// TODO this is a temporary wrapping function
	// in order to test. Rebuild it.
	src := strings.Split(p.RawContent, "\n")
	out := []string{}
	for _, ln := range src {
		if len([]rune(ln)) <= width {
			out = append(out, ln)
		} else {
			words := strings.SplitAfter(ln, " ")
			var subout bytes.Buffer
			for i, wd := range words {
				sublen := subout.Len()
				wdlen := len([]rune(wd))
				if sublen+wdlen <= width {
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
					}
				} else {
					out = append(out, subout.String())
					subout.Reset()
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
						subout.Reset()
					}
				}
			}
		}
	}
	p.WrappedContent = out
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakePage(url Url, content string, links []string) Page {
	p := Page{make([]string, 0), content, links, url, 0}
	return p
}

