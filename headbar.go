package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Headbar struct {
	title string
	url string
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (h *Headbar) Render(width int, theme string) string {
	maxMsgWidth := width - len([]rune(h.title)) - 2
	if theme == "inverse" {
		return fmt.Sprintf("\033[7m%s▟\033[27m %-*.*s\033[0m", h.title, maxMsgWidth, maxMsgWidth, h.url)
	} else {
		return fmt.Sprintf("%s▟\033[7m %-*.*s\033[0m", h.title, maxMsgWidth, maxMsgWidth, h.url)
	}
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeHeadbar(title string) Headbar {
	return Headbar{title, ""}
}

