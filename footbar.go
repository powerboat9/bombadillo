package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Footbar struct {
	PercentRead int
	PageType string
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (f *Footbar) SetPercentRead(p int) {
	f.PercentRead = p
}

func (f *Footbar) SetPageType(t string) {
	f.PageType = t
}

func (f *Footbar) Draw() {
	// TODO this will actually draw the bar
	// without having to redraw everything else
}

func (f *Footbar) Render(termWidth int) string {
	return fmt.Sprintf("\033[7m%-*.*s\033[0m", termWidth, termWidth, "")
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeFootbar() Footbar {
	return Footbar{100, "N/A"}
}

