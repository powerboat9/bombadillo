package http

import (
	"fmt"
	"os/exec"
	"strings"
)

type page struct {
	Content   string
	Links			[]string
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
