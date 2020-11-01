package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tildegit.org/sloum/bombadillo/cmdparse"
	"tildegit.org/sloum/bombadillo/cui"
	"tildegit.org/sloum/bombadillo/finger"
	"tildegit.org/sloum/bombadillo/gemini"
	"tildegit.org/sloum/bombadillo/gopher"
	"tildegit.org/sloum/bombadillo/http"
	"tildegit.org/sloum/bombadillo/local"
	"tildegit.org/sloum/bombadillo/telnet"
	"tildegit.org/sloum/bombadillo/termios"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type client struct {
	Height       int
	Width        int
	Options      map[string]string
	Message      string
	MessageIsErr bool
	PageState    Pages
	BookMarks    Bookmarks
	TopBar       Headbar
	FootBar      Footbar
	Certs        gemini.TofuDigest
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (c *client) GetSizeOnce() {
	var w, h = termios.GetWindowSize()
	c.Height = h
	c.Width = w
}

func (c *client) GetSize() {
	c.GetSizeOnce()
	c.SetMessage("Loading...", false)
	c.Draw()

	for {
		var w, h = termios.GetWindowSize()
		if h != c.Height || w != c.Width {
			c.Height = h
			c.Width = w
			c.SetPercentRead()
			c.Draw()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (c *client) Draw() {
	var screen strings.Builder
	screen.Grow(c.Height*c.Width + c.Width)
	screen.WriteString("\033[0m")
	screen.WriteString(c.TopBar.Render(c.Width, c.Options["theme"]))
	screen.WriteString("\n")
	pageContent := c.PageState.Render(c.Height, c.Width-1, (c.Options["theme"] == "color"))
	var re *regexp.Regexp
	if c.Options["theme"] == "inverse" {
		screen.WriteString("\033[7m")
	}
	re = regexp.MustCompile(`\033\[(?:\d*;?)+[A-Za-z]`)
	if c.BookMarks.IsOpen {
		bm := c.BookMarks.Render(c.Width, c.Height)
		bmWidth := len([]rune(bm[0]))
		for i := 0; i < c.Height-3; i++ {
			if c.Width > bmWidth {
				contentWidth := c.Width - bmWidth
				if i < len(pageContent) {
					extra := 0
					if c.Options["theme"] == "color" {
						escapes := re.FindAllString(pageContent[i], -1)
						for _, esc := range escapes {
							extra += len(esc)
						}
					}
					screen.WriteString(fmt.Sprintf("%-*.*s", contentWidth+extra, contentWidth+extra, pageContent[i]))
				} else {
					screen.WriteString(fmt.Sprintf("%-*.*s", contentWidth, contentWidth, " "))
				}
				screen.WriteString("\033[500C\033[39D")
			}

			if c.Options["theme"] == "inverse" && !c.BookMarks.IsFocused {
				screen.WriteString("\033[2;7m")
			} else if !c.BookMarks.IsFocused {
				screen.WriteString("\033[2m")
			}

			if c.Options["theme"] == "color" {
				screen.WriteString("\033[0m")
			}
			screen.WriteString(bm[i])

			if c.Options["theme"] == "inverse" && !c.BookMarks.IsFocused {
				screen.WriteString("\033[7;22m")
			} else if !c.BookMarks.IsFocused {
				screen.WriteString("\033[0m")
			}

			screen.WriteString("\n")
		}
	} else {
		for i := 0; i < c.Height-3; i++ {
			if i < len(pageContent) {
				screen.WriteString("\033[0K")
				screen.WriteString(pageContent[i])
				screen.WriteString("\n")
			} else {
				screen.WriteString("\033[0K")
				screen.WriteString("\n")
			}
		}
	}
	screen.WriteString("\033[0m")
	// TODO using message here breaks on resize, must regenerate
	screen.WriteString(c.RenderMessage())
	screen.WriteString("\n") // for the input line
	screen.WriteString(c.FootBar.Render(c.Width, c.PageState.Position, c.Options["theme"]))
	// cui.Clear("screen")
	cui.MoveCursorTo(0, 0)
	fmt.Print(screen.String())
}

func (c *client) TakeControlInput() {
	input := cui.Getch()

	switch input {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		if input == '0' {
			c.goToLink("10")
		} else {
			c.goToLink(string(input))
		}
	case 'j':
		// scroll down one line
		c.ClearMessage()
		c.Scroll(1)
	case 'k':
		// scroll up one line
		c.ClearMessage()
		c.Scroll(-1)
	case 'q':
		// quit bombadillo
		cui.Exit(0, "")
	case 'g':
		// scroll to top
		c.ClearMessage()
		c.Scroll(-len(c.PageState.History[c.PageState.Position].WrappedContent))
	case 'G':
		// scroll to bottom
		c.ClearMessage()
		c.Scroll(len(c.PageState.History[c.PageState.Position].WrappedContent))
	case 'd':
		// scroll down 75%
		c.ClearMessage()
		distance := c.Height - c.Height/4
		c.Scroll(distance)
	case 'u':
		// scroll up 75%
		c.ClearMessage()
		distance := c.Height - c.Height/4
		c.Scroll(-distance)
	case 'b', 'h':
		// go back
		c.ClearMessage()
		err := c.PageState.NavigateHistory(-1)
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		} else {
			c.SetHeaderUrl()
			c.SetPercentRead()
			c.Draw()
		}
	case 'R':
		c.ClearMessage()
		err := c.ReloadPage()
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		} else {
			c.Draw()
		}
	case 'B':
		// open the bookmarks browser
		c.BookMarks.ToggleOpen()
		c.Draw()
	case 'f', 'l':
		// go forward
		c.ClearMessage()
		err := c.PageState.NavigateHistory(1)
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		} else {
			c.SetHeaderUrl()
			c.SetPercentRead()
			c.Draw()
		}
	case '\t':
		// Toggle bookmark browser focus on/off
		c.BookMarks.ToggleFocused()
		c.Draw()
	case 'n':
		// Next search item
		c.ClearMessage()
		err := c.NextSearchItem(1)
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		}
	case 'N':
		// Previous search item
		c.ClearMessage()
		err := c.NextSearchItem(-1)
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		}
	case '/':
		// Search for text
		c.ClearMessage()
		c.ClearMessageLine()
		if c.Options["theme"] == "normal" || c.Options["theme"] == "color" {
			fmt.Printf("\033[7m%*.*s\r", c.Width, c.Width, "")
		}
		entry, err := cui.GetLine("/")
		c.ClearMessageLine()
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			break
		}
		err = c.find(entry)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
		}
		err = c.NextSearchItem(0)
		if err != nil {
			c.Draw()
		}
	case ':', ' ':
		// Process a command
		c.ClearMessage()
		c.ClearMessageLine()
		if c.Options["theme"] == "normal" || c.Options["theme"] == "color" {
			fmt.Printf("\033[7m%*.*s\r", c.Width, c.Width, "")
		}
		entry, err := cui.GetLine(": ")
		c.ClearMessageLine()
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			break
		} else if strings.TrimSpace(entry) == "" {
			c.DrawMessage()
			break
		}

		parser := cmdparse.NewParser(strings.NewReader(entry))
		p, err := parser.Parse()
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
		} else {
			err := c.routeCommandInput(p)
			if err != nil {
				c.SetMessage(err.Error(), true)
				c.DrawMessage()
			}
		}
	}
}

func (c *client) routeCommandInput(com *cmdparse.Command) error {
	switch com.Type {
	case cmdparse.SIMPLE:
		c.simpleCommand(com.Action)
	case cmdparse.GOURL:
		c.goToURL(com.Target)
	case cmdparse.GOLINK:
		c.goToLink(com.Target)
	case cmdparse.DO:
		c.doCommand(com.Action, com.Value)
	case cmdparse.DOLINK:
		c.doLinkCommand(com.Action, com.Target)
	case cmdparse.DOAS:
		c.doCommandAs(com.Action, com.Value)
	case cmdparse.DOLINKAS:
		c.doLinkCommandAs(com.Action, com.Target, com.Value)
	default:
		return fmt.Errorf("Unknown command entry")
	}

	return nil

}

func (c *client) simpleCommand(action string) {
	action = strings.ToUpper(action)
	switch action {
	case "Q", "QUIT":
		cui.Exit(0, "")
	case "H", "HOME":
		if c.Options["homeurl"] != "unset" {
			go c.Visit(c.Options["homeurl"])
		} else {
			c.SetMessage(fmt.Sprintf("No home address has been set"), false)
			c.DrawMessage()
		}
	case "B", "BOOKMARKS":
		c.BookMarks.ToggleOpen()
		c.Draw()
	case "R", "RELOAD":
		c.ClearMessage()
		err := c.ReloadPage()
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		} else {
			c.Draw()
		}
	case "SEARCH":
		c.search("", "", "?")
	case "HELP", "?":
		c.Visit(helplocation)
	default:
		c.SetMessage(syntaxErrorMessage(action), true)
		c.DrawMessage()
	}
}

func (c *client) doCommand(action string, values []string) {
	switch action {
	case "C", "CHECK":
		c.displayConfigValue(values[0])
		c.DrawMessage()
	case "HELP", "?":
		if val, ok := ERRS[values[0]]; ok {
			c.SetMessage(val, false)
		} else {
			msg := fmt.Sprintf("%q is not a valid command; help syntax: %s", values[0], ERRS[action])
			c.SetMessage(msg, false)
		}
		c.DrawMessage()
	case "PURGE", "P":
		err := c.Certs.Purge(values[0])
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		if values[0] == "*" {
			c.SetMessage("All certificates have been purged", false)
			c.DrawMessage()
		} else {
			c.SetMessage(fmt.Sprintf("The certificate for %q has been purged", strings.ToLower(values[0])), false)
			c.DrawMessage()
		}
		err = saveConfig()
		if err != nil {
			c.SetMessage("Error saving purge to file", true)
			c.DrawMessage()
		}
	case "SEARCH":
		c.search(values[0], "", "")
	case "WRITE", "W":
		if values[0] == "." {
			values[0] = c.PageState.History[c.PageState.Position].Location.Full
		}
		u, err := MakeUrl(values[0])
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		fns := strings.Split(u.Resource, "/")
		var fn string
		if len(fns) > 0 {
			fn = strings.Trim(fns[len(fns)-1], "\t\r\n \a\f\v")
		} else {
			fn = "index"
		}
		if fn == "" {
			fn = "index"
		}
		c.saveFile(u, fn)

	default:
		c.SetMessage(syntaxErrorMessage(action), true)
		c.DrawMessage()
	}
}

func (c *client) doCommandAs(action string, values []string) {
	switch action {
	case "ADD", "A":
		if len(values) < 2 {
			c.SetMessage(syntaxErrorMessage(action), true)
			c.DrawMessage()
			return
		}
		if values[0] == "." {
			values[0] = c.PageState.History[c.PageState.Position].Location.Full
		}
		msg, err := c.BookMarks.Add(values)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		c.SetMessage(msg, false)
		c.DrawMessage()

		err = saveConfig()
		if err != nil {
			c.SetMessage("Error saving bookmark to file", true)
			c.DrawMessage()
		}
		if c.BookMarks.IsOpen {
			c.Draw()
		}
	case "SEARCH":
		if len(values) < 2 {
			c.SetMessage(syntaxErrorMessage(action), true)
			c.DrawMessage()
			return
		}
		c.search(strings.Join(values, " "), "", "")
	case "SET", "S":
		if len(values) < 2 {
			c.SetMessage(syntaxErrorMessage(action), true)
			c.DrawMessage()
			return
		}
		if _, ok := c.Options[values[0]]; ok {
			val := strings.Join(values[1:], " ")
			if !validateOpt(values[0], val) {
				c.SetMessage(fmt.Sprintf("Invalid setting for %q", values[0]), true)
				c.DrawMessage()
				return
			}
			c.Options[values[0]] = lowerCaseOpt(values[0], val)
			if values[0] == "geminiblocks" {
				gemini.BlockBehavior = c.Options[values[0]]
			} else if values[0] == "timeout" {
				updateTimeouts(c.Options[values[0]])
			} else if values[0] == "configlocation" {
				c.SetMessage("Cannot set READ ONLY setting 'configlocation'", true)
				c.DrawMessage()
				return
			}
			err := saveConfig()
			if err != nil {
				c.SetMessage("Value set, but error saving config to file", true)
				c.DrawMessage()
			} else {
				c.SetMessage(fmt.Sprintf("%s is now set to %q", values[0], c.Options[values[0]]), false)
				c.Draw()
			}
			return
		}
		c.SetMessage(fmt.Sprintf("Unable to set %s, it does not exist", values[0]), true)
		c.DrawMessage()
	default:
		c.SetMessage(syntaxErrorMessage(action), true)
		c.DrawMessage()
	}
}

func (c *client) doLinkCommandAs(action, target string, values []string) {
	num, err := strconv.Atoi(target)
	if err != nil {
		c.SetMessage(fmt.Sprintf("Expected link number, got %q", target), true)
		c.DrawMessage()
		return
	}

	num--

	links := c.PageState.History[c.PageState.Position].Links
	if num >= len(links) || num < 0 {
		c.SetMessage(fmt.Sprintf("Invalid link id: %s", target), true)
		c.DrawMessage()
		return
	}

	switch action {
	case "ADD", "A":
		bm := make([]string, 0, 5)
		bm = append(bm, links[num])
		bm = append(bm, values...)
		msg, err := c.BookMarks.Add(bm)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		c.SetMessage(msg, false)
		c.DrawMessage()

		err = saveConfig()
		if err != nil {
			c.SetMessage("Error saving bookmark to file", true)
			c.DrawMessage()
		}
		if c.BookMarks.IsOpen {
			c.Draw()
		}
	case "WRITE", "W":
		out := make([]string, 0, len(values)+1)
		out = append(out, links[num])
		out = append(out, values...)
		c.doCommandAs(action, out)
	default:
		c.SetMessage(syntaxErrorMessage(action), true)
		c.DrawMessage()
	}
}

func (c *client) saveFile(u Url, name string) {
	var file []byte
	var err error
	c.SetMessage(fmt.Sprintf("Saving %s ...", name), false)
	c.DrawMessage()
	switch u.Scheme {
	case "gopher":
		file, err = gopher.Retrieve(u.Host, u.Port, u.Resource)
	case "gemini":
		file, err = gemini.Fetch(u.Host, u.Port, u.Resource, &c.Certs)
	case "http", "https":
		file, err = http.Fetch(u.Full)
	default:
		c.SetMessage(fmt.Sprintf("Saving files over %s is not supported", u.Scheme), true)
		c.DrawMessage()
		return
	}

	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}

	// We are ignoring the error here since WriteFile will
	// generate the same error, and will handle the messaging
	savePath, _ := findAvailableFileName(c.Options["savelocation"], name)
	err = ioutil.WriteFile(savePath, file, 0644)
	if err != nil {
		c.SetMessage("Error writing file: "+err.Error(), true)
		c.DrawMessage()
		return
	}

	c.SetMessage(fmt.Sprintf("File saved to: %s", savePath), false)
	c.DrawMessage()
}

func (c *client) saveFileFromData(d, name string) {
	data := []byte(d)
	c.SetMessage(fmt.Sprintf("Saving %s ...", name), false)
	c.DrawMessage()

	// We are ignoring the error here since WriteFile will
	// generate the same error, and will handle the messaging
	savePath, _ := findAvailableFileName(c.Options["savelocation"], name)
	err := ioutil.WriteFile(savePath, data, 0644)
	if err != nil {
		c.SetMessage("Error writing file: "+err.Error(), true)
		c.DrawMessage()
		return
	}

	c.SetMessage(fmt.Sprintf("File saved to: %s", savePath), false)
	c.DrawMessage()
}

func (c *client) doLinkCommand(action, target string) {
	num, err := strconv.Atoi(target)
	if err != nil {
		c.SetMessage(fmt.Sprintf("Expected number, got %q", target), true)
		c.DrawMessage()
	}

	switch action {
	case "DELETE", "D":
		msg, err := c.BookMarks.Delete(num)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		c.SetMessage(msg, false)
		c.DrawMessage()

		err = saveConfig()
		if err != nil {
			c.SetMessage("Error saving bookmark deletion to file", true)
			c.DrawMessage()
		}
		if c.BookMarks.IsOpen {
			c.Draw()
		}
	case "BOOKMARKS", "B":
		if num > len(c.BookMarks.Links)-1 {
			c.SetMessage(fmt.Sprintf("There is no bookmark with ID %d", num), true)
			c.DrawMessage()
			return
		}
		c.Visit(c.BookMarks.Links[num])
	case "CHECK", "C":
		num--

		links := c.PageState.History[c.PageState.Position].Links
		if num >= len(links) || num < 0 {
			c.SetMessage(fmt.Sprintf("Invalid link id: %s", target), true)
			c.DrawMessage()
			return
		}
		link := links[num]
		c.SetMessage(fmt.Sprintf("[%d] %s", num+1, link), false)
		c.DrawMessage()
	case "WRITE", "W":
		links := c.PageState.History[c.PageState.Position].Links
		if len(links) < num || num < 1 {
			c.SetMessage("Invalid link ID", true)
			c.DrawMessage()
			return
		}
		u, err := MakeUrl(links[num-1])
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		fns := strings.Split(u.Resource, "/")
		var fn string
		if len(fns) > 0 {
			fn = strings.Trim(fns[len(fns)-1], "\t\r\n \a\f\v")
		} else {
			fn = "index"
		}
		if fn == "" {
			fn = "index"
		}
		c.saveFile(u, fn)
	default:
		c.SetMessage(syntaxErrorMessage(action), true)
		c.DrawMessage()
	}

}

func (c *client) search(query, uri, question string) {
	var entry string
	var err error
	if query == "" {
		c.ClearMessage()
		c.ClearMessageLine()
		if c.Options["theme"] == "normal" || c.Options["theme"] == "color" {
			fmt.Printf("\033[7m%*.*s\r", c.Width, c.Width, "")
		}
		fmt.Print(question)
		entry, err = cui.GetLine("? ")
		c.ClearMessageLine()
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		} else if strings.TrimSpace(entry) == "" {
			c.ClearMessage()
			c.DrawMessage()
			return
		}
	} else {
		entry = query
	}
	if uri == "" {
		uri = c.Options["searchengine"]
	}
	u, err := MakeUrl(uri)
	if err != nil {
		c.SetMessage("The search url is not valid", true)
		c.DrawMessage()
		return
	}
	var rootUrl string
	switch u.Scheme {
	case "gopher":
		if ind := strings.Index(u.Full, "\t"); ind >= 0 {
			rootUrl = u.Full[:ind]
		} else {
			rootUrl = u.Full
		}
		c.Visit(fmt.Sprintf("%s\t%s", rootUrl, entry))
	case "gemini":
		if ind := strings.Index(u.Full, "?"); ind >= 0 {
			rootUrl = u.Full[:ind]
		} else {
			rootUrl = u.Full
		}
		escapedEntry := url.PathEscape(entry)
		c.Visit(fmt.Sprintf("%s?%s", rootUrl, escapedEntry))
	case "http", "https":
		c.Visit(u.Full)
	default:
		c.SetMessage(fmt.Sprintf("%q is not a supported protocol", u.Scheme), true)
		c.DrawMessage()
	}
}

func (c *client) Scroll(amount int) {
	if c.BookMarks.IsFocused {
		bottom := len(c.BookMarks.Titles) - c.Height + 5 // 3 for the three bars: top, msg, bottom
		if amount < 0 && c.BookMarks.Position == 0 {
			c.SetMessage("The bookmark ladder does not go up any further", false)
			c.DrawMessage()
			fmt.Print("\a")
			return
		} else if (amount > 0 && c.BookMarks.Position == bottom) || bottom < 0 {
			c.SetMessage("Feel the ground beneath your bookmarks", false)
			c.DrawMessage()
			fmt.Print("\a")
			return
		}

		newScrollPosition := c.BookMarks.Position + amount
		if newScrollPosition < 0 {
			newScrollPosition = 0
		} else if newScrollPosition > bottom {
			newScrollPosition = bottom
		}

		c.BookMarks.Position = newScrollPosition
	} else {
		var percentRead int
		page := c.PageState.History[c.PageState.Position]
		bottom := len(page.WrappedContent) - c.Height + 3 // 3 for the three bars: top, msg, bottom
		if amount < 0 && page.ScrollPosition == 0 {
			c.SetMessage("You are already at the top", false)
			c.DrawMessage()
			fmt.Print("\a")
			return
		} else if (amount > 0 && page.ScrollPosition == bottom) || bottom < 0 {
			c.FootBar.SetPercentRead(100)
			c.SetMessage("You are already at the bottom", false)
			c.DrawMessage()
			fmt.Print("\a")
			return
		}

		newScrollPosition := page.ScrollPosition + amount
		if newScrollPosition < 0 {
			newScrollPosition = 0
		} else if newScrollPosition > bottom {
			newScrollPosition = bottom
		}

		c.PageState.History[c.PageState.Position].ScrollPosition = newScrollPosition

		if len(page.WrappedContent) < c.Height-3 {
			percentRead = 100
		} else {
			percentRead = int(float32(newScrollPosition+c.Height-3) / float32(len(page.WrappedContent)) * 100.0)
		}
		c.FootBar.SetPercentRead(percentRead)
	}
	c.Draw()
}

func (c *client) ReloadPage() error {
	if c.PageState.Length < 1 {
		return fmt.Errorf("There is no page to reload")
	}
	url := c.PageState.History[c.PageState.Position].Location.Full
	if c.PageState.Position == 0 {
		c.PageState.Position--
	} else {
		err := c.PageState.NavigateHistory(-1)
		if err != nil {
			return err
		}
	}
	length := c.PageState.Length
	c.Visit(url)
	c.PageState.Length = length
	return nil
}

func (c *client) SetPercentRead() {
	page := c.PageState.History[c.PageState.Position]
	var percentRead int
	if len(page.WrappedContent) < c.Height-3 {
		percentRead = 100
	} else {
		percentRead = int(float32(page.ScrollPosition+c.Height-3) / float32(len(page.WrappedContent)) * 100.0)
	}
	c.FootBar.SetPercentRead(percentRead)
}

func (c *client) displayConfigValue(setting string) {
	if val, ok := c.Options[setting]; ok {
		c.SetMessage(fmt.Sprintf("%s is set to: %q", setting, val), false)
		c.DrawMessage()
	} else {
		c.SetMessage(fmt.Sprintf("Invalid: %q does not exist", setting), true)
		c.DrawMessage()
	}
}

func (c *client) SetMessage(msg string, isError bool) {
	c.MessageIsErr = isError
	c.Message = strings.Replace(msg, "\t", "%09", -1)
}

func (c *client) DrawMessage() {
	cui.MoveCursorTo(c.Height-1, 0)
	fmt.Print(c.RenderMessage())
}

func (c *client) RenderMessage() string {
	leadIn, leadOut := "", ""
	if c.Options["theme"] == "normal" || c.Options["theme"] == "color" {
		leadIn = "\033[7m"
		leadOut = "\033[0m"
	}

	if c.MessageIsErr {
		leadIn = "\033[31;1m"
		leadOut = "\033[0m"

		if c.Options["theme"] == "normal" || c.Options["theme"] == "color" {
			leadIn = "\033[41;1;7m"
		}
	}

	return fmt.Sprintf("%s%-*.*s%s", leadIn, c.Width, c.Width, c.Message, leadOut)
}

func (c *client) ClearMessage() {
	c.SetMessage("", false)
}

func (c *client) ClearMessageLine() {
	cui.MoveCursorTo(c.Height-1, 0)
	cui.Clear("line")
}

func (c *client) goToURL(u string) {
	if num, _ := regexp.MatchString(`^-?\d+.?\d*$`, u); num {
		c.goToLink(u)
		return
	}

	c.Visit(u)
}

func (c *client) goToLink(l string) {
	if num, _ := regexp.MatchString(`^-?\d+$`, l); num && c.PageState.Length > 0 {
		linkcount := len(c.PageState.History[c.PageState.Position].Links)
		item, err := strconv.Atoi(l)
		if err != nil {
			c.SetMessage(fmt.Sprintf("Invalid link id: %s", l), true)
			c.DrawMessage()
			return
		}
		if item <= linkcount && item > 0 {
			linkurl := c.PageState.History[c.PageState.Position].Links[item-1]
			c.Visit(linkurl)
		} else {
			c.SetMessage(fmt.Sprintf("Invalid link id: %s", l), true)
			c.DrawMessage()
			return
		}
	}
}

func (c *client) SetHeaderUrl() {
	if c.PageState.Length > 0 {
		u := c.PageState.History[c.PageState.Position].Location.Full
		c.TopBar.url = strings.Replace(u, "\t", "%09", -1)
	} else {
		c.TopBar.url = ""
	}
}

// Visit functions as a controller/router to the
// appropriate protocol handler
func (c *client) Visit(url string) {
	c.SetMessage("Loading...", false)
	c.DrawMessage()

	url = strings.Replace(url, "%09", "\t", -1)
	u, err := MakeUrl(url)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}

	switch u.Scheme {
	case "gopher":
		c.handleGopher(u)
	case "gemini":
		c.handleGemini(u)
	case "telnet":
		c.handleTelnet(u)
	case "http", "https":
		c.handleWeb(u)
	case "local":
		c.handleLocal(u)
	case "finger":
		c.handleFinger(u)
	default:
		c.SetMessage(fmt.Sprintf("%q is not a supported protocol", u.Scheme), true)
		c.DrawMessage()
	}
}

// +++ Begin Protocol Handlers +++

func (c *client) handleGopher(u Url) {
	if u.DownloadOnly || (c.Options["showimages"] == "false" && (u.Mime == "I" || u.Mime == "g")) {
		nameSplit := strings.Split(u.Resource, "/")
		filename := nameSplit[len(nameSplit)-1]
		filename = strings.Trim(filename, " \t\r\n\v\f\a")
		if filename == "" {
			filename = "gopherfile"
		}
		c.saveFile(u, filename)
	} else if u.Mime == "7" {
		c.search("", u.Full, "?")
	} else {
		content, links, err := gopher.Visit(u.Mime, u.Host, u.Port, u.Resource)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		pg := MakePage(u, content, links)
		if u.Mime == "I" || u.Mime == "g" {
			pg.FileType = "image"
		} else {
			pg.FileType = "text"
		}
		pg.WrapContent(c.Width-1, (c.Options["theme"] == "color"))
		c.PageState.Add(pg)
		c.SetPercentRead()
		c.ClearMessage()
		c.SetHeaderUrl()
		c.Draw()
	}
}

func (c *client) handleGemini(u Url) {
	capsule, err := gemini.Visit(u.Host, u.Port, u.Resource, &c.Certs)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}
	go saveConfig()
	switch capsule.Status {
	case 1:
		// Query
		c.search("", u.Full, capsule.Content)
	case 2:
		// Success
		if capsule.MimeMaj == "text" || (c.Options["showimages"] == "true" && capsule.MimeMaj == "image") {
			pg := MakePage(u, capsule.Content, capsule.Links)
			pg.FileType = capsule.MimeMaj
			pg.WrapContent(c.Width-1, (c.Options["theme"] == "color"))
			c.PageState.Add(pg)
			c.SetPercentRead()
			c.ClearMessage()
			c.SetHeaderUrl()
			c.Draw()
		} else {
			c.SetMessage("The file is non-text: writing to disk...", false)
			c.DrawMessage()
			nameSplit := strings.Split(u.Resource, "/")
			filename := nameSplit[len(nameSplit)-1]
			c.saveFileFromData(capsule.Content, filename)
		}
	case 3:
		// Redirect
		lowerRedirect := strings.ToLower(capsule.Content)
		lowerOriginal := strings.ToLower(u.Full)
		if strings.Replace(lowerRedirect, lowerOriginal, "", 1) == "/" {
			c.Visit(capsule.Content)
		} else {
			if !strings.Contains(capsule.Content, "://") {
				lnk, lnkErr := gemini.HandleRelativeUrl(capsule.Content, u.Full)
				if lnkErr == nil {
					capsule.Content = lnk
				}
			}

			c.SetMessage(fmt.Sprintf("Follow redirect (y/n): %s?", capsule.Content), false)
			c.DrawMessage()
			ch := cui.Getch()
			if ch == 'y' || ch == 'Y' {
				c.Visit(capsule.Content)
			} else {
				c.SetMessage("Redirect aborted", false)
				c.DrawMessage()
			}
		}
	}
}

func (c *client) handleTelnet(u Url) {
	c.SetMessage("Attempting to start telnet session", false)
	c.DrawMessage()
	msg, err := telnet.StartSession(u.Host, u.Port)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
	} else {
		c.SetMessage(msg, true)
		c.DrawMessage()
	}
	c.Draw()
}

func (c *client) handleLocal(u Url) {
	content, links, err := local.Open(u.Resource)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}
	pg := MakePage(u, content, links)
	ext := strings.ToLower(filepath.Ext(u.Full))
	if ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".png" {
		pg.FileType = "image"
	}
	pg.WrapContent(c.Width-1, (c.Options["theme"] == "color"))
	c.PageState.Add(pg)
	c.SetPercentRead()
	c.ClearMessage()
	c.SetHeaderUrl()
	c.Draw()
}

func (c *client) handleFinger(u Url) {
	content, err := finger.Finger(u.Host, u.Port, u.Resource)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}
	pg := MakePage(u, content, []string{})
	pg.WrapContent(c.Width-1, (c.Options["theme"] == "color"))
	c.PageState.Add(pg)
	c.SetPercentRead()
	c.ClearMessage()
	c.SetHeaderUrl()
	c.Draw()
}

func (c *client) handleWeb(u Url) {
	wm := strings.ToLower(c.Options["webmode"])
	switch wm {
	case "lynx", "w3m", "elinks":
		if http.IsTextFile(u.Full) {
			page, err := http.Visit(wm, u.Full, c.Width-1)
			if err != nil {
				c.SetMessage(fmt.Sprintf("%s error: %s", wm, err.Error()), true)
				c.DrawMessage()
				return
			}
			pg := MakePage(u, page.Content, page.Links)
			pg.WrapContent(c.Width-1, (c.Options["theme"] == "color"))
			c.PageState.Add(pg)
			c.SetPercentRead()
			c.ClearMessage()
			c.SetHeaderUrl()
			c.Draw()
		} else {
			c.SetMessage("The file is non-text: writing to disk...", false)
			c.DrawMessage()
			var fn string
			if i := strings.LastIndex(u.Full, "/"); i > 0 && i+1 < len(u.Full) {
				fn = u.Full[i+1:]
			} else {
				fn = "bombadillo.download"
			}
			c.saveFile(u, fn)
		}
	case "gui":
		c.SetMessage("Attempting to open in gui web browser", false)
		c.DrawMessage()
		msg, err := http.OpenInBrowser(u.Full)
		if err != nil {
			c.SetMessage(err.Error(), true)
		} else {
			c.SetMessage(msg, false)
		}
		c.DrawMessage()
	default:
		c.SetMessage("Current 'webmode' setting does not allow http/https", false)
		c.DrawMessage()
	}
}

func (c *client) find(s string) error {
	c.PageState.History[c.PageState.Position].SearchTerm = s
	c.PageState.History[c.PageState.Position].FindText()
	if s == "" {
		return nil
	}
	if len(c.PageState.History[c.PageState.Position].FoundLinkLines) == 0 {
		return fmt.Errorf("No text matching %q was found", s)
	}
	return nil
}

func (c *client) NextSearchItem(dir int) error {
	page := c.PageState.History[c.PageState.Position]
	if len(page.FoundLinkLines) == 0 {
		return fmt.Errorf("The search is over before it has begun")
	}
	c.PageState.History[c.PageState.Position].SearchIndex += dir
	page.SearchIndex += dir
	if page.SearchIndex < 0 {
		c.PageState.History[c.PageState.Position].SearchIndex = 0
		page.SearchIndex = 0
	}

	if page.SearchIndex >= len(page.FoundLinkLines) {
		c.PageState.History[c.PageState.Position].SearchIndex = len(page.FoundLinkLines) - 1
		return fmt.Errorf("The search path goes no further")
	} else if page.SearchIndex < 0 {
		c.PageState.History[c.PageState.Position].SearchIndex = 0
		return fmt.Errorf("You are at the beginning of the search path")
	}

	diff := page.FoundLinkLines[page.SearchIndex] - page.ScrollPosition
	c.ScrollForSearch(diff)
	c.Draw()
	return nil
}

func (c *client) ScrollForSearch(amount int) {
	var percentRead int
	page := c.PageState.History[c.PageState.Position]
	bottom := len(page.WrappedContent) - c.Height + 3 // 3 for the three bars: top, msg, bottom

	newScrollPosition := page.ScrollPosition + amount
	if newScrollPosition < 0 {
		newScrollPosition = 0
	} else if newScrollPosition > bottom {
		newScrollPosition = bottom
	}

	c.PageState.History[c.PageState.Position].ScrollPosition = newScrollPosition

	if len(page.WrappedContent) < c.Height-3 {
		percentRead = 100
	} else {
		percentRead = int(float32(newScrollPosition+c.Height-3) / float32(len(page.WrappedContent)) * 100.0)
	}
	c.FootBar.SetPercentRead(percentRead)
	c.Draw()
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

// MakeClient returns a client struct and names the client after
// the string that is passed in
func MakeClient(name string) *client {
	c := client{0, 0, defaultOptions, "", false, MakePages(), MakeBookmarks(), MakeHeadbar(name), MakeFootbar(), gemini.MakeTofuDigest()}
	return &c
}

func findAvailableFileName(fpath, fname string) (string, error) {
	savePath := filepath.Join(fpath, fname)
	_, fileErr := os.Stat(savePath)

	for suffix := 1; fileErr == nil; suffix++ {
		fn := fmt.Sprintf("%s.%d", fname, suffix)
		savePath = filepath.Join(fpath, fn)
		_, fileErr = os.Stat(savePath)

		if !os.IsNotExist(fileErr) && fileErr != nil {
			return savePath, fileErr
		}
	}

	return savePath, nil
}

func syntaxErrorMessage(action string) string {
	if val, ok := ERRS[action]; ok {
		return fmt.Sprintf("Incorrect syntax. Try: %s", val)
	}
	return fmt.Sprintf("Unknown command %q", action)
}

func updateTimeouts(timeoutString string) error {
	sec, err := strconv.Atoi(timeoutString)
	if err != nil {
		return err
	}
	timeout := time.Duration(sec) * time.Second

	gopher.Timeout = timeout
	gemini.TlsTimeout = timeout

	return nil
}
