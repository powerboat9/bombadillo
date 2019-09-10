package main


//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Page struct {
	WrappedContent []string
	RawContent string
	Links []string
	Location Url
	ScrollPosition int
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\



//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakePage(url Url, content string) Page {
	p := Page{make([]string, 0), content, make([]string, 0), url, 0}
	return p
}

