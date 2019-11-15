package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

// Page represents the contents and links or an http/https document
type Page struct {
	Content string
	Links   []string
}

// Visit is the main entry to viewing a web document in bombadillo.
// It takes a url, a terminal width, and which web backend the user
// currently has set. Visit returns a Page and an error
func Visit(webmode, url string, width int) (Page, error) {
	if width > 80 {
		width = 80
	}
	var w string
	switch webmode {
	case "lynx":
		w = "-width"
	case "w3m":
		w = "-cols"
	case "elinks":
		w = "-dump-width"
	default:
		return Page{}, fmt.Errorf("Invalid webmode setting")
	}
	c, err := exec.Command(webmode, "-dump", w, fmt.Sprintf("%d", width), url).Output()
	if err != nil {
		return Page{}, err
	}
	return parseLinks(string(c)), nil
}

// IsTextFile makes an http(s) head request to a given URL
// and determines if the content-type is text based. It then
// returns a bool
func IsTextFile(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	ctype := resp.Header.Get("content-type")
	if strings.Contains(ctype, "text") || ctype == "" {
		return true
	}

	return false
}

func parseLinks(c string) Page {
	var out Page
	contentUntil := strings.LastIndex(c, "References")
	if contentUntil >= 1 {
		out.Content = c[:contentUntil]
	} else {
		out.Content = c
		out.Links = make([]string, 0)
		return out
	}
	links := c[contentUntil+11:]
	links = strings.TrimSpace(links)
	linkSlice := strings.Split(links, "\n")
	out.Links = make([]string, 0, len(linkSlice))
	for _, link := range linkSlice {
		ls := strings.SplitN(link, ".", 2)
		if len(ls) < 2 {
			continue
		}
		out.Links = append(out.Links, strings.TrimSpace(ls[1]))
	}
	return out
}

// Fetch makes an http(s) request and returns the []bytes
// for the response and an error. Fetch is used for saving
// the source file of an http(s) document
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return bodyBytes, nil
}
