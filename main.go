package main

// Bombadillo is a gopher and gemini client for the terminal of unix or unix-like systems.
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
	"strings"
	"syscall"

	"tildegit.org/sloum/bombadillo/config"
	"tildegit.org/sloum/bombadillo/cui"
	_ "tildegit.org/sloum/bombadillo/gemini"
	"tildegit.org/sloum/mailcap"
)

var version string
var build string

var bombadillo *client
var helplocation string = "gopher://bombadillo.colorfield.space:70/1/user-guide.map"
var settings config.Config
var mc *mailcap.Mailcap

func saveConfig() error {
	var opts strings.Builder
	bkmrks := bombadillo.BookMarks.IniDump()
	certs := bombadillo.Certs.IniDump()

	opts.WriteString("\n[SETTINGS]\n")
	for k, v := range bombadillo.Options {
		if k == "theme" && v != "normal" && v != "inverse" {
			v = "normal"
			bombadillo.Options["theme"] = "normal"
		}
		opts.WriteString(k)
		opts.WriteRune('=')
		opts.WriteString(v)
		opts.WriteRune('\n')
	}

	opts.WriteString(bkmrks)

	opts.WriteString(certs)

	return ioutil.WriteFile(bombadillo.Options["configlocation"]+"/.bombadillo.ini", []byte(opts.String()), 0644)
}

func validateOpt(opt, val string) bool {
	var validOpts = map[string][]string{
		"openhttp":     []string{"true", "false"},
		"theme":        []string{"normal", "inverse"},
		"terminalonly": []string{"true", "false"},
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
	} else {
		return true
	}
}

func lowerCaseOpt(opt, val string) string {
	switch opt {
	case "openhttp", "theme", "terminalonly":
		return strings.ToLower(val)
	default:
		return val
	}
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
			// The config defaults to the home folder.
			// Users cannot really edit this value. But
			// a compile time override is available.
			// It is still stored in the ini and as a part
			// of the options map.
			continue
		}

		if _, ok := bombadillo.Options[lowerkey]; ok {
			if validateOpt(lowerkey, v.Value) {
				bombadillo.Options[lowerkey] = v.Value
			} else {
				bombadillo.Options[lowerkey] = defaultOptions[lowerkey]
			}
		}
	}

	for i, v := range settings.Bookmarks.Titles {
		bombadillo.BookMarks.Add([]string{v, settings.Bookmarks.Links[i]})
	}

	for _, v := range settings.Certs {
		bombadillo.Certs.Add(v.Key, v.Value)
	}

	return nil
}

func initClient() error {
	bombadillo = MakeClient("  ((( Bombadillo )))  ")
	err := loadConfig()
	if bombadillo.Options["tlscertificate"] != "" && bombadillo.Options["tlskey"] != "" {
		bombadillo.Certs.LoadCertificate(bombadillo.Options["tlscertificate"], bombadillo.Options["tlskey"])
	}
	return err
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
			syscall.Kill(syscall.Getpid(), syscall.SIGSTOP)
		case syscall.SIGCONT:
			cui.InitTerm()
			bombadillo.Draw()
		case syscall.SIGINT:
			cui.Exit()
		}
	}
}

//printHelp produces a nice display message when the --help flag is used
func printHelp() {
	art := `Bombadillo - a non-web client

Syntax:   bombadillo [url] 
          bombadillo [options...]

Examples: bombadillo gopher://bombadillo.colorfield.space
          bombadillo -v

Options: 
`
	fmt.Fprint(os.Stdout, art)
	flag.PrintDefaults()
}

func main() {
	getVersion := flag.Bool("v", false, "Display version information and exit")
	flag.Usage = printHelp
	flag.Parse()
	if *getVersion {
		fmt.Printf("Bombadillo %s - build %s\n", version, build)
		os.Exit(0)
	}
	args := flag.Args()

	// Build the mailcap db
	// So that we can open files from gemini
	mc = mailcap.NewMailcap()

	cui.InitTerm()
	defer cui.Exit()
	err := initClient()
	if err != nil {
		// if we can't initialize we should bail out
		panic(err)
	}

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
