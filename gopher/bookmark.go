package gopher

import (
	"fmt"
	"strings"
)

type Bookmarks struct {
	Titles	[]string
	Links		[]string
}

func (b *Bookmarks) Add(v []string) error {
	if len(v) < 2 {
		return fmt.Errorf("Received %d arguments, expected 2 or more", len(v))
	}
	b.Titles = append(b.Titles, strings.Join(v[1:], " "))
	b.Links = append(b.Links, v[0])
	return nil
}

func (b *Bookmarks) Del(i int) error {
	if i < len(b.Titles) && i < len(b.Links) {
		b.Titles = append(b.Titles[:i], b.Titles[i + 1:]...)
		b.Links = append(b.Links[:i], b.Links[i + 1:]...)
		return nil
	}
	return fmt.Errorf("Bookmark %d does not exist", i)
}

func (b Bookmarks) List() []string {
	var out []string
	for i, t := range b.Titles {
		out = append(out, fmt.Sprintf("[%d] %s", i, t))
	}
	return out
}


func (b Bookmarks) IniDump() string {
	if len(b.Titles) < 0 {
		return ""
	}
	out := "[BOOKMARKS]\n"
	for i := 0; i < len(b.Titles); i++ {
		out += b.Titles[i]
		out += "="
		out += b.Links[i]
		out += "\n"
	}
	return out
}


func MakeBookmarks() Bookmarks {
	return Bookmarks{[]string{}, []string{}}
}
