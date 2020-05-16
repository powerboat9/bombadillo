package main

import (
	"fmt"
	"strings"

	"tildegit.org/sloum/bombadillo/tdiv"
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
	FoundLinkLines []int
	SearchTerm     string
	SearchIndex    int
	FileType       string
	WrapWidth      int
	Color          bool
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

// ScrollPositionRange may not be in actual usage....
// TODO: find where this is being used
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

func (p *Page) RenderImage(width int) {
	w := (width - 5) * 2
	if w > 300 {
		w = 300
	}
	p.WrappedContent = tdiv.Render([]byte(p.RawContent), w)
	p.WrapWidth = width
}

// WrapContent performs a hard wrap to the requested
// width and updates the WrappedContent
// of the Page struct width a string slice
// of the wrapped data
func (p *Page) WrapContent(width int, color bool) {
	width = min(width, 100)
	if p.FileType == "image" {
		p.RenderImage(width)
		return
	}
	counter := 0
	spacer := "           "
	var content strings.Builder
	var esc strings.Builder
	escape := false
	content.Grow(len(p.RawContent))
	for _, ch := range []rune(p.RawContent) {
		if escape {
			if color {
				esc.WriteRune(ch)
			}
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				escape = false
				if ch == 'm' {
					content.WriteString(esc.String())
					esc.Reset()
				}
			}
			continue
		}
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
		} else if ch == '\r' || ch == '\v' || ch == '\b' || ch == '\f' || ch == '\a' {
			// Get rid of control characters we dont want
			continue
		} else if ch == 27 {
			if p.Location.Scheme == "local" {
				if counter+4 >= width {
					content.WriteRune('\n')
				}
				content.WriteString("\\033")
				continue
			}
			escape = true
			if color {
				esc.WriteRune(ch)
			}
			continue
		} else {
			if counter <= width {
				content.WriteRune(ch)
				counter++
			} else {
				content.WriteRune('\n')
				counter = 0
				if p.Location.Mime == "1" {
					content.WriteString(spacer)
					counter += len(spacer)
				}
				content.WriteRune(ch)
			}
		}
	}

	p.WrappedContent = strings.Split(content.String(), "\n")
	p.WrapWidth = width
	p.Color = color
	p.HighlightFoundText()
}

func (p *Page) HighlightFoundText() {
	if p.SearchTerm == "" {
		return
	}
	for i, ln := range p.WrappedContent {
		found := strings.Index(ln, p.SearchTerm)
		if found < 0 {
			continue
		}
		format := "\033[7m%s\033[27m"
		if bombadillo.Options["theme"] == "inverse" {
			format = "\033[27m%s\033[7m"
		}
		ln = strings.Replace(ln, p.SearchTerm, fmt.Sprintf(format, p.SearchTerm), -1)
		p.WrappedContent[i] = ln
	}
}

func (p *Page) FindText() {
	p.FoundLinkLines = make([]int, 0, 10)
	s := p.SearchTerm
	p.SearchIndex = 0
	if s == "" {
		return
	}
	format := "\033[7m%s\033[27m"
	if bombadillo.Options["theme"] == "inverse" {
		format = "\033[27m%s\033[7m"
	}
	for i, ln := range p.WrappedContent {
		found := strings.Index(ln, s)
		if found < 0 {
			continue
		}
		ln = strings.Replace(ln, s, fmt.Sprintf(format, s), -1)
		p.WrappedContent[i] = ln
		p.FoundLinkLines = append(p.FoundLinkLines, i)
	}
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

// MakePage returns a Page struct with default values
func MakePage(url Url, content string, links []string) Page {
	p := Page{make([]string, 0), content, links, url, 0, make([]int, 0), "", 0, "", 40, false}
	return p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
