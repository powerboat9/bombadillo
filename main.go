package main

// Bombadillo is distributed under the "Non-Profit Open Source Software License 3.0"
// The license is included with the source code in the file LICENSE. The basic
// takeway: use, remix, and share this software for any purpose that is not a commercial
// purpose as defined by the above mentioned license and is itself distributed udner
// the terms of said license with said license file included.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"tildegit.org/sloum/bombadillo/config"
	"tildegit.org/sloum/bombadillo/cui"
)

const version = "2.0.0"

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
	var opts strings.Builder
	bkmrks := bombadillo.BookMarks.IniDump()

	opts.WriteString(bkmrks)
	opts.WriteString("\n[SETTINGS]\n")
	for k, v := range bombadillo.Options {
		opts.WriteString(k)
		opts.WriteRune('=')
		opts.WriteString(v)
		opts.WriteRune('\n')
	}

	return ioutil.WriteFile(bombadillo.Options["configlocation"] + "/.bombadillo.ini", []byte(opts.String()), 0644)
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
	getVersion := flag.Bool("v", false, "See version number")
	flag.Parse()
	if *getVersion {
		fmt.Printf("Bombadillo v%s\n", version)
		os.Exit(0)
	}
	args := flag.Args()

	cui.Tput("rmam") // turn off line wrapping
	cui.Tput("smcup") // use alternate screen
	defer cui.Exit()
	err := initClient()
	if err != nil {
		// if we can't initialize we should bail out
		panic(err)
	}

	// Start polling for terminal size changes
	go bombadillo.GetSize()

	if len(args) > 0 {
		// If a url was passed, move it down the line
		// Goroutine so keypresses can be made during
		// page load
		bombadillo.Visit(args[0])
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
