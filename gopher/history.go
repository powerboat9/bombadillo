package gopher

import (
	"fmt"
	"errors"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

// The history struct represents the history of the browsing
// session. It contains the current history position, the
// length of the active history space (this can be different
// from the available capacity in the Collection), and a
// collection array containing View structs representing
// each page in the current history. In general usage this
// struct should be initialized via the MakeHistory function.
type History struct {
	Position		int
	Length			int
	Collection	[20]View
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

// The "Add" receiver takes a view and adds it to
// the history struct that called it. "Add" returns
// nothing. "Add" will shift history down if the max
// history length would be exceeded, and will reset
// history length if something is added in the middle.
func (h *History) Add(v View) {
	v.ParseMap()
	if h.Position == h.Length - 1 && h.Length < len(h.Collection) {
		h.Collection[h.Length] = v
		h.Length++
		h.Position++
	} else if h.Position == h.Length - 1 && h.Length == 20  {
		for x := 1; x < len(h.Collection); x++ {
			h.Collection[x-1] = h.Collection[x]
		}
		h.Collection[len(h.Collection)-1] = v
	} else {
		h.Position += 1
		h.Length = h.Position + 1
		h.Collection[h.Position] = v
	}
}

// The "Get" receiver is called by a history struct
// and returns a View from the current position, will
// return an error if history is empty and there is
// nothing to get.
func (h History) Get() (*View, error) {
	if h.Position < 0 {
		return nil, errors.New("History is empty, cannot get item from empty history.")
	}

	return &h.Collection[h.Position], nil
}

// The "GoBack" receiver is called by a history struct.
// When called it decrements the current position and
// displays the content for the View in that position.
// If history is at position 0, no action is taken.
func (h *History) GoBack() bool {
	if h.Position > 0 {
		h.Position--
		return true
	}

	fmt.Print("\a")
	return false
}


// The "GoForward" receiver is called by a history struct.
// When called it increments the current position and
// displays the content for the View in that position.
// If history is at position len - 1, no action is taken.
func (h *History) GoForward() bool {
	if h.Position + 1 < h.Length {
		h.Position++
		return true
	}

	fmt.Print("\a")
	return false
}

// The "DisplayCurrentView" receiver is called by a history
// struct. It calls the Display receiver for th view struct
// at the current history position. "DisplayCurrentView" does
// not return anything, and does nothing if position is less
// that 0.
func (h *History) DisplayCurrentView() {
	h.Collection[h.Position].Display()
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\


// Constructor function for History struct.
// This is used to initialize history position
// as -1, which is needed. Returns a copy of
// initialized History struct (does NOT return
// a pointer to the struct).
func MakeHistory() History {
	return History{-1, 0, [20]View{}}
}
