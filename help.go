package main

import (
	"fmt"
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
	default:
		return "", fmt.Errorf("No help section for %q exists", section)
	}
	return filepath.Join(helpRoot, addr), nil
}
