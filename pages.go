package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Pages struct {
	Position int
	Length int
	History [20]Page
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (p *Pages) NavigateHistory(qty int) error {
	newPosition := p.Position + qty
	if newPosition < 0 {
		return fmt.Errorf("You are already at the beginning of history")
	} else if newPosition > p.Length - 1 {
		return fmt.Errorf("Your way is blocked by void, there is nothing forward")
	}

	p.Position = newPosition
	return nil
}

func (p *Pages) Add(pg Page) error {
	// TODO add the given page onto the pages struct
	// handling truncation of the history as needed.
	return fmt.Errorf("")
}

func (p *Pages) Render() ([]string, error) {
	// TODO grab the current page as wrappedContent
	// May need to handle spacing at end of lines.
	return []string{}, fmt.Errorf("")
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakePages() Pages {
	return Pages{-1, 0, [20]Page{}}
}


