package main

import (
	"os"
	"os/user"
	"path/filepath"
)

var defaultOptions = map[string]string{
	// The configuration options below control the default settings for
	// users of Bombadillo.
	//
	// Changes take effect when Bombadillo is built. Follow the standard
	// install instructions after making a change.
	//
	// Most options can be changed by a user in the Bombadillo client, and
	// changes made here will not overwrite an existing user's settings.
	// The exception to both cases is "configlocation" which controls where
	// .bombadillo.ini is stored. If you make changes to this setting,
	// consider moving bombadillo.ini to the new location as well, so you
	// (or your users) do not loose bookmarks or other preferences.
	//
	// Further explanation of each option is available in the man page.

	// Basic Usage
	//
	// Any option can be defined as a string, like this:
	// "option": "value"
	//
	// Options can also have values calculated on startup. There are two
	// functions below that do just this: homePath() and xdgConfigPath()
	// You can set any value to use these functions like this:
	// "option": homePath()
	// "option": xdgConfigPath()
	// See the comments for these functions for more information on what
	// they do.
	//
	// You can also use `filepath.Join()` if you want to build a file path.
	// For example, specify "~/bombadillo" like so:
	// "option": filepath.Join(homePath(), bombadillo)

	// Moving .bombadillo.ini out of your home directory
	//
	// To ensure .bombadillo.ini is saved as per XDG config spec, change
	// the "configlocation" as follows:
	// "configlocation": xdgConfigPath()

	"configlocation": xdgConfigPath(),
	"defaultscheme":  "gopher", // "gopher", "gemini", "http", "https"
	"geminiblocks":   "block",  // "block", "alt", "neither", "both"
	"homeurl":        "gopher://bombadillo.colorfield.space:70/1/user-guide.map",
	"savelocation":   homePath(),
	"searchengine":   "gopher://gopher.floodgap.com:70/7/v2/vs",
	"showimages":     "true",
	"telnetcommand":  "telnet",
	"theme":          "normal", // "normal", "inverted", "color"
	"timeout":        "15",     // connection timeout for gopher/gemini in seconds
	"webmode":        "none",   // "none", "gui", "lynx", "w3m", "elinks"
}

// homePath will return the path to your home directory as a string
// Usage:
//	"configlocation": homeConfigPath()
func homePath() string {
	var userinfo, _ = user.Current()
	return userinfo.HomeDir
}

// xdgConfigPath returns the path to your XDG base directory for configuration
// i.e the contents of environment variable XDG_CONFIG_HOME, or ~/.config/
// Usage:
//	"configlocation": xdgConfigPath()
func xdgConfigPath() string {
	configPath := os.Getenv("XDG_CONFIG_HOME")
	if configPath == "" {
		return filepath.Join(homePath(), ".config")
	}
	return configPath
}
