package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ERRS = map[string]string{
	"ADD": "`add [target] [name...]`",
	"DELETE": "`delete [bookmark-id]`",
	"BOOKMARKS": "`bookmarks [[bookmark-id]]`",
	"CHECK": "`check [link_id]` or `check [setting]`",
	"HOME": "`home`",
	"PURGE": "`purge [host]`",
	"QUIT": "`quit`",
	"RELOAD": "`reload`",
	"SEARCH": "`search [[keyword(s)...]]`",
	"SET": "`set [setting] [value]`",
	"WRITE": "`write [target]`",
	"HELP": "`help [[topic]]`",
}

var helpRoot string = "/usr/local/share/bombadillo/help"

func helpAddress(section string) (string, error) {
	var addr string
	switch strings.ToLower(section) {
	case "add", "a", "delete", "d", "bookmarks", "bookmark", "b":
		addr = "bookmarks.help"
	case "quit", "quitting", "q", "flags", "runtime", "options", "exiting", "exit", "general", "startup", "version", "title":
		addr = "general.help"
	case "help", "info", "?", "information":
		addr = "help.help"
	case "write", "save", "saving", "w", "file", "writing", "download", "downloading", "downloads":
		addr = "saving.help"
	case "license":
		addr = "license.help"
	case "local", "file":
		addr = "local.help"
	case "finger":
		addr = "finger.help"
	case "gemini", "text/gemini", "tls", "tofu":
		addr = "gemini.help"
	case "gopher":
		addr = "gopher.help"
	case "keys", "key", "hotkeys", "hotkey", "keymap", "controls":
		addr = "keys.help"
	case "telnet":
		addr = "telnet.help"
	case "navigating", "navigation", "scroll", "scrolling", "history","links", "link":
		addr = "navigation.help"
	case "command", "commands", "functions":
		addr = "commands.help"
	case "protocol", "protocols":
		addr = "protocols.help"
	case "resources", "links":
		addr = "resources.go"
	default:
		return "", fmt.Errorf("No help section for %q exists", section)
	}

	fp := filepath.Join(helpRoot, addr)

	_, err := os.Stat(fp)

	if err != nil {
		return "", fmt.Errorf("No help section for %q exists", section)
	}

	return fp, nil
}
