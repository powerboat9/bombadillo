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

func Visit(url string, width int) (page, error) {
	if width > 80 {
		width = 80
	}
	w := fmt.Sprintf("-width=%d", width)
	c, err := exec.Command("lynx", "-dump", w, url).Output()
	if err != nil {
		return page{}, err
	}
	return parseLinks(string(c)), nil
}

// Returns false on err or non-text type
// Else returns true
func IsTextFile(url string) bool {
	c, err := exec.Command("lynx", "-dump", "-head", url).Output()
	if err != nil {
		return false
	}
	content := string(c)
	content = strings.ToLower(content)
	headers := strings.Split(content, "\n")
	for _, header := range headers {
		if strings.Contains(header, "content-type:") && strings.Contains(header, "text") {
			return true
		} else if strings.Contains(header, "content-type:") {
			return false
		}
	}

	// If we made it here, there is no content-type header.
	// So in the event of the unknown, lets render to the 
	// screen. This will allow redirects to get rendered
	// as well.
	return true
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

