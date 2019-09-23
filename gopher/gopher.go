// Contains the building blocks of a gopher client: history, url, and view.
// History handles the browsing session and view represents individual
// text based resources, the url represents a parsed url.
package gopher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

//------------------------------------------------\\
// + + +          V A R I A B L E S          + + + \\
//--------------------------------------------------\\

// types is a map of gophertypes to a string representing their
// type, to be used when displaying gophermaps
var types = map[string]string{
	"0": "TXT",
	"1": "MAP",
	"3": "ERR",
	"4": "BIN",
	"5": "DOS",
	"6": "UUE",
	"7": "FTS",
	"8": "TEL",
	"9": "BIN",
	"g": "GIF",
	"G": "GEM",
	"h": "HTM",
	"I": "IMG",
	"p": "PNG",
	"s": "SND",
	"S": "SSH",
	"T": "TEL",
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

// Retrieve makes a request to a Url and resturns
// the response as []byte/error. This function is
// available to use directly, but in most implementations
// using the "Visit" receiver of the History struct will
// be better.
func Retrieve(host, port, resource string) ([]byte, error) {
	nullRes := make([]byte, 0)
	timeOut := time.Duration(5) * time.Second

	if host == "" || port == "" {
		return nullRes, errors.New("Incomplete request url")
	}

	addr := host + ":" + port

	conn, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		return nullRes, err
	}

	send := resource + "\n"

	_, err = conn.Write([]byte(send))
	if err != nil {
		return nullRes, err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nullRes, err
	}

	return result, nil
}

// Visit handles the making of the request, parsing of maps, and returning
// the correct information to the client
func Visit(gophertype, host, port, resource string) (string, []string, error) {
	resp, err := Retrieve(host, port, resource)
	if err != nil {
		return "", []string{}, err
	} 
	
	text := string(resp)
	links := []string{}

	if IsDownloadOnly(gophertype) {
		return text, []string{}, nil
	}


	if gophertype == "1" {
		text, links = parseMap(text)
	}

	return text, links, nil
}

func getType(t string) string {
	if val, ok := types[t]; ok {
		return val
	}

	return "???"
}

func isWebLink(resource string) (string, bool) {
	split := strings.SplitN(resource, ":", 2)
	if first := strings.ToUpper(split[0]); first == "URL" && len(split) > 1 {
		return split[1], true
	}
	return "", false
}

func parseMap(text string) (string, []string)  {
	splitContent := strings.Split(text, "\n")
	links := make([]string, 0, 10)

	for i, e := range splitContent {
		e = strings.Trim(e, "\r\n")
		if e == "." {
			splitContent[i] = ""
			continue
		}

		line := strings.Split(e, "\t")
		var title string

		if len(line[0]) > 1 {
			title = line[0][1:]
		} else {
			title = ""
		}
		
		if len(line) > 1 && len(line[0]) > 0 && string(line[0][0]) == "i" {
			splitContent[i] = "           " + string(title)
		} else if len(line) >= 4 {
			link := buildLink(line[2], line[3], string(line[0][0]), line[1])
			links = append(links, link)
			linktext := fmt.Sprintf("(%s) %2d   %s", getType(string(line[0][0])), len(links), title)
			splitContent[i] = linktext
		}
	}
	return strings.Join(splitContent, "\n"), links
}

// Returns false for all text formats (including html
// even though it may link out. Things like telnet
// should never make it into the retrieve call for
// this module, having been handled in the client
// based on their protocol.
func IsDownloadOnly(gophertype string) bool {
	switch gophertype {
	case "0", "1", "3", "7", "h":
		return false
	default:
		return true
	}
}

func buildLink(host, port, gtype, resource string) string {
	switch gtype {
	case "8", "T":
		return fmt.Sprintf("telnet://%s:%s", host, port)
	case "G":
		return fmt.Sprintf("gemini://%s:%s%s", host, port, resource)
	case "h":
		u, tf := isWebLink(resource)
		if tf {
			if strings.Index(u, "://") > 0 {
				return u
			} else {
				return fmt.Sprintf("http://%s", u)
			}
		}
		return fmt.Sprintf("gopher://%s:%s/h%s", host, port, resource)
	default:
		return fmt.Sprintf("gopher://%s:%s/%s%s", host, port, gtype, resource)
	}
}
