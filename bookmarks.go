package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Bookmarks struct {
	IsOpen bool
	IsFocused bool
	Position int
	Length int
	Titles []string
	Links []string
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (b *Bookmarks) Add([]string) error {
	// TODO add a bookmark
	return fmt.Errorf("")
}

func (b *Bookmarks) Delete(int) error {
	// TODO delete a bookmark
	return fmt.Errorf("")
}

func (b *Bookmarks) ToggleOpen() {
	b.IsOpen = !b.IsOpen
	if b.IsOpen {
		b.IsFocused = true
	} else {
		b.IsFocused = false
	}
}

func (b *Bookmarks) ToggleFocused() {
	if b.IsOpen {
		b.IsFocused = !b.IsFocused
	}
}

func (b *Bookmarks) IniDump() string {
	// TODO create dump of values for INI file
	return ""
}

func (b *Bookmarks) Render() ([]string, error) {
	// TODO grab all of the bookmarks as a fixed
	// width string including border and spacing
	return []string{}, fmt.Errorf("")
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeBookmarks() Bookmarks {
	return Bookmarks{false, false, 0, 0, make([]string, 0), make([]string, 0)}
}

