package main

import (
	"os/user"
)

var userinfo, _ = user.Current()
var defaultOptions = map[string]string{
	//
	// General configuration options
	//
	// Edit these values before compile to have different default values
	// ... though they can always be edited from within bombadillo as well
	// it just may take more time/work.
  //
  // To change the default location for the config you can enter
  // any valid path as a string, if you want an absolute, or
  // concatenate with the main default: `userinfo.HomeDir` like so:
  // "configlocation": userinfo.HomeDir + "/config/"
	"homeurl":      "gopher://bombadillo.colorfield.space:70/1/user-guide.map",
	"savelocation": userinfo.HomeDir,
	"searchengine": "gopher://gopher.floodgap.com:70/7/v2/vs",
	"openhttp":     "false",
	"telnetcommand": "telnet",
	"configlocation": userinfo.HomeDir,
	"theme": "normal", // "normal", "inverted"
	"terminalonly": "true",
	"tlscertificate": "",
	"tlskey": "",
	"lynxmode": "false",
}

