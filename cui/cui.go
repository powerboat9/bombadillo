package cui

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

var shapes = map[string]string{
	"wall":     "╵",
	"ceiling":  "╴",
	"tl":       "┌",
	"tr":       "┐",
	"bl":       "└",
	"br":       "┘",
	"awall":    "║",
	"aceiling": "═",
	"atl":      "╔",
	"atr":      "╗",
	"abl":      "╚",
	"abr":      "╝",
}

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
		"up":    "A",
		"down":  "B",
		"left":  "D",
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
	HandleAlternateScreen("rmcup")
	os.Exit(0)
}

func Clear(dir string) {
	directions := map[string]string{
		"up":     "\033[1J",
		"down":   "\033[0J",
		"left":   "\033[1K",
		"right":  "\033[0K",
		"line":   "\033[2K",
		"screen": "\033[2J",
	}

	if val, ok := directions[dir]; ok {
		fmt.Print(val)
	}

}

// takes the document content (as a slice) and modifies any lines that are longer
// than the specified console width, splitting them over two lines. returns the
// amended document content as a slice.
// word wrapping uses a "greedy" algorithm
func wrapLines(s []string, consolewidth int) []string {
	out := []string{}
	for _, ln := range s {
		if len([]rune(ln)) <= consolewidth {
			out = append(out, ln)
		} else {
			words := strings.SplitAfter(ln, " ")
			var subout bytes.Buffer
			for i, wd := range words {
				sublen := subout.Len()
				wdlen := len([]rune(wd))
				if sublen+wdlen <= consolewidth {
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
					}
				} else {
					out = append(out, subout.String())
					subout.Reset()
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
						subout.Reset()
					}
				}
			}
		}
	}
	return out
}

func Getch() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return '@'
	}
	return char
}

func GetLine() (string, error) {
	SetLineMode()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	SetCharMode()
	return text[:len(text)-1], nil
}

func SetCharMode() {
	cmd := exec.Command("stty", "cbreak", "-echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Print("\033[?25l")
}

func SetLineMode() {
	cmd := exec.Command("stty", "-cbreak", "echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func HandleAlternateScreen(opt string) {
	cmd := exec.Command("tput", opt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	// explicitly ignoring the error here as
	// the alternate screen is an optional feature
	// that may not be available everywhere we expect
	// to run
	_ = cmd.Run()
}

func wrapLines2(s []string, consolewidth int) []string {
	out := []string{}
	for _, ln := range s {
		if utf8.RuneCountInString(ln) <= consolewidth {
			out = append(out, ln)
		} else {
			words := strings.SplitAfter(ln, " ")
			var subout bytes.Buffer
			for i, wd := range words {
				sublen := subout.Len()
				wdlen := utf8.RuneCountInString(wd)
				if sublen+wdlen <= consolewidth {
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
					}
				} else {
					out = append(out, subout.String())
					subout.Reset()
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
						subout.Reset()
					}
				}
			}
		}
	}
	return out
}

func wrapLines3(s []string, consolewidth int) []string {
	out := []string{}
	for _, ln := range s {
		if len(ln) <= consolewidth {
			out = append(out, ln)
		} else {
			words := strings.SplitAfter(ln, " ")
			var subout bytes.Buffer
			for i, wd := range words {
				sublen := subout.Len()
				wdlen := len(wd)
				if sublen+wdlen <= consolewidth {
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
					}
				} else {
					out = append(out, subout.String())
					subout.Reset()
					subout.WriteString(wd)
					if i == len(words)-1 {
						out = append(out, subout.String())
						subout.Reset()
					}
				}
			}
		}
	}
	return out
}
