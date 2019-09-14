package main

import (
	"fmt"
	"strconv"
)


//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Footbar struct {
	PercentRead string
	PageType string
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (f *Footbar) SetPercentRead(p int) {
	if p > 100 {
		p = 100
	} else if p < 0 {
		p = 0
	}
	f.PercentRead = strconv.Itoa(p) + "%"
}

func (f *Footbar) SetPageType(t string) {
	f.PageType = t
}

func (f *Footbar) Render(termWidth, position int, theme string) string {
	pre := fmt.Sprintf("HST: (%2.2d) - - - %4s Read ", position + 1, f.PercentRead)
	out := "\033[0m%*.*s "
	if theme == "inverse" {
		out = "\033[7m%*.*s \033[0m"
	}
	return fmt.Sprintf(out, termWidth - 1, termWidth - 1, pre)
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeFootbar() Footbar {
	return Footbar{"---", "N/A"}
}

