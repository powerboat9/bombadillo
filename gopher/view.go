package gopher

import (
	"fmt"
	"strings"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

// View is a struct representing a gopher page. It contains
// the page content as a string slice, a list of link URLs
// as string slices, and the Url struct representing the page.
type View struct {
	Content []string
	Links   []string
	Address Url
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

// ParseMap is called by a view struct to parse a gophermap.
// It checks if the view is for a gophermap. If not,it does
// nothing. If so, it parses the gophermap into comment lines
// and link lines. For link lines it adds a link to the links
// slice and changes the content value to just the printable
// string plus a gophertype indicator and a link number that
// relates to the link position in the links slice. This
// receiver does not return anything.
func (v *View) ParseMap() {
	if v.Address.Gophertype == "1" || v.Address.Gophertype == "7" {
		for i, e := range v.Content {
			e = strings.Trim(e, "\r\n")
			if e == "." {
				v.Content[i] = " "
				continue
			}

			line := strings.Split(e, "\t")
			var title string
			if len(line[0]) > 1 {
				title = line[0][1:]
			} else {
				title = ""
			}
			if len(line[0]) > 0 && string(line[0][0]) == "i" {
				v.Content[i] = "           " + string(title)
				continue
			} else if len(line) >= 4 {
				fulllink := fmt.Sprintf("%s:%s/%s%s", line[2], line[3], string(line[0][0]), line[1])
				v.Links = append(v.Links, fulllink)
				linktext := fmt.Sprintf("(%s) %2d   %s", getType(string(line[0][0])), len(v.Links), title)
				v.Content[i] = linktext
			}
		}
	}
}

// Display is called on a view struct to print the contents of the view.
// This receiver does not return anything.
func (v View) Display() {
	fmt.Println()
	for _, el := range v.Content {
		fmt.Println(el)
	}
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

// MakeView creates and returns a new View struct from
// a Url and a string splice of content. This is used to
// initialize a View with a Url struct, links, and content.
// It takes a Url struct and a content []string and returns
// a View (NOT a pointer to a View).
func MakeView(url Url, content []string) View {
	v := View{content, make([]string, 0), url}
	v.ParseMap()
	return v
}
