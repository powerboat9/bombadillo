package main

// Bombadillo is an internet client for the terminal of unix or
// unix-like systems.
//
// Copyright (C) 2019 Brian Evans
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"tildegit.org/sloum/bombadillo/config"
	"tildegit.org/sloum/bombadillo/cui"
	"tildegit.org/sloum/bombadillo/gemini"
)

var version string = "2.3.3"

var bombadillo *client
var helplocation string = "gopher://bombadillo.colorfield.space:70/1/user-guide.map"
var settings config.Config

func saveConfig() error {
	var opts strings.Builder
	bkmrks := bombadillo.BookMarks.IniDump()
	certs := bombadillo.Certs.IniDump()

	opts.WriteString("\n[SETTINGS]\n")
	for k, v := range bombadillo.Options {
		opts.WriteString(k)
		opts.WriteRune('=')
		opts.WriteString(v)
		opts.WriteRune('\n')
	}

	opts.WriteString(bkmrks)

	opts.WriteString(certs)

	return ioutil.WriteFile(filepath.Join(bombadillo.Options["configlocation"], ".bombadillo.ini"), []byte(opts.String()), 0644)
}

func validateOpt(opt, val string) bool {
	var validOpts = map[string][]string{
		"webmode":       []string{"none", "gui", "lynx", "w3m", "elinks"},
		"theme":         []string{"normal", "inverse", "color"},
		"defaultscheme": []string{"gopher", "gemini", "http", "https"},
		"showimages":    []string{"true", "false"},
		"geminiblocks":  []string{"block", "neither", "alt", "both"},
	}

	opt = strings.ToLower(opt)
	val = strings.ToLower(val)

	if _, ok := validOpts[opt]; ok {
		for _, item := range validOpts[opt] {
			if item == val {
				return true
			}
		}
		return false
	}

	if opt == "timeout" {
		_, err := strconv.Atoi(val)
		if err != nil {
			return false
		}
	}

	return true
}

func lowerCaseOpt(opt, val string) string {
	switch opt {
	case "webmode", "theme", "defaultscheme", "showimages", "geminiblocks":
		return strings.ToLower(val)
	default:
		return val
	}
}

func loadConfig() {
	err := os.MkdirAll(bombadillo.Options["configlocation"], 0755)
	if err != nil {
		exitMsg := fmt.Sprintf("Error creating 'configlocation' directory: %s", err.Error())
		cui.Exit(3, exitMsg)
	}

	fp := filepath.Join(bombadillo.Options["configlocation"], ".bombadillo.ini")
	file, err := os.Open(fp)
	if err != nil {
		err = saveConfig()
		if err != nil {
			exitMsg := fmt.Sprintf("Error writing config file during bootup: %s", err.Error())
			cui.Exit(4, exitMsg)
		}
	}

	confparser := config.NewParser(file)
	settings, _ = confparser.Parse()
	_ = file.Close()
	for _, v := range settings.Settings {
		lowerkey := strings.ToLower(v.Key)
		if lowerkey == "configlocation" {
			// Read only
			continue
		}

		if _, ok := bombadillo.Options[lowerkey]; ok {
			if validateOpt(lowerkey, v.Value) {
				bombadillo.Options[lowerkey] = v.Value
				if lowerkey == "geminiblocks" {
					gemini.BlockBehavior = v.Value
				} else if lowerkey == "timeout" {
					updateTimeouts(v.Value)
				}
			} else {
				bombadillo.Options[lowerkey] = defaultOptions[lowerkey]
			}
		}
	}

	for i, v := range settings.Bookmarks.Titles {
		_, _ = bombadillo.BookMarks.Add([]string{v, settings.Bookmarks.Links[i]})
	}

	for _, v := range settings.Certs {
		// Remove expired certs
		vals := strings.SplitN(v.Value, "|", -1)
		if len(vals) < 2 {
			continue
		}
		now := time.Now()
		ts, err := strconv.ParseInt(vals[1], 10, 64)
		if err != nil || now.Unix() > ts {
			continue
		}
		// Satisfied that the cert is not expired
		// or malformed: add to the current client
		// instance
		bombadillo.Certs.Add(v.Key, vals[0], ts)
	}
}

func initClient() {
	bombadillo = MakeClient("  ((( Bombadillo )))  ")
	loadConfig()
}

// In the event of specific signals, ensure the display is shown correctly.
// Accepts a signal, blocking until it is received.  Once not blocked, corrects
// terminal display settings as appropriate for that signal. Loops
// indefinitely, does not return.
func handleSignals(c <-chan os.Signal) {
	for {
		switch <-c {
		case syscall.SIGTSTP:
			cui.CleanupTerm()
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGSTOP)
		case syscall.SIGCONT:
			cui.InitTerm()
			bombadillo.Draw()
		case syscall.SIGINT:
			cui.Exit(130, "")
		}
	}
}

//printHelp produces a nice display message when the --help flag is used
func printHelp() {
	art := `Bombadillo - a non-web browser

Syntax:   bombadillo [options] [url] 

Examples: bombadillo gopher://bombadillo.colorfield.space
          bombadillo -t 
          bombadillo -v

Options: 
`
	_, _ = fmt.Fprint(os.Stdout, art)
	flag.PrintDefaults()
}

func main() {
	getVersion := flag.Bool("v", false, "Display version information and exit")
	addTitleToXWindow := flag.Bool("t", false, "Set the window title to 'Bombadillo'. Can be used in a GUI environment, however not all terminals support this feature.")
	flag.Usage = printHelp
	flag.Parse()
	if *getVersion {
		fmt.Printf("Bombadillo %s\n", version)
		os.Exit(0)
	}
	args := flag.Args()

	cui.InitTerm()

	if *addTitleToXWindow {
		fmt.Print("\033[22;0t")            // Store window title on terminal stack
		fmt.Print("\033]0;Bombadillo\007") // Update window title
	}

	defer cui.Exit(0, "")
	initClient()

	// watch for signals, send them to be handled
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTSTP, syscall.SIGCONT, syscall.SIGINT)
	go handleSignals(c)

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
