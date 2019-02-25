// Contains the building blocks of a gopher client: history and view.
// History handles the browsing session and view represents individual
// text based resources.
package gopher

import (
	"fmt"
	"strings"
	"errors"
	"regexp"
	"net"
	"io/ioutil"
	"time"
)


//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

// The history struct represents the history of the browsing
// session. It contains the current history position, the
// length of the active history space (this can be different
// from the available capacity in the Collection), and a
// collection array containing View structs representing
// each page in the current history. In general usage this
// struct should be initialized via the MakeHistory function.
type History struct {
	Position	int
	Length		int
	Collection [20]View
}

// The view struct represents a gopher page. It contains
// the page content as a string slice, a list of link URLs
// as string slices, and the Url struct representing the page.
type View struct {
	Content		[]string
	Links			[]string
	Address		Url
}

// The url struct represents a URL for the rest of the system.
// It includes component parts as well as a full URL string.
type Url struct {
	Scheme      string
	Host        string
	Port        string
	Gophertype  string
	Resource    string
	Full				string
	IsBinary		bool
}



//------------------------------------------------\\
// + + +          V A R I A B L E S          + + + \\
//--------------------------------------------------\\

// Types is a map of gophertypes to a string representing their
// type, to be used when displaying gophermaps
var Types = map[string]string{
	"0": "TXT",
	"1": "MAP",
	"h": "HTM",
	"3": "ERR",
	"4": "BIN",
	"5": "DOS",
	"s": "SND",
	"g": "GIF",
	"I": "IMG",
	"9": "BIN",
	"7": "FTS",
	"6": "UUE",
	"p": "PNG",
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

// The "Add" receiver takes a view and adds it to
// the history struct that called it. "Add" returns
// nothing. "Add" will shift history down if the max
// history length would be exceeded, and will reset
// history length if something is added in the middle.
func (h *History) Add(v View) {
	v.ParseMap()
	if h.Position == h.Length - 1 && h.Length < len(h.Collection) {
		h.Collection[h.Length] = v
		h.Length++
		h.Position++
	} else if h.Position == h.Length - 1 && h.Length == 20  {
		for x := 1; x < len(h.Collection); x++ {
			h.Collection[x-1] = h.Collection[x]
		}
		h.Collection[len(h.Collection)-1] = v
	} else {
		h.Position += 1
		h.Length = h.Position + 1
		h.Collection[h.Position] = v
	}
}

// The "Get" receiver is called by a history struct
// and returns a View from the current position, will
// return an error if history is empty and there is
// nothing to get.
func (h History) Get() (*View, error) {
	if h.Position < 0 {
		return nil, errors.New("History is empty, cannot get item from empty history.")
	}

	return &h.Collection[h.Position], nil
}

// The "GoBack" receiver is called by a history struct.
// When called it decrements the current position and
// displays the content for the View in that position.
// If history is at position 0, no action is taken.
func (h *History) GoBack() {
	if h.Position > 0 {
		h.Position--
		h.DisplayCurrentView()
	}
}


// The "GoForward" receiver is called by a history struct.
// When called it increments the current position and
// displays the content for the View in that position.
// If history is at position len - 1, no action is taken.
func (h *History) GoForward() {
	if h.Position + 1 < h.Length {
		h.Position++
		h.DisplayCurrentView()
	}
}

// The "DisplayCurrentView" receiver is called by a history
// struct. It calls the Display receiver for th view struct
// at the current history position. "DisplayCurrentView" does
// not return anything, and does nothing if position is less
// that 0.
func (h *History) DisplayCurrentView() {
	h.Collection[h.Position].Display()
}

// The "Visit" receiver is a high level combination of a few
// different receivers that makes it easy to create a Url,
// make a request to that Url, and add the response and Url
// to a View. That View then gets added to the History struct
// that the Visit receiver was called on. Returns a boolean
// value indicating whether or not the content is binary or
// textual data.
func (h *History) Visit(addr string) (View, error) {
	u, err := MakeUrl(addr)
	if err != nil {
		return View{}, err
	}

	text, err := Retrieve(u)
	if err != nil {
		return View{}, err
	}

	var pageContent []string
	if u.IsBinary {
		pageContent = []string{string(text)}
	} else {
		pageContent = strings.Split(string(text), "\n")
	}

	return MakeView(u, pageContent), nil
}

// The "ParseMap" receiver is called by a view struct. It
// checks if the view is for a gophermap. If not,it does
// nothing. If so, it parses the gophermap into comment lines
// and link lines. For link lines it adds a link to the links
// slice and changes the content value to just the printable
// string plus a gophertype indicator and a link number that
// relates to the link position in the links slice. This
// receiver does not return anything.
func (v *View) ParseMap() {
	if v.Address.Gophertype == "1" {
		for i, e := range v.Content {
			e = strings.Trim(e, "\r\n")
			line := strings.Split(e,"\t")
			if len(line[0]) > 0 && string(line[0][0]) == "i" {
				v.Content[i] = "           " + string(line[0][1:])
				continue
			} else if len(line) >= 4 {
				fulllink := fmt.Sprintf("%s:%s/%s%s", line[2], line[3], string(line[0][0]), line[1])
				v.Links = append(v.Links, fulllink)
				linktext := fmt.Sprintf("(%s) %2d   %s", Types[string(line[0][0])], len(v.Links), string(line[0][1:]))
				v.Content[i] = linktext
			}
		}
	}	
}

// The "Display" receiver is called on a view struct.
// It prints the content, line by line, of the View.
// This receiver does not return anything.
func (v View) Display() {
	fmt.Println()
	for _, el := range v.Content {
		if el != "." {
			fmt.Println(el)
		}
	}
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

// MakeUrl is a Url constructor that takes in a string 
// representation of a url and returns a Url struct and 
// an error (or nil).
func MakeUrl(u string) (Url, error) {
	var out Url
	re := regexp.MustCompile(`^((?P<scheme>gopher|http|https|ftp|telnet):\/\/)?(?P<host>[\w\-\.\d]+)(?::(?P<port>\d+)?)?(?:/(?P<type>[01345679gIhisp])?)?(?P<resource>(?:\/.*)?)?$`)
	match := re.FindStringSubmatch(u)

	if valid := re.MatchString(u); valid != true {
		return out, errors.New("Invalid URL or command character")
	}

	for i, name := range re.SubexpNames() {
		switch name {
			case "scheme":
				out.Scheme = match[i]
			case "host":
				out.Host = match[i]
			case "port":
				out.Port = match[i]
			case "type":
				out.Gophertype = match[i]
			case "resource":
				out.Resource = match[i]
		}
	}

	if out.Scheme == "" {
		out.Scheme = "gopher"
	}

	if out.Host == "" {
		return out, errors.New("No host.")
	}

	if out.Scheme == "gopher" && out.Port == "" {
		out.Port = "70"
	} else if out.Scheme == "http" || out.Scheme == "https" && out.Port == "" {
		out.Port = "80"
	}

	if out.Gophertype == "" && (out.Resource == "" || out.Resource == "/") {
		out.Gophertype = "1"
	} 

	if out.Gophertype == "1" || out.Gophertype == "0" {
		out.IsBinary = false
	} else {
		out.IsBinary = true
	}

	if out.Scheme == "gopher" && out.Gophertype == "" {
		out.Gophertype = "0"
	}

	out.Full = out.Scheme + "://" + out.Host + ":" + out.Port + "/" + out.Gophertype + out.Resource

	return out, nil
}


// Constructor function for History struct.
// This is used to initialize history position
// as -1, which is needed. Returns a copy of
// initialized History struct (does NOT return
// a pointer to the struct).
func MakeHistory() History {
	return History{-1, 0, [20]View{}}
}


// Constructor function for View struct.
// This is used to initialize a View with
// a Url struct, links, and content. It takes
// a Url struct and a content []string and returns
// a View (NOT a pointer to a View).
func MakeView(url Url, content []string) View {
	v := View{content, make([]string, 0), url}
	v.ParseMap()
	return v
}


// Retrieve makes a request to a Url and resturns
// the response as []byte/error. This function is
// available to use directly, but in most implementations
// using the "Visit" receiver of the History struct will
// be better.
func Retrieve(u Url) ([]byte, error) {
  nullRes := make([]byte, 0)
	timeOut := time.Duration(5) * time.Second

  if u.Host == "" || u.Port == "" {
		return nullRes, errors.New("Incomplete request url")
  }

	addr := u.Host + ":" + u.Port

	conn, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		return nullRes, err
	}

	send := u.Resource + "\n"
	if u.Scheme == "http" || u.Scheme == "https" {
		send = u.Gophertype
	}

	_, err = conn.Write([]byte(send))
	if err != nil {
		return nullRes, err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nullRes, err
	}

	return result, err
}


