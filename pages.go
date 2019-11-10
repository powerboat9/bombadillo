package main

import (
	"fmt"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Pages struct {
	Position int
	Length   int
	History  [20]Page
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (p *Pages) NavigateHistory(qty int) error {
	newPosition := p.Position + qty
	if newPosition < 0 {
		return fmt.Errorf("You are already at the beginning of history")
	} else if newPosition > p.Length-1 {
		return fmt.Errorf("Your way is blocked by void, there is nothing forward")
	}

	p.Position = newPosition
	return nil
}

func (p *Pages) Add(pg Page) {
	if p.Position == p.Length-1 && p.Length < len(p.History) {
		p.History[p.Length] = pg
		p.Length++
		p.Position++
	} else if p.Position == p.Length-1 && p.Length == 20 {
		for x := 1; x < len(p.History); x++ {
			p.History[x-1] = p.History[x]
		}
		p.History[len(p.History)-1] = pg
	} else {
		p.Position += 1
		p.Length = p.Position + 1
		p.History[p.Position] = pg
	}
}

func (p *Pages) Render(termHeight, termWidth int) []string {
	if p.Length < 1 {
		return make([]string, 0)
	}
	pos := p.History[p.Position].ScrollPosition
	prev := len(p.History[p.Position].WrappedContent)
	p.History[p.Position].WrapContent(termWidth)
	now := len(p.History[p.Position].WrappedContent)
	if prev > now {
		diff := prev - now
		pos = pos - diff
	} else if prev < now {
		diff := now - prev
		pos = pos + diff
		if pos > now-termHeight {
			pos = now - termHeight
		}
	}

	if pos < 0 || now < termHeight-3 {
		pos = 0
	}

	p.History[p.Position].ScrollPosition = pos

	return p.History[p.Position].WrappedContent[pos:]
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakePages() Pages {
	return Pages{-1, 0, [20]Page{}}
}
