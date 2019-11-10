package main

import (
	"strings"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

// Page represents a visited URL's contents; including
// the raw content, wrapped content, link slice, URL,
// and the current scroll position
type Page struct {
	WrappedContent []string
	RawContent     string
	Links          []string
	Location       Url
	ScrollPosition int
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (p *Page) ScrollPositionRange(termHeight int) (int, int) {
	termHeight -= 3
	if len(p.WrappedContent)-p.ScrollPosition < termHeight {
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

// WrapContent performs a hard wrap to the requested
// width and updates the WrappedContent
// of the Page struct width a string slice
// of the wrapped data
func (p *Page) WrapContent(width int) {
	counter := 0
	var content strings.Builder
	content.Grow(len(p.RawContent))
	for _, ch := range []rune(p.RawContent) {
		if ch == '\n' {
			content.WriteRune(ch)
			counter = 0
		} else if ch == '\t' {
			if counter+4 < width {
				content.WriteString("    ")
				counter += 4
			} else {
				content.WriteRune('\n')
				counter = 0
			}
		} else if ch == '\r' || ch == '\v' || ch == '\b' || ch == '\f' || ch == 27 {
			// Get rid of control characters we dont want
			continue
		} else {
			if counter < width {
				content.WriteRune(ch)
				counter++
			} else {
				content.WriteRune('\n')
				counter = 0
				if p.Location.Mime == "1" {
					spacer := "           "
					content.WriteString(spacer)
					counter += len(spacer)
				}
				content.WriteRune(ch)
			}
		}
	}

	p.WrappedContent = strings.Split(content.String(), "\n")
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakePage(url Url, content string, links []string) Page {
	p := Page{make([]string, 0), content, links, url, 0}
	return p
}
