package cui

import (
	"strings"
	"fmt"
	"os"
	"os/exec"
	"bytes"
)


var screenInit bool = false

type Screen struct {
	Height				int
	Width					int
	Windows				[]*Window
  Activewindow	int
	Bars					[]*MsgBar
}


func (s *Screen) AddWindow(r1, c1, r2, c2 int, scroll, border, show bool) {
	w := Window{box{r1, c1, r2, c2}, scroll, 0, []string{}, border, false, show}
	s.Windows = append(s.Windows, &w)
}

func (s *Screen) AddMsgBar(row int, title, msg string, showTitle bool) {
	b := MsgBar{row, title, msg, showTitle}
	s.Bars = append(s.Bars, &b)
}

func (s Screen) DrawAllWindows() {
	s.Clear()
	for _, w := range s.Windows {
		if w.Show {
			w.DrawWindow()
		}
	}
	MoveCursorTo(s.Height - 1, 1)
}

func (s Screen) Clear() {
	fill := strings.Repeat(" ", s.Width)
	for i := 0; i <= s.Height; i++ {
		MoveCursorTo(i, 0)
		fmt.Print(fill)
	}
}

func (s Screen) SetCharMode() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	fmt.Print("\033[?25l")
}

// Checks for a screen resize and resizes windows if needed
// Then redraws the screen. Takes a bool to decide whether
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
			title = string(bar.title[:s.Width - 3]) + "..."
		}
		_, _ = buf.WriteString(title)
		msg := bar.message
		if len(bar.message) > s.Width - len(title) {
			msg = string(bar.message[:s.Width - len(title) - 3]) + "..."
		}
		_, _ = buf.WriteString(msg)

		MoveCursorTo(bar.row, 1)
		fmt.Print(buf.String())
		fmt.Print("\033[0m")

	}
}

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


func NewScreen() *Screen {
	if screenInit {
		fmt.Println("Fatal error: Cannot create multiple screens")
		os.Exit(1)
	}
	var s Screen
	s.GetSize()
	for i := 0; i < s.Height; i++ {
		fmt.Println()
	}
	SetCharMode()
	Clear("screen")
	screenInit = true
	return &s
}
