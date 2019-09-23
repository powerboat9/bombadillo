package main

import (
	"fmt"
	"strings"

	"tildegit.org/sloum/bombadillo/cui"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Bookmarks struct {
	IsOpen bool
	IsFocused bool
	Position int
	Length int
	Titles []string
	Links []string
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (b *Bookmarks) Add(v []string) (string, error) {
	if len(v) < 2 {
		return "", fmt.Errorf("Received %d arguments, expected 2+", len(v))
	}
	b.Titles = append(b.Titles, strings.Join(v[1:], " "))
	b.Links = append(b.Links, v[0])
	b.Length = len(b.Titles)
	return "Bookmark added successfully", nil
}

func (b *Bookmarks) Delete(i int) (string, error) {
	if i < len(b.Titles) && len(b.Titles) == len(b.Links) {
		b.Titles = append(b.Titles[:i], b.Titles[i+1:]...)
		b.Links = append(b.Links[:i], b.Links[i+1:]...)
		b.Length = len(b.Titles)
		return "Bookmark deleted successfully", nil
	}
	return "", fmt.Errorf("Bookmark %d does not exist", i)
}

func (b *Bookmarks) ToggleOpen() {
	b.IsOpen = !b.IsOpen
	if b.IsOpen {
		b.IsFocused = true
	} else {
		b.IsFocused = false
		cui.Clear("screen")
	}
}

func (b *Bookmarks) ToggleFocused() {
	if b.IsOpen {
		b.IsFocused = !b.IsFocused
	}
}

func (b Bookmarks) IniDump() string {
	if len(b.Titles) < 0 {
		return ""
	}
	out := "[BOOKMARKS]\n"
	for i := 0; i < len(b.Titles); i++ {
		out += b.Titles[i]
		out += "="
		out += b.Links[i]
		out += "\n"
	}
	return out
}

// Get a list, including link nums, of bookmarks
// as a string slice
func (b Bookmarks) List() []string {
	var out []string
	for i, t := range b.Titles {
		out = append(out, fmt.Sprintf("[%d] %s", i, t))
	}
	return out
}

func (b Bookmarks) Render(termwidth, termheight int) []string {
	width := 40
	termheight -= 3
	var walll, wallr, floor, ceil, tr, tl, br, bl string
	if termwidth < 40 {
		width = termwidth
	} 
	if b.IsFocused {
		walll = cui.Shapes["awalll"]
		wallr = cui.Shapes["awallr"]
		ceil = cui.Shapes["aceiling"]
		floor = cui.Shapes["afloor"]
		tr = cui.Shapes["atr"]
		br = cui.Shapes["abr"]
		tl = cui.Shapes["atl"]
		bl = cui.Shapes["abl"]
	} else {
		walll = cui.Shapes["walll"]
		wallr = cui.Shapes["wallr"]
		ceil = cui.Shapes["ceiling"]
		floor = cui.Shapes["floor"]
		tr = cui.Shapes["tr"]
		br = cui.Shapes["br"]
		tl = cui.Shapes["tl"]
		bl = cui.Shapes["bl"]
	}

	out := make([]string, 0, 5)
	contentWidth := width - 2
	top := fmt.Sprintf("%s%s%s", tl, strings.Repeat(ceil, contentWidth), tr)
	out = append(out, top)
	marks := b.List()
	for i := 0; i < termheight - 2; i++ {
		if i + b.Position >= len(b.Titles) {
			out = append(out, fmt.Sprintf("%s%-*.*s%s", walll, contentWidth, contentWidth, "", wallr))
		} else {
			out = append(out, fmt.Sprintf("%s%-*.*s%s", walll, contentWidth, contentWidth, marks[i + b.Position], wallr))
		}
	}

	bottom := fmt.Sprintf("%s%s%s", bl, strings.Repeat(floor, contentWidth), br)
	out = append(out, bottom)
	return out
}

// TODO handle scrolling of the bookmarks list
// either here with a scroll up/down or in the client
// code for scroll


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeBookmarks() Bookmarks {
	return Bookmarks{false, false, 0, 0, make([]string, 0), make([]string, 0)}
}

