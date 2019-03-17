// Contains the building blocks of a gopher client: history and view.
// History handles the browsing session and view represents individual
// text based resources.
package gopher

import (
	"strings"
	"errors"
	"net"
	"io/ioutil"
	"time"
)


//------------------------------------------------\\
// + + +          V A R I A B L E S          + + + \\
//--------------------------------------------------\\

// Types is a map of gophertypes to a string representing their
// type, to be used when displaying gophermaps
var types = map[string]string{
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
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\


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


// The "Visit" function is a high level combination of a few
// different types that makes it easy to create a Url,
// make a request to that Url, and add the response and Url
// to a View. Returns a copy of the view and an error (or nil).
func Visit(addr string) (View, error) {
	u, err := MakeUrl(addr)
	if err != nil {
		return View{}, err
	}

	text, err := Retrieve(u)
	if err != nil {
		return View{}, err
	}

	var pageContent []string
	if u.IsBinary && u.Gophertype != "7" {
		pageContent = []string{string(text)}
	} else {
		pageContent = strings.Split(string(text), "\n")
	}

	return MakeView(u, pageContent), nil
}

func GetType(t string) string {
	if val, ok := types[t]; ok {
		return val
	} 
	return "???"

}

