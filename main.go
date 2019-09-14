package main

// Bombadillo is distributed under the "Non-Profit Open Source Software License 3.0"
// The license is included with the source code in the file LICENSE. The basic
// takeway: use, remix, and share this software for any purpose that is not a commercial
// purpose as defined by the above mentioned license and is itself distributed udner
// the terms of said license with said license file included.

import (
	"io/ioutil"
	"os"
	"strings"

	"tildegit.org/sloum/bombadillo/config"
	"tildegit.org/sloum/bombadillo/cui"
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

	// TODO find out why the loading message
	// has disappeared on initial load...

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
