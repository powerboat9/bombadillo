package cui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"tildegit.org/sloum/bombadillo/termios"
)

var Shapes = map[string]string{
	"walll":    "╎",
	"wallr":    " ",
	"ceiling":  " ",
	"floor":    " ",
	"tl":       "╎",
	"tr":       " ",
	"bl":       "╎",
	"br":       " ",
	"awalll":   "▌",
	"awallr":   "▐",
	"aceiling": "▀",
	"afloor":   "▄",
	"atl":      "▞",
	"atr":      "▜",
	"abl":      "▚",
	"abr":      "▟",
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

// Exit performs cleanup operations before exiting the application
func Exit(exitCode int, msg string) {
	CleanupTerm()
	if msg != "" {
		fmt.Print(msg, "\n")
	}
	fmt.Print("\033[23;0t") // Restore window title from terminal stack
	os.Exit(exitCode)
}

// InitTerm sets the terminal modes appropriate for Bombadillo
func InitTerm() {
	termios.SetCharMode()
	Tput("smcup")          // use alternate screen
	Tput("rmam")           // turn off line wrapping
	fmt.Print("\033[?25l") // hide cursor
}

// CleanupTerm reverts changs to terminal mode made by InitTerm
func CleanupTerm() {
	moveCursorToward("down", 500)
	moveCursorToward("right", 500)
	termios.SetLineMode()

	fmt.Print("\n")
	fmt.Print("\033[?25h") // reenables cursor blinking
	Tput("smam")           // turn on line wrap
	Tput("rmcup")          // stop using alternate screen
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

func Getch() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return '@'
	}
	return char
}

func GetLine(prefix string) (string, error) {
	termios.SetLineMode()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prefix)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	termios.SetCharMode()
	return text[:len(text)-1], nil
}

func Tput(opt string) {
	cmd := exec.Command("tput", opt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	// explicitly ignoring the error here as
	// the alternate screen is an optional feature
	// that may not be available everywhere we expect
	// to run
	_ = cmd.Run()
}
