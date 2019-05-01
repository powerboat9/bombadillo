package cui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// screenInit records whether or not the screen has been initialized
// this is used to prevent more than one screen from being used
var screenInit bool = false

// Screen represent the top level abstraction for a cui application.
// It takes up the full width and height of the terminal window and
// holds the various Windows and MsgBars for the application as well
// as a record of which window is active for control purposes.
type Screen struct {
	Height       int
	Width        int
	Windows      []*Window
	Activewindow int
	Bars         []*MsgBar
}

// AddWindow adds a new window to the Screen struct in question
func (s *Screen) AddWindow(r1, c1, r2, c2 int, scroll, border, show bool) {
	w := Window{box{r1, c1, r2, c2}, scroll, 0, []string{}, border, false, show}
	s.Windows = append(s.Windows, &w)
}

// AddMsgBar adds a new MsgBar to the Screen struct in question
func (s *Screen) AddMsgBar(row int, title, msg string, showTitle bool) {
	b := MsgBar{row, title, msg, showTitle}
	s.Bars = append(s.Bars, &b)
}

// DrawAllWindows loops over every window in the Screen struct and
// draws it to screen in index order (smallest to largest)
func (s Screen) DrawAllWindows() {
	s.Clear()
	for _, w := range s.Windows {
		if w.Show {
			w.DrawWindow()
		}
	}
	MoveCursorTo(s.Height-1, 1)
}

// Clear removes all content from the interior of the screen
func (s Screen) Clear() {
	for i := 0; i <= s.Height; i++ {
		MoveCursorTo(i, 0)
		Clear("line")
	}
}

// ReflashScreen checks for a screen resize and resizes windows if
// needed then redraws the screen. It takes a bool to decide whether
// to redraw the full screen or just the content. On a resize
// event, the full screen will always be redrawn.
func (s *Screen) ReflashScreen(clearScreen bool) {
	oldh, oldw := s.Height, s.Width
	s.GetSize()
	if s.Height != oldh || s.Width != oldw {
		// TODO this should be pure library code and not rely on
		// specific windows being present with specific behaviors.
		// Maybe allow windows to have a resize function that can
		// be declared within the application?
		// For now this will be ok though.
		s.Windows[0].Box.row2 = s.Height - 2
		s.Windows[0].Box.col2 = s.Width
		bookmarksWidth := 40
		if s.Width < 40 {
			bookmarksWidth = s.Width
		}
		s.Windows[1].Box.row2 = s.Height - 2
		s.Windows[1].Box.col1 = s.Width - bookmarksWidth
		s.Windows[1].Box.col2 = s.Width

		s.DrawAllWindows()
		s.DrawMsgBars()
	} else if clearScreen {
		s.DrawAllWindows()
		s.DrawMsgBars()
	} else {
		for _, w := range s.Windows {
			if w.Show {
				w.DrawWindow()
			}
		}
	}
}

// DrawMsgBars draws all MsgBars present in the Screen struct.
// All MsgBars are looped over and drawn in index order (sm - lg).
func (s *Screen) DrawMsgBars() {
	for _, bar := range s.Bars {
		MoveCursorTo(bar.row, 1)
		Clear("line")
		fmt.Print("\033[7m")
		fmt.Print(strings.Repeat(" ", s.Width))
		MoveCursorTo(bar.row, 1)
		var buf bytes.Buffer
		title := bar.title
		fmt.Print(title)
		if len(bar.title) > s.Width {
			title = string(bar.title[:s.Width-3]) + "..."
		}
		_, _ = buf.WriteString(title)
		msg := bar.message
		if len(bar.message) > s.Width-len(title) {
			msg = string(bar.message[:s.Width-len(title)-3]) + "..."
		}
		_, _ = buf.WriteString(msg)

		MoveCursorTo(bar.row, 1)
		fmt.Print(buf.String())
		fmt.Print("\033[0m")

	}
}

// GetSize retrieves the terminal size and sets the Screen
// width and height to that size
func (s *Screen) GetSize() {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Fatal error: Unable to retrieve terminal size")
		os.Exit(1)
	}
	var h, w int
	fmt.Sscan(string(out), &h, &w)
	s.Height = h
	s.Width = w
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -

// NewScreen is a constructor function that returns a pointer
// to a Screen struct
func NewScreen() (*Screen, error) {
	if screenInit {
		fmt.Println("Fatal error: Cannot create multiple screens")
		os.Exit(1)
	}
	var s Screen
	s.GetSize()
	for i := 0; i < s.Height; i++ {
		fmt.Println()
	}
	err := SetCharMode()
	if err != nil {
		return nil, err
	}

	Clear("screen")
	screenInit = true
	return &s, nil
}
