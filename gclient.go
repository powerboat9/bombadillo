package main

import (
	"gsock/gopher"
	"gsock/cmdparse"
	"gsock/config"
	"gsock/cui"
	"os/user"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"regexp"
	"strconv"
)

var history gopher.History = gopher.MakeHistory()
var screen *cui.Screen
var userinfo, _ = user.Current()
var settings config.Config
var options = map[string]string{
	"homeurl": "",
	"savelocation": userinfo.HomeDir + "/Downloads/",
	"searchengine": "gopher://gopher.floodgap.com:70/7/v2/vs",
	"openhttp": "false",
	"httpbrowser": "lynx",
}

func err_exit(err string, code int) {
	fmt.Println(err)
	os.Exit(code)
}

func save_file(address, name string) error {
	quickMessage("Saving file...", false)
	defer quickMessage("Saving file...", true)

	url, err := gopher.MakeUrl(address)
	if err != nil {
		return err
	}

	data, err := gopher.Retrieve(url)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(userinfo.HomeDir + "/" + name, data, 0644)
	if err != nil {
		return err
	}

	return nil
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


func route_input(com *cmdparse.Command) error {
	var err error
	switch com.Type {
		case cmdparse.SIMPLE:
			err = simple_command(com.Action)
		case cmdparse.GOURL:
			err = go_to_url(com.Target)
		case cmdparse.GOLINK:
			err = go_to_link(com.Target)
		case cmdparse.DOLINK:
			err = do_link_command(com.Action, com.Target)
		case cmdparse.DOAS:
			err = do_command_as(com.Action, com.Value)
		case cmdparse.DOLINKAS:
			err = do_link_command_as(com.Action, com.Target, com.Value)
		default:
			return fmt.Errorf("Unknown command entry!")
	}

	return err
}

func toggle_bookmarks() {
	bookmarks := screen.Windows[1]
	if bookmarks.Show {
		bookmarks.Show = false
	} else {
		bookmarks.Show = true
	}

	if screen.Activewindow == 0 {
		screen.Activewindow = 1
	} else {
		screen.Activewindow = 0
	}
}

func simple_command(a string) error {
	a = strings.ToUpper(a)
	switch a {
		case "Q", "QUIT":
			cui.Exit()
		case "H", "HOME":
			return go_home()
		case "B", "BOOKMARKS":
			toggle_bookmarks()
		default:
			return fmt.Errorf("Unknown action %q", a)
	}
	return nil
}

func go_to_url(u string) error {
		quickMessage("Loading...", false)
		v, err := gopher.Visit(u)
		if err != nil {
			quickMessage("Loading...", true)
			return err
		}
		quickMessage("Loading...", true)

		if v.Address.Gophertype == "7" {
			err := search(v.Address.Full)
			if err != nil {
				return err
			}
		} else if v.Address.IsBinary {
			// TO DO: run this into the write to file method
		} else {
			history.Add(v)
		}
		return nil
}

func go_to_link(l string) error {
	if num, _ := regexp.MatchString(`^\d+$`, l); num && history.Length > 0 {
		linkcount := len(history.Collection[history.Position].Links)
		item, _ := strconv.Atoi(l)
		if item <= linkcount {
			linkurl := history.Collection[history.Position].Links[item - 1]
			quickMessage("Loading...", false)
			v, err := gopher.Visit(linkurl)
			if err != nil {
				quickMessage("Loading...", true)
				return err
			}
			quickMessage("Loading...", true)

			if v.Address.Gophertype == "7" {
				err := search(linkurl)
				if err != nil {
					return err
				}
			} else if v.Address.IsBinary {
				// TO DO: run this into the write to file method
			} else {
				history.Add(v)
			}
		} else {
			return fmt.Errorf("Invalid link id: %s", l)
		}
	} else {
			return fmt.Errorf("Invalid link id: %s", l)
	}
	return nil
}

func go_home() error {
	if options["homeurl"] != "" {
		return go_to_url(options["homeurl"])
	}
	return fmt.Errorf("No home address has been set")
}

func do_link_command(action, target string) error {
	num, err := strconv.Atoi(target)
	if err != nil {
		return fmt.Errorf("Expected number, got %q", target)
	}

	switch action {
		case "DELETE", "D":
			err := settings.Bookmarks.Del(num)
			screen.Windows[1].Content = settings.Bookmarks.List()
			save_config()
			return err
		case "BOOKMARKS", "B":
			if num > len(settings.Bookmarks.Links) - 1 {
				return fmt.Errorf("There is no bookmark with ID %d", num)
			}
			err := go_to_url(settings.Bookmarks.Links[num])
			return err
	}

	return fmt.Errorf("This method has not been built")
}

func do_command_as(action string, values []string) error {
	if len(values) < 2 {
		return fmt.Errorf("%q", values)
	}

	switch action {
		case "ADD", "A":
			if values[0] == "." {
				values[0] = history.Collection[history.Position].Address.Full
			}
			err := settings.Bookmarks.Add(values)
			if err != nil {
				return err
			}
			screen.Windows[1].Content = settings.Bookmarks.List()
			save_config()
			return nil
		case "WRITE", "W":
			return  save_file(values[0], strings.Join(values[1:], " "))
		case "SET", "S":
			if _, ok := options[values[0]]; ok {
				options[values[0]] = strings.Join(values[1:], " ")
				save_config()
				return nil
			}
			return fmt.Errorf("Unable to set %s, it does not exist",values[0])
	}
	return fmt.Errorf("This method has not been built")
}

func do_link_command_as(action, target string, values []string) error {
	num, err := strconv.Atoi(target)
	if err != nil {
		return fmt.Errorf("Expected number, got %q", target)
	}

	links := history.Collection[history.Position].Links
	if num >= len(links) {
		return fmt.Errorf("Invalid link id: %s", target)
	}

	switch action {
		case "ADD", "A":
			newBookmark := append([]string{links[num - 1]}, values...)
			err := settings.Bookmarks.Add(newBookmark)
			if err != nil {
				return err
			}
			screen.Windows[1].Content = settings.Bookmarks.List()
			save_config()
			return nil
		case "WRITE", "W":
			return save_file(links[num - 1], strings.Join(values, " "))
	}

	return fmt.Errorf("This method has not been built")
}

func clearInput(incError bool) {
	cui.MoveCursorTo(screen.Height - 1, 0)
	cui.Clear("line")
	if incError {
		cui.MoveCursorTo(screen.Height, 0)
		cui.Clear("line")
	}
}

func quickMessage(msg string, clearMsg bool) {
	cui.MoveCursorTo(screen.Height, screen.Width - 2 - len(msg))
	if clearMsg {
		cui.Clear("right")
	} else {
		fmt.Print("\033[48;5;21m\033[38;5;15m", msg, "\033[0m")
	}
}

func save_config() {
	bkmrks := settings.Bookmarks.IniDump()
	opts := "\n[SETTINGS]\n"
	for k, v := range options {
		opts += k
		opts += "="
		opts += v
		opts += "\n"
	}
	ioutil.WriteFile(userinfo.HomeDir + "/.badger.ini", []byte(bkmrks+opts), 0644)
}

func load_config() {
	file, _ := os.Open(userinfo.HomeDir + "/.badger.ini")
	confparser := config.NewParser(file)
	settings, _ = confparser.Parse()
	file.Close()
	screen.Windows[1].Content = settings.Bookmarks.List()
	for _, v := range settings.Settings {
		lowerkey := strings.ToLower(v.Key)
		if _, ok := options[lowerkey]; ok {
			options[lowerkey] = v.Value
		}
	}
}

func initClient() {
	history.Position = -1
	screen = cui.NewScreen()
	screen.SetCharMode()
	screen.AddWindow(2, 1, screen.Height - 2, screen.Width, false, false, true)
	screen.AddMsgBar(1, "  ((( Badger )))  ", "  A fun gopher client!", true)
	bookmarksWidth := 40
	if screen.Width < 40 {
		bookmarksWidth = screen.Width
	}
	screen.AddWindow(2, screen.Width - bookmarksWidth, screen.Height - 2, screen.Width, false, true, false)
	load_config()
}

func main() {
	defer cui.Exit()
	initClient()
	mainWindow := screen.Windows[0]
	first_load := true

	redrawScreen := true

	for {
		screen.ReflashScreen(redrawScreen)

		if first_load {
			go_home()
			first_load = false
			mainWindow.Content = history.Collection[history.Position].Content
			screen.Bars[0].SetMessage(history.Collection[history.Position].Address.Full)
			continue
		}

		redrawScreen = false

		c := cui.Getch()
		switch c {
			case 'j', 'J':
				screen.Windows[screen.Activewindow].ScrollDown()
			case 'k', 'K':
				screen.Windows[screen.Activewindow].ScrollUp()
			case 'q', 'Q':
				cui.Exit()
			case 'b':
				history.GoBack()
				mainWindow.Scrollposition = 0
				redrawScreen = true
			case 'B':
				toggle_bookmarks()
			case 'f', 'F':
				history.GoForward()
				mainWindow.Scrollposition = 0
				redrawScreen = true
			case ':':
				redrawScreen = true
				cui.MoveCursorTo(screen.Height - 1, 0)
				entry := cui.GetLine()
				// Clear entry line and error line
				clearInput(true)
				if entry == "" {
					redrawScreen = false
					continue
				}
				parser := cmdparse.NewParser(strings.NewReader(entry))
				p, err := parser.Parse()
				if err != nil {
					cui.MoveCursorTo(screen.Height, 0)
					fmt.Print("\033[41m\033[37m", err, "\033[0m")
					// Set screen to not reflash
					redrawScreen = false
				} else {
					err := route_input(p)
					if err != nil {
						cui.MoveCursorTo(screen.Height, 0)
						fmt.Print("\033[41m\033[37m", err, "\033[0m")
						redrawScreen = false
					} else {
						mainWindow.Scrollposition = 0
					}
				}
		}
		if history.Position >= 0 {
			mainWindow.Content = history.Collection[history.Position].Content
			screen.Bars[0].SetMessage(history.Collection[history.Position].Address.Full)
		}
	}
}
