package main


//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Headbar struct {
	title string
	url string
	content string
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (h *Headbar) SetUrl(u string) {
	h.url = u
}

func (h *Headbar) Build(width string) string {
	// TODO Build out header to specified width
	h.content = "" // This is a temp value to show intention
	return ""
}

func (h *Headbar) Draw() {
	// TODO this will actually draw the bar
	// without having to redraw everything else
}

func (h *Headbar) Render() string {
	// TODO returns the content value
	return ""
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeHeadbar(title string) Headbar {
	return Headbar{title, "", title}
}

