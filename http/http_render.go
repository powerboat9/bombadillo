package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type page struct {
	Content string
	Links   []string
}

func Visit(webmode, url string, width int) (page, error) {
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
		return page{}, fmt.Errorf("Invalid webmode setting")
	}
	c, err := exec.Command(webmode, "-dump", w, fmt.Sprintf("%d", width), url).Output()
	if err != nil {
		return page{}, err
	}
	return parseLinks(string(c)), nil
}

// Returns false on err or non-text type
// Else returns true
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

func parseLinks(c string) page {
	var out page
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
