package main

import (
	"io/ioutil"
	"os"
	// "strconv"
	"strings"

	"tildegit.org/sloum/bombadillo/config"
	"tildegit.org/sloum/bombadillo/cui"
	// "tildegit.org/sloum/bombadillo/gopher"
)

var bombadillo *client
var helplocation string = "gopher://colorfield.space:70/1/bombadillo-info"
var settings config.Config


// func saveFileFromData(v gopher.View) error {
	// quickMessage("Saving file...", false)
	// urlsplit := strings.Split(v.Address.Full, "/")
	// filename := urlsplit[len(urlsplit)-1]
	// saveMsg := fmt.Sprintf("Saved file as %q", options["savelocation"]+filename)
	// err := ioutil.WriteFile(options["savelocation"]+filename, []byte(strings.Join(v.Content, "")), 0644)
	// if err != nil {
		// quickMessage("Saving file...", true)
		// return err
	// }

	// quickMessage(saveMsg, false)
	// return nil
// }



// func doLinkCommand(action, target string) error {
	// num, err := strconv.Atoi(target)
	// if err != nil {
		// return fmt.Errorf("Expected number, got %q", target)
	// }

	// switch action {
	// case "DELETE", "D":
		// err := settings.Bookmarks.Del(num)
		// if err != nil {
			// return err
		// }

		// screen.Windows[1].Content = settings.Bookmarks.List()
		// err = saveConfig()
		// if err != nil {
			// return err
		// }

		// screen.ReflashScreen(false)
		// return nil
	// case "BOOKMARKS", "B":
		// if num > len(settings.Bookmarks.Links)-1 {
			// return fmt.Errorf("There is no bookmark with ID %d", num)
		// }
		// err := goToURL(settings.Bookmarks.Links[num])
		// return err
	// }

	// return fmt.Errorf("This method has not been built")
// }


// func doCommand(action string, values []string) error {
	// if length := len(values); length != 1 {
		// return fmt.Errorf("Expected 1 argument, received %d", length)
	// }

	// switch action {
	// case "CHECK", "C":
		// err := checkConfigValue(values[0])
		// if err != nil {
			// return err
		// }
		// return nil
	// }
	// return fmt.Errorf("Unknown command structure")
// }

// func doLinkCommandAs(action, target string, values []string) error {
	// num, err := strconv.Atoi(target)
	// if err != nil {
		// return fmt.Errorf("Expected number, got %q", target)
	// }

	// links := history.Collection[history.Position].Links
	// if num >= len(links) {
		// return fmt.Errorf("Invalid link id: %s", target)
	// }

	// switch action {
	// case "ADD", "A":
		// newBookmark := append([]string{links[num-1]}, values...)
		// err := settings.Bookmarks.Add(newBookmark)
		// if err != nil {
			// return err
		// }

		// screen.Windows[1].Content = settings.Bookmarks.List()

		// err = saveConfig()
		// if err != nil {
			// return err
		// }

		// screen.ReflashScreen(false)
		// return nil
	// case "WRITE", "W":
		// return saveFile(links[num-1], strings.Join(values, " "))
	// }

	// return fmt.Errorf("This method has not been built")
// }

// func updateMainContent() {
	// screen.Windows[0].Content = history.Collection[history.Position].Content
	// screen.Bars[0].SetMessage(history.Collection[history.Position].Address.Full)
// }

func saveConfig() error {
	bkmrks := bombadillo.BookMarks.IniDump()
	// TODO opts becomes a string builder rather than concat
	opts := "\n[SETTINGS]\n"
	for k, v := range bombadillo.Options {
		opts += k
		opts += "="
		opts += v
		opts += "\n"
	}

	return ioutil.WriteFile(bombadillo.Options["configlocation"] + "/.bombadillo.ini", []byte(bkmrks+opts), 0644)
}

func loadConfig() error {
	file, err := os.Open(bombadillo.Options["configlocation"] + "/.bombadillo.ini")
	if err != nil {
		err = saveConfig()
		if err != nil {
			return err
		}
	}

	confparser := config.NewParser(file)
	settings, _ = confparser.Parse()
	file.Close()
	for _, v := range settings.Settings {
		lowerkey := strings.ToLower(v.Key)
		if lowerkey == "configlocation" {
			// The config should always be stored in home
			// folder. Users cannot really edit this value.
			// It is still stored in the ini and as a part
			// of the options map.
			continue
		}

		if _, ok := bombadillo.Options[lowerkey]; ok {
			bombadillo.Options[lowerkey] = v.Value
		}
	}

	for i, v := range settings.Bookmarks.Titles {
		bombadillo.BookMarks.Add([]string{v, settings.Bookmarks.Links[i]})
	}

	return nil
}

func initClient() error {
	bombadillo = MakeClient("  ((( Bombadillo )))  ")
	cui.SetCharMode()
	err := loadConfig()
	return err
}

func main() {
	cui.HandleAlternateScreen("smcup")
	defer cui.Exit()
	err := initClient()
	if err != nil {
		// if we can't initialize we should bail out
		panic(err)
	}

	// Start polling for terminal size changes
	go bombadillo.GetSize()

	if len(os.Args) > 1 {
		// If a url was passed, move it down the line
		// Goroutine so keypresses can be made during
		// page load
		bombadillo.Visit(os.Args[1])
	} else {
		// Otherwise, load the homeurl
		// Goroutine so keypresses can be made during
		// page load
		bombadillo.Visit(bombadillo.Options["homeurl"])
	}

	// Loop indefinitely on user input
	for {
		bombadillo.TakeControlInput()
	}
}
