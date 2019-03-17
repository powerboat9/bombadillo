package cui

import (
	"fmt"
	"strings"
)


type box struct {
	row1			int
	col1			int
	row2			int
	col2			int
}

// TODO add coloring
type Window struct {
	Box							box
	Scrollbar				bool
	Scrollposition	int
	Content					[]string
	drawBox					bool
	Active					bool
	Show						bool
}

func (w *Window) DrawWindow() {
	w.DrawContent()

	if w.drawBox {
		w.DrawBox()
	}
}

func (w *Window) DrawBox(){
	moveThenDrawShape(w.Box.row1, w.Box.col1, "tl")
	moveThenDrawShape(w.Box.row1, w.Box.col2, "tr")
	moveThenDrawShape(w.Box.row2, w.Box.col1, "bl")
	moveThenDrawShape(w.Box.row2, w.Box.col2, "br")
	for i := w.Box.col1 + 1; i < w.Box.col2; i++ {
		moveThenDrawShape(w.Box.row1, i, "ceiling")
		moveThenDrawShape(w.Box.row2, i, "ceiling")
	}

	for i:= w.Box.row1 + 1; i < w.Box.row2; i++ {
		moveThenDrawShape(i, w.Box.col1, "wall")
		moveThenDrawShape(i, w.Box.col2, "wall")
	}
}

func (w *Window) DrawContent(){
	var maxlines, border_thickness, contenth int
	var short_content bool = false

	if w.drawBox {
		border_thickness, contenth = -1, 1
	} else {
		border_thickness, contenth = 1, 0
	}

	height := w.Box.row2 - w.Box.row1 + border_thickness
	width := w.Box.col2 - w.Box.col1 + border_thickness

	content := WrapLines(w.Content, width)

	if len(content) < w.Scrollposition + height {
		maxlines = len(content)
		short_content = true
	} else {
		maxlines = w.Scrollposition + height
	}

	for i := w.Scrollposition; i < maxlines; i++ {
		MoveCursorTo(w.Box.row1 + contenth + i - w.Scrollposition, w.Box.col1 + contenth)
		fmt.Print( strings.Repeat(" ", width) )
		MoveCursorTo(w.Box.row1 + contenth + i - w.Scrollposition, w.Box.col1 + contenth)
		fmt.Print(content[i])
	}
	if short_content {
		for i := len(content); i <= height; i++ {
			MoveCursorTo(w.Box.row1 + contenth + i - w.Scrollposition, w.Box.col1 + contenth)
			fmt.Print( strings.Repeat(" ", width) )
		}
	}
}

func (w *Window) ScrollDown() {
	height := w.Box.row2 - w.Box.row1 - 1
	contentLength := len(w.Content)
	if w.Scrollposition < contentLength - height {
		w.Scrollposition++
	} else {
		fmt.Print("\a")
	}
}

func (w *Window) ScrollUp() {
	if w.Scrollposition > 0 {
		w.Scrollposition--
	} else {
		fmt.Print("\a")
	}
}
