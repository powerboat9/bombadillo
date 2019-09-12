package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Headbar struct {
	title string
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (h *Headbar) Build(width string) string {
	// TODO Build out header to specified width
	return ""
}

func (h *Headbar) Draw() {
	// TODO this will actually draw the bar
	// without having to redraw everything else
}

func (h *Headbar) Render(width int, message string) string {
	maxMsgWidth := width - len([]rune(h.title))
	return fmt.Sprintf("\033[7m%s%-*.*s\033[0m", h.title, maxMsgWidth, maxMsgWidth, message)
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeHeadbar(title string) Headbar {
	return Headbar{title}
}

