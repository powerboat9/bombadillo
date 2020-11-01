package main

// ERRS maps commands to their syntax error message
var ERRS = map[string]string{
	"A":         "`a [target] [name...]`",
	"ADD":       "`add [target] [name...]`",
	"D":         "`d [bookmark-id]`",
	"DELETE":    "`delete [bookmark-id]`",
	"B":         "`b [[bookmark-id]]`",
	"BOOKMARKS": "`bookmarks [[bookmark-id]]`",
	"C":         "`c [link_id]` or `c [setting]`",
	"CHECK":     "`check [link_id]` or `check [setting]`",
	"H":         "`h`",
	"HOME":      "`home`",
	"P":         "`p [host]`",
	"PURGE":     "`purge [host]`",
	"Q":         "`q`",
	"QUIT":      "`quit`",
	"R":         "`r`",
	"RELOAD":    "`reload`",
	"SEARCH":    "`search [[keyword(s)...]]`",
	"S":         "`s [setting] [value]`",
	"SET":       "`set [setting] [value]`",
	"W":         "`w [target]`",
	"WRITE":     "`write [target]`",
	"?":         "`? [[command]]`",
	"HELP":      "`help [[command]]`",
}
