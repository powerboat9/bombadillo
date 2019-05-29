package cui

import (
	"fmt"
	"strings"
)

type box struct {
	Row1 int
	Col1 int
	Row2 int
	Col2 int
}

// TODO add coloring
type Window struct {
	Box            box
	Scrollbar      bool
	Scrollposition int
	Content        []string
	drawBox        bool
	Active         bool
	Show           bool
  tempContentLen  int
}

func (w *Window) DrawWindow() {
	w.DrawContent()

	if w.drawBox {
		w.DrawBox()
	}
}

func (w *Window) DrawBox() {
	lead := ""
	if w.Active {
		lead = "a"
	}
	moveThenDrawShape(w.Box.Row1, w.Box.Col1, lead+"tl")
	moveThenDrawShape(w.Box.Row1, w.Box.Col2, lead+"tr")
	moveThenDrawShape(w.Box.Row2, w.Box.Col1, lead+"bl")
	moveThenDrawShape(w.Box.Row2, w.Box.Col2, lead+"br")
	for i := w.Box.Col1 + 1; i < w.Box.Col2; i++ {
		moveThenDrawShape(w.Box.Row1, i, lead+"ceiling")
		moveThenDrawShape(w.Box.Row2, i, lead+"ceiling")
	}

	for i := w.Box.Row1 + 1; i < w.Box.Row2; i++ {
		moveThenDrawShape(i, w.Box.Col1, lead+"wall")
		moveThenDrawShape(i, w.Box.Col2, lead+"wall")
	}
}

func (w *Window) DrawContent() {
	var maxlines, borderThickness, contenth int
	var short_content bool = false

	if w.drawBox {
		borderThickness, contenth = -1, 1
	} else {
		borderThickness, contenth = 1, 0
	}

	height := w.Box.Row2 - w.Box.Row1 + borderThickness
	width := w.Box.Col2 - w.Box.Col1 + borderThickness

	content := wrapLines(w.Content, width)
  w.tempContentLen = len(content)

	if w.Scrollposition > w.tempContentLen-height {
		w.Scrollposition = w.tempContentLen-height
		if w.Scrollposition < 0 {
			w.Scrollposition = 0
		}
	}

	if len(content) < w.Scrollposition+height {
		maxlines = len(content)
		short_content = true
	} else {
		maxlines = w.Scrollposition + height
	}

	for i := w.Scrollposition; i < maxlines; i++ {
		MoveCursorTo(w.Box.Row1+contenth+i-w.Scrollposition, w.Box.Col1+contenth)
		fmt.Print(strings.Repeat(" ", width))
		MoveCursorTo(w.Box.Row1+contenth+i-w.Scrollposition, w.Box.Col1+contenth)
		fmt.Print(content[i])
	}
	if short_content {
		for i := len(content); i <= height; i++ {
			MoveCursorTo(w.Box.Row1+contenth+i-w.Scrollposition, w.Box.Col1+contenth)
			fmt.Print(strings.Repeat(" ", width))
		}
	}
}

func (w *Window) ScrollDown() {
	var borderThickness int
	if w.drawBox {
		borderThickness = -1
	} else {
		borderThickness = 1
	}

	height := w.Box.Row2 - w.Box.Row1 + borderThickness

	if w.Scrollposition < w.tempContentLen-height {
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

func (w *Window) PageDown() {
	var borderThickness int
	if w.drawBox {
		borderThickness = -1
	} else {
		borderThickness = 1
	}

	height := w.Box.Row2 - w.Box.Row1 + borderThickness

	if w.Scrollposition < w.tempContentLen-height {
		w.Scrollposition += height
		if w.Scrollposition > w.tempContentLen-height {
			w.Scrollposition = w.tempContentLen-height
		}
	} else {
		fmt.Print("\a")
	}
}

func (w *Window) PageUp() {
	var borderThickness int
	if w.drawBox {
		borderThickness = -1
	} else {
		borderThickness = 1
	}

	height := w.Box.Row2 - w.Box.Row1 + borderThickness
	contentLength := len(w.Content)
	if w.Scrollposition > 0 && height < contentLength {
		w.Scrollposition -= height
		if w.Scrollposition < 0 {
			w.Scrollposition = 0
		}
	} else {
		fmt.Print("\a")
	}
}

func (w *Window) ScrollHome() {
	if w.Scrollposition > 0 {
		w.Scrollposition = 0
	} else {
		fmt.Print("\a")
	}
}

func (w *Window) ScrollEnd() {
	var borderThickness int
	if w.drawBox {
		borderThickness = -1
	} else {
		borderThickness = 1
	}

	height := w.Box.Row2 - w.Box.Row1 + borderThickness

	if w.Scrollposition < w.tempContentLen-height {
		w.Scrollposition = w.tempContentLen-height
	} else {
		fmt.Print("\a")
	}
}
