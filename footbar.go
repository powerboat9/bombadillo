package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Footbar struct {
	PercentRead string
	PageType string
	Content string
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (f *Footbar) SetPercentRead(p int) {
	f.PercentRead = fmt.Sprintf("%d%%", p)
}

func (f *Footbar) SetPageType(t string) {
	f.PageType = t
}

func (f *Footbar) Draw() {
	// TODO this will actually draw the bar
	// without having to redraw everything else
}

func (f *Footbar) Build(width string) string {
	// TODO Build out header to specified width
	f.Content = "" // This is a temp value to show intention
	return ""
}

func (f *Footbar) Render() string {
	// TODO returns a full line
	return ""
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeFootbar() Footbar {
	return Footbar{"", "N/A", ""}
}

