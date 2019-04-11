package main

import (
	"bombadillo/gopher"
	"bombadillo/cmdparse"
	"bombadillo/config"
	"bombadillo/cui"
	"os/user"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"regexp"
	"strconv"
)

var helplocation string ="gopher://colorfield.space:70/1/bombadillo-info"
var history gopher.History = gopher.MakeHistory()
var screen *cui.Screen
var userinfo, _ = user.Current()
var settings config.Config
var options = map[string]string{
	"homeurl": "gopher://colorfield.space:70/1/bombadillo-info",
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

	err = ioutil.WriteFile(options["savelocation"] + name, data, 0644)
	if err != nil {
		return err
	}

	return fmt.Errorf("Saved file to " + options["savelocation"] + name)
}

func save_file_from_data(v gopher.View) error {
	urlsplit := strings.Split(v.Address.Full, "/")
	filename := urlsplit[len(urlsplit) - 1]
	save_msg := fmt.Sprintf("Saved file as %q", options["savelocation"] + filename)
	quickMessage(save_msg, false)
	defer quickMessage(save_msg, true)
	err := ioutil.WriteFile(options["savelocation"] + filename, []byte(strings.Join(v.Content, "")), 0644)
	if err != nil {
		return err
	}

	return fmt.Errorf(save_msg)
}

func search(u string) error {
	cui.MoveCursorTo(screen.Height - 1, 0)
	cui.Clear("line")
	fmt.Print("Enter form input: ")
	cui.MoveCursorTo(screen.Height - 1, 17)
	entry := cui.GetLine()
	quickMessage("Searching...", false)
	searchurl := fmt.Sprintf("%s\t%s", u, entry)
	sv, err := gopher.Visit(searchurl, options["openhttp"])
	if err != nil {
		quickMessage("Searching...", true)
		return err
	}
	history.Add(sv)
	quickMessage("Searching...", true)
	updateMainContent()
	screen.ReflashScreen(true)
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
	main := screen.Windows[0]
	if bookmarks.Show {
		bookmarks.Show = false
		screen.Activewindow = 0
		main.Active = true
		bookmarks.Active = false
	} else {
		bookmarks.Show = true
		screen.Activewindow = 1
		main.Active = false
		bookmarks.Active = true
	}

	if screen.Activewindow == 0 {
	} else {
	}
	screen.ReflashScreen(false)
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
		case "SEARCH":
			return search(options["searchengine"])
		case "HELP":
			return go_to_url(helplocation)

		default:
			return fmt.Errorf("Unknown action %q", a)
	}
	return nil
}

func go_to_url(u string) error {
	quickMessage("Loading...", false)
	v, err := gopher.Visit(u, options["openhttp"])
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
		return save_file_from_data(v)
	} else {
		history.Add(v)
	}
	updateMainContent()
	screen.ReflashScreen(true)
	return nil
}

func go_to_link(l string) error {
	if num, _ := regexp.MatchString(`^\d+$`, l); num && history.Length > 0 {
		linkcount := len(history.Collection[history.Position].Links)
		item, _ := strconv.Atoi(l)
		if item <= linkcount {
			linkurl := history.Collection[history.Position].Links[item - 1]
			quickMessage("Loading...", false)
			v, err := gopher.Visit(linkurl, options["openhttp"])
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
				return save_file_from_data(v)
			} else {
				history.Add(v)
			}
		} else {
			return fmt.Errorf("Invalid link id: %s", l)
		}
	} else {
			return fmt.Errorf("Invalid link id: %s", l)
	}
	updateMainContent()
	screen.ReflashScreen(true)
	return nil
}

func go_home() error {
	if options["homeurl"] != "unset" {
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
			screen.ReflashScreen(false)
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

	if values[0] == "." {
		values[0] = history.Collection[history.Position].Address.Full
	}

	switch action {
		case "ADD", "A":
			err := settings.Bookmarks.Add(values)
			if err != nil {
				return err
			}
			screen.Windows[1].Content = settings.Bookmarks.List()
			save_config()
			screen.ReflashScreen(false)
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
	return fmt.Errorf("Unknown command structure")
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
			screen.ReflashScreen(false)
			return nil
		case "WRITE", "W":
			return save_file(links[num - 1], strings.Join(values, " "))
	}

	return fmt.Errorf("This method has not been built")
}


func updateMainContent() {
	screen.Windows[0].Content = history.Collection[history.Position].Content
	screen.Bars[0].SetMessage(history.Collection[history.Position].Address.Full)
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
	ioutil.WriteFile(userinfo.HomeDir + "/.bombadillo.ini", []byte(bkmrks+opts), 0644)
}


func load_config() {
	file, err := os.Open(userinfo.HomeDir + "/.bombadillo.ini")
	if err != nil {
		save_config()
	}
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

func toggleActiveWindow(){
	if screen.Windows[1].Show {
		if screen.Windows[0].Active {
			screen.Windows[0].Active = false
			screen.Windows[1].Active = true
			screen.Activewindow = 1
		} else {
			screen.Windows[0].Active = true
			screen.Windows[1].Active = false
			screen.Activewindow = 0
		}
		screen.Windows[1].DrawWindow()
	}
}


func display_error(err error) {
	cui.MoveCursorTo(screen.Height, 0)
	fmt.Print("\033[41m\033[37m", err, "\033[0m")
}


func initClient() {
	history.Position = -1
	screen = cui.NewScreen()
	cui.SetCharMode()
	screen.AddWindow(2, 1, screen.Height - 2, screen.Width, false, false, true)
	screen.Windows[0].Active = true
	screen.AddMsgBar(1, "  ((( Bombadillo )))  ", "  A fun gopher client!", true)
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

	for {
		if first_load {
			first_load = false
			err := go_home()

			if err == nil {
				updateMainContent()
			}
			continue
		}

		c := cui.Getch()
		switch c {
			case 'j', 'J':
				screen.Windows[screen.Activewindow].ScrollDown()
				screen.ReflashScreen(false)
			case 'k', 'K':
				screen.Windows[screen.Activewindow].ScrollUp()
				screen.ReflashScreen(false)
			case 'q', 'Q':
				cui.Exit()
			case 'b':
				success := history.GoBack()
				if success {
					mainWindow.Scrollposition = 0
					updateMainContent()
					screen.ReflashScreen(true)
				}
			case 'B':
				toggle_bookmarks()
			case 'f', 'F':
				success := history.GoForward()
				if success {
					mainWindow.Scrollposition = 0
					updateMainContent()
					screen.ReflashScreen(true)
				}
			case '\t':
				toggleActiveWindow()
			case ':',' ':
				cui.MoveCursorTo(screen.Height - 1, 0)
				entry := cui.GetLine()
				// Clear entry line and error line
				clearInput(true)
				if entry == "" {
					continue
				}
				parser := cmdparse.NewParser(strings.NewReader(entry))
				p, err := parser.Parse()
				if err != nil {
					display_error(err)
				} else {
					err := route_input(p)
					if err != nil {
						display_error(err)
					}
				}
		}
	}
}
