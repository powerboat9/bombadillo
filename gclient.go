package main

import (
	"fmt"
	"gsock/gopher"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"gsock/cui"
	"errors"
)

var history gopher.History = gopher.MakeHistory()
var screen *cui.Screen
var userinfo, _ = user.Current()

func err_exit(err string, code int) {
	fmt.Println(err)
	os.Exit(code)
}

func save_file() {
	//TODO add a way to save a file...
	//eg. :save 5 test.txt
}

func search(u string) error {
	cui.MoveCursorTo(screen.Height - 1, 0)
	cui.Clear("line")
	fmt.Print("Enter form input: ")
	cui.MoveCursorTo(screen.Height - 1, 17)
	entry := cui.GetLine()
	searchurl := fmt.Sprintf("%s\t%s", u, entry[:len(entry) - 1])
	sv, err := gopher.Visit(searchurl)
	if err != nil {
		return err
	}
	history.Add(sv)

	return nil
}


func route_input(s string) error {
	if num, _ := regexp.MatchString(`^\d+$`, s); num && history.Length > 0 {
		linkcount := len(history.Collection[history.Position].Links)
		item, _ := strconv.Atoi(s)
		if item <= linkcount {
			linkurl := history.Collection[history.Position].Links[item - 1]
			v, err := gopher.Visit(linkurl)
			if err != nil {
				return err
			}

			if v.Address.Gophertype == "7" {
				err := search(linkurl)
				if err != nil {
					return err
				}
			} else if v.Address.IsBinary {
				// TODO add download querying here
			} else {
				history.Add(v)
			}
		} else {
			errname := fmt.Sprintf("Invalid link id: %s", s)
			return errors.New(errname)
		}
	} else {
		v, err := gopher.Visit(s)
		if err != nil {
			return err
		}
		if v.Address.Gophertype == "7" {
			err := search(v.Address.Full)
			if err != nil {
				return err
			}
		} else if v.Address.IsBinary {
			// TODO add download querying here
		} else {
			history.Add(v)
		}
	}
	return nil
}


func main() {
	// fmt.Println(userinfo.HomeDir)
	history.Position = -1
	screen = cui.NewScreen()
	screen.SetCharMode()
	defer cui.Exit()
	screen.AddWindow(1, 1, screen.Height - 2, screen.Width, false, false)
	mainWindow := screen.Windows[0]
	redrawScreen := true

	for {
		screen.ReflashScreen(redrawScreen)

		redrawScreen = false

		c := cui.Getch()
		switch c {
			case 'j', 'J':
				mainWindow.ScrollDown()
			case 'k', 'K':
				mainWindow.ScrollUp()
			case 'q', 'Q':
				cui.Exit()
			case 'b', 'B':
				history.GoBack()
				mainWindow.Scrollposition = 0
				redrawScreen = true
			case 'f', 'F':
				history.GoForward()
				mainWindow.Scrollposition = 0
				redrawScreen = true
			case ':':
				redrawScreen = true
				cui.MoveCursorTo(screen.Height - 1, 0)
				entry := cui.GetLine()
				// Clear entry line
				cui.MoveCursorTo(screen.Height - 1, 0)
				cui.Clear("line")
				if entry == "" {
					cui.MoveCursorTo(screen.Height - 1, 0)
					fmt.Print(" ")
					continue
				}
				err := route_input(entry)
				if err != nil {
					// Display error
					cui.MoveCursorTo(screen.Height, 0)
					fmt.Print("\033[41m\033[37m", err, "\033[0m")
					// Set screen to not reflash
					redrawScreen = false
				} else {
					mainWindow.Scrollposition = 0
					// screen.Clear()
				}
		}
		if history.Position >= 0 {
			mainWindow.Content = history.Collection[history.Position].Content
		}
	}
}
