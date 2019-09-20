package main

import (
	"os/user"
)

var userinfo, _ = user.Current()
var defaultOptions = map[string]string{
	//
	// General configuration options
	//
	"homeurl":      "gopher://colorfield.space:70/1/bombadillo-info",
	"savelocation": userinfo.HomeDir,
	"searchengine": "gopher://gopher.floodgap.com:70/7/v2/vs",
	"openhttp":     "false",
	"httpbrowser":  "lynx",
	"telnetcommand": "telnet",
	"configlocation": userinfo.HomeDir,
	"theme": "normal", // "normal", "inverted"
}

