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

// TODO decide whether or not to institute a color theme
// system. Preliminary testing implies it should be very
// doable.
var theme = map[string]string{
	"topbar_title_bg": "",
	"topbar_link_fg": "",
	"body_bg": "237",
	"body_fg": "",
	"bookmarks_bg": "",
	"bookmarks_fg": "",
	"command_bg": "",
	"message_fg": "",
	"error_fg": "",
	"bottombar_bg": "",
	"bottombar_fg": "",
	//
	// text style options
	//
	"topbar_title_style": "bold",
	"topbar_link_style": "plain",
	"body_style": "plain",
	"bookmark_body_style": "plain",
	"bookmark_border_style": "plain",
	"message_style": "italic",
	"error_style": "bold",
	"command_style": "plain",
	"bottom_bar_style": "plain",
}
