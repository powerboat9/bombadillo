package cui

import (
	"strings"
	"bytes"
	"fmt"
	"bufio"
	"os"
	"os/exec"
)

var shapes = map[string]string{
	"wall": "│",
	"ceiling": "─",
	"tl": "┌",
	"tr": "┐",
	"bl": "└",
	"br": "┘",
	"scroll-thumb": "▉",
	"scroll-track": "░",
}

var screenInit bool = false

type Screen struct {
	Height				int
	Width					int
	Windows				[]*Window
  Activewindow	int
}

type box struct {
	row1			int
	col1			int
	row2			int
	col2			int
}

type Window struct {
	Box							box
	Scrollbar				bool
	Scrollposition	int
	Content					[]string
	drawBox					bool
	Active					bool
}

func (s *Screen) AddWindow(r1, c1, r2, c2 int, scroll, border bool) {
	w := Window{box{r1, c1, r2, c2}, scroll, 0, []string{}, border, false}
	s.Windows = append(s.Windows, &w)
}

func (s Screen) DrawFullScreen() {
	s.Clear()
	// w := s.Windows[s.Activewindow]
	for _, w := range s.Windows {
		if w.drawBox {
			w.DrawBox()
		}

		w.DrawContent()
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
		for _, w := range s.Windows {
			w.Box.row2 = s.Height - 2
			w.Box.col2 = s.Width
		}
		s.DrawFullScreen()
	} else if clearScreen {
		s.DrawFullScreen()
	} else {
		s.Windows[s.Activewindow].DrawContent()
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

func (w *Window) DrawBox(){
	moveThenDrawShape(w.Box.row1, w.Box.col1, "tl")
	moveThenDrawShape(w.Box.row1, w.Box.col2, "tr")
	moveThenDrawShape(w.Box.row2, w.Box.col1, "bl")
	moveThenDrawShape(w.Box.row2, w.Box.col2, "br")
	for i := w.Box.col1 + 1; i < w.Box.col2; i++ {
		moveThenDrawShape(w.Box.row1, i, "ceiling")
		moveThenDrawShape(w.Box.row2, i, "ceiling")
	}

	for i:= w.Box.row1 + 1; i < w.Box.row2; i++ {
		moveThenDrawShape(i, w.Box.col1, "wall")
		moveThenDrawShape(i, w.Box.col2, "wall")
	}
}

func (w *Window) DrawContent(){
	var maxlines, borderw, contenth int
	if w.drawBox {
		borderw, contenth = -1, 1
	} else {
		borderw, contenth = 1, 0
	}
	height, width := w.Box.row2 - w.Box.row1 + borderw, w.Box.col2 - w.Box.col1 + borderw
	content := WrapLines(w.Content, width)
	if len(content) < w.Scrollposition + height {
		maxlines = len(content)
	} else {
		maxlines = w.Scrollposition + height
	}

	for i := w.Scrollposition; i < maxlines; i++ {
		MoveCursorTo(w.Box.row1 + contenth + i - w.Scrollposition, w.Box.col1 + contenth)
		fmt.Print( strings.Repeat(" ", width) )
		MoveCursorTo(w.Box.row1 + contenth + i - w.Scrollposition, w.Box.col1 + contenth)
		fmt.Print(content[i])
	}
}

func (w *Window) ScrollDown() {
	height := w.Box.row2 - w.Box.row1 - 1
	contentLength := len(w.Content)
	if w.Scrollposition < contentLength - height {
		w.Scrollposition++
	} else {
		fmt.Print("\a")
	}
}

func (w *Window) ScrollUp() {
	if w.Scrollposition > 0 {
		w.Scrollposition--
	} else {
		fmt.Print("\a")
	}
}




//--------------------------------------------------------------------------//
//                                                                          //
//                          F U N C T I O N S                               //
//                                                                          //
//--------------------------------------------------------------------------//


func drawShape(shape string) {
	if val, ok := shapes[shape]; ok {
		fmt.Printf("%s", val)
	} else {
		fmt.Print("x")
	}
}

func moveThenDrawShape(r, c int, s string) {
	MoveCursorTo(r, c)
	drawShape(s)
}

func MoveCursorTo(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func moveCursorToward(dir string, amount int) {
	directions := map[string]string{
		"up": "A",
		"down": "B",
		"left": "D",
		"right": "C",
	}

	if val, ok := directions[dir]; ok {
		fmt.Printf("\033[%d%s", amount, val)
	}
}

func Exit() {
	moveCursorToward("down", 500)
	moveCursorToward("right", 500)
	SetLineMode()
	fmt.Print("\n")
	fmt.Print("\033[?25h")
	os.Exit(0)
}

func Clear(dir string) {
	directions := map[string]string{
		"up": "\033[1J",
		"down": "\033[0J",
		"left": "\033[1K",
		"right": "\033[0K",
		"line": "\033[2K",
		"screen": "\033[2J",
	}

	if val, ok := directions[dir]; ok {
		fmt.Print(val)
	}

}

func WrapLines(s []string, length int) []string {
	out := []string{}
	for _, ln := range s {
		if len(ln) <= length {
			out = append(out, ln)
		} else {
			words := strings.Split(ln, " ")
			var subout bytes.Buffer
			for i, wd := range words {
				sublen := subout.Len()
				if sublen + len(wd) + 1 <= length {
					if sublen > 0 {
						subout.WriteString(" ")
					}
					subout.WriteString(wd)	
					if i == len(words) - 1 {
						out = append(out, subout.String())
					}
				} else {
						out = append(out, subout.String())
						subout.Reset()
						subout.WriteString(wd)
						if i == len(words) - 1 {
							out = append(out, subout.String())
							subout.Reset()
						}
				}
			}
		}
	}
	return out
}


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

func Getch() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return '@'
	}
	return char
}

func GetLine() string {
	SetLineMode()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(": ")
	text, _ := reader.ReadString('\n')
	SetCharMode()
	return text[:len(text)-1]
}

func SetCharMode() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	fmt.Print("\033[?25l")
}

func SetLineMode() {
	exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
