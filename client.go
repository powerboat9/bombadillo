package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tildegit.org/sloum/bombadillo/cmdparse"
	"tildegit.org/sloum/bombadillo/cui"
	"tildegit.org/sloum/bombadillo/gemini"
	"tildegit.org/sloum/bombadillo/gopher"
	"tildegit.org/sloum/bombadillo/http"
	"tildegit.org/sloum/bombadillo/telnet"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type client struct {
	Height int
	Width int
	Options map[string]string
	Message string
	MessageIsErr bool
	PageState Pages
	BookMarks Bookmarks
	TopBar Headbar
	FootBar Footbar
	Certs gemini.TofuDigest
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (c *client) GetSizeOnce() {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Fatal error: Unable to retrieve terminal size")
		os.Exit(5)
	}
	var h, w int
	fmt.Sscan(string(out), &h, &w)
	c.Height = h
	c.Width = w
}

func (c *client) GetSize() {
	c.GetSizeOnce()
	c.SetMessage("Loading...", false)
	c.Draw()

	for {
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("Fatal error: Unable to retrieve terminal size")
			os.Exit(5)
		}
		var h, w int
		fmt.Sscan(string(out), &h, &w)
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
	screen.Grow(c.Height * c.Width + c.Width)
	screen.WriteString("\033[0m")
	screen.WriteString(c.TopBar.Render(c.Width, c.Options["theme"]))
	screen.WriteString("\n")
	pageContent := c.PageState.Render(c.Height, c.Width - 1)
  if c.Options["theme"] == "inverse" {
		screen.WriteString("\033[7m")
	}
	if c.BookMarks.IsOpen {
		bm := c.BookMarks.Render(c.Width, c.Height)
		bmWidth := len([]rune(bm[0]))
		for i := 0; i < c.Height - 3; i++ {
			if c.Width > bmWidth {
				contentWidth := c.Width - bmWidth
				if i < len(pageContent)  {
					screen.WriteString(fmt.Sprintf("%-*.*s", contentWidth, contentWidth, pageContent[i]))
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
			
			screen.WriteString(bm[i])

			if c.Options["theme"] == "inverse" && !c.BookMarks.IsFocused {
				screen.WriteString("\033[7;22m")
			} else if !c.BookMarks.IsFocused {
				screen.WriteString("\033[0m")
			}

			screen.WriteString("\n")
		}
	} else {
		for i := 0; i < c.Height - 3; i++ {
			if i < len(pageContent) {
				screen.WriteString(fmt.Sprintf("%-*.*s", c.Width - 1, c.Width - 1, pageContent[i]))
				screen.WriteString("\n")
			} else {
				screen.WriteString(fmt.Sprintf("%-*.*s", c.Width, c.Width, " "))
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
	cui.MoveCursorTo(0,0)
	fmt.Print(screen.String())
}

func (c *client) TakeControlInput() {
	input := cui.Getch()

	switch input {
	case 'j', 'J':
		// scroll down one line
		c.ClearMessage()
		c.Scroll(1)
	case 'k', 'K':
		// scroll up one line
		c.ClearMessage()
		c.Scroll(-1)
	case 'q', 'Q':
		// quit bombadillo
		cui.Exit()
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
		distance := c.Height - c.Height / 4
		c.Scroll(distance)
	case 'u':
		// scroll up 75%
		c.ClearMessage()
		distance := c.Height - c.Height / 4
		c.Scroll(-distance)
	case 'b':
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
	case 'B':
		// open the bookmarks browser
		c.BookMarks.ToggleOpen()
		c.Draw()
	case 'f', 'F':
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
	case ':', ' ':
		// Process a command
		c.ClearMessage()
		c.ClearMessageLine()
		if c.Options["theme"] == "normal" {
			fmt.Printf("\033[7m%*.*s\r", c.Width, c.Width, "")
		}
		entry, err := cui.GetLine()
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
				c.Draw()
			}
		}
	}
}


func (c *client) routeCommandInput(com *cmdparse.Command) error {
	var err error
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
		return fmt.Errorf("Unknown command entry!")
	}

	return err
}

func (c *client) simpleCommand(action string) {
	action = strings.ToUpper(action)
	switch action {
	case "Q", "QUIT":
		cui.Exit()
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
	case "R", "REFRESH":
		// TODO build refresh code
	case "SEARCH":
		c.search("", "", "?")
	case "HELP", "?":
		go c.Visit(helplocation)
	default:
	  c.SetMessage(fmt.Sprintf("Unknown action %q", action), true)
		c.DrawMessage()
	}
}

func (c *client) doCommand(action string, values []string) {
	if length := len(values); length != 1 {
		c.SetMessage(fmt.Sprintf("Expected 1 argument, received %d", len(values)), true)
		c.DrawMessage()
		return
	}

	switch action {
	case "CHECK", "C":
		c.displayConfigValue(values[0])
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
			fn = strings.Trim(fns[len(fns) - 1], "\t\r\n \a\f\v")
		} else {
			fn = "index"
		}
		if fn == "" {
			fn = "index"
		}
		c.saveFile(u, fn)

	default:
	  c.SetMessage(fmt.Sprintf("Unknown action %q", action), true)
		c.DrawMessage()
	}
}

func (c *client) doCommandAs(action string, values []string) {
	if len(values) < 2 {
		c.SetMessage(fmt.Sprintf("Expected 1 argument, received %d", len(values)), true)
		c.DrawMessage()
		return
	}

	if values[0] == "." {
		values[0] = c.PageState.History[c.PageState.Position].Location.Full
	}

	switch action {
	case "ADD", "A":
		msg, err := c.BookMarks.Add(values)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		} else {
			c.SetMessage(msg, false)
			c.DrawMessage()
		}

		err = saveConfig()
		if err != nil {
			c.SetMessage("Error saving bookmark to file", true)
			c.DrawMessage()
		}
		if c.BookMarks.IsOpen {
			c.Draw()
		}
	case "SEARCH":
		c.search(strings.Join(values, " "), "", "")
	case "WRITE", "W":
		u, err := MakeUrl(values[0])
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		fileName := strings.Join(values[1:], "-")
		fileName = strings.Trim(fileName, " \t\r\n\a\f\v")
		c.saveFile(u, fileName)

	case "SET", "S":
		if _, ok := c.Options[values[0]]; ok {
			val := strings.Join(values[1:], " ")
			if !validateOpt(values[0], val) {
				c.SetMessage(fmt.Sprintf("Invalid setting for %q", values[0]), true)
				c.DrawMessage()
				return
			}
			c.Options[values[0]] = lowerCaseOpt(values[0], val)
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
		return
	}
	c.SetMessage(fmt.Sprintf("Unknown command structure"), true)
}

func (c *client) doLinkCommandAs(action, target string, values []string) {
	num, err := strconv.Atoi(target)
	if err != nil {
		c.SetMessage(fmt.Sprintf("Expected link number, got %q", target), true)
		c.DrawMessage()
		return
	}

	num -= 1

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
		} else {
			c.SetMessage(msg, false)
			c.DrawMessage()
		}

		err = saveConfig()
		if err != nil {
			c.SetMessage("Error saving bookmark to file", true)
			c.DrawMessage()
		}
		if c.BookMarks.IsOpen {
			c.Draw()
		}
	case "WRITE", "W":
		out := make([]string, 0, len(values) + 1)
		out = append(out, links[num])
		out = append(out, values...)
		c.doCommandAs(action, out)
	default:
		c.SetMessage(fmt.Sprintf("Unknown command structure"), true)
	}
}


func (c *client) getCurrentPageUrl() (string, error) {
	if c.PageState.Length < 1 {
		return "", fmt.Errorf("There are no pages in history")
	}
	return c.PageState.History[c.PageState.Position].Location.Full, nil
}

func (c *client) getCurrentPageRawData() (string, error) {
	if c.PageState.Length < 1 {
		return "", fmt.Errorf("There are no pages in history")
	}
	return c.PageState.History[c.PageState.Position].RawContent, nil
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
	savePath := c.Options["savelocation"] + name
	err = ioutil.WriteFile(savePath, file, 0644)
	if err != nil {
		c.SetMessage("Error writing file to disk", true)
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
	savePath := c.Options["savelocation"] + name
	err := ioutil.WriteFile(savePath, data, 0644)
	if err != nil {
		c.SetMessage("Error writing file to disk", true)
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
		} else {
			c.SetMessage(msg, false)
			c.DrawMessage()
		}

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
		num -= 1

		links := c.PageState.History[c.PageState.Position].Links
		if num >= len(links) || num < 1 {
			c.SetMessage(fmt.Sprintf("Invalid link id: %s", target), true)
			c.DrawMessage()
			return
		}
		link := links[num]
		c.SetMessage(fmt.Sprintf("[%d] %s", num + 1, link), false)
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
			fn = strings.Trim(fns[len(fns) - 1], "\t\r\n \a\f\v")
		} else {
			fn = "index"
		}
		if fn == "" {
			fn = "index"
		}
		c.saveFile(u, fn)
	default:
	  c.SetMessage(fmt.Sprintf("Action %q does not exist for target %q", action, target), true)
		c.DrawMessage()
	}

}

func (c *client) search(query, url, question string) {
	var entry string
	var err error
	if query == "" {
		c.ClearMessage()
		c.ClearMessageLine()
		if c.Options["theme"] == "normal" {
			fmt.Printf("\033[7m%*.*s\r", c.Width, c.Width, "")
		}
		fmt.Print(question)
		entry, err = cui.GetLine()
		c.ClearMessageLine()
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		} else if strings.TrimSpace(entry) == "" {
			return
		}
	} else {
		entry = query
	}
	if url == "" {
		url = c.Options["searchengine"]
	}
	u, err := MakeUrl(url)
	if err != nil {
		c.SetMessage("The search url is not a valid url", true)
		c.DrawMessage()
		return
	}
	switch u.Scheme {
	case "gopher":
		go c.Visit(fmt.Sprintf("%s\t%s",u.Full,entry))
	case "gemini":
		// TODO url escape the entry variable
		escapedEntry := entry
		go c.Visit(fmt.Sprintf("%s?%s",u.Full,escapedEntry))
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

		if len(page.WrappedContent) < c.Height - 3 {
			percentRead = 100
		} else {
			percentRead = int(float32(newScrollPosition + c.Height - 3) / float32(len(page.WrappedContent)) * 100.0)
		}
		c.FootBar.SetPercentRead(percentRead)
	}
	c.Draw()
}

func (c *client) SetPercentRead() {
	page := c.PageState.History[c.PageState.Position]
	var percentRead int
	if len(page.WrappedContent) < c.Height - 3 {
		percentRead = 100
	} else {
		percentRead = int(float32(page.ScrollPosition + c.Height - 3) / float32(len(page.WrappedContent)) * 100.0)
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
	c.Message = strings.ReplaceAll(msg, "\t", "%09")
}

func (c *client) DrawMessage() {
	cui.MoveCursorTo(c.Height-1, 0)
	fmt.Print(c.RenderMessage())
}

func (c *client) RenderMessage() string {
	leadIn, leadOut := "", ""
	if c.Options["theme"] == "normal" {
		leadIn = "\033[7m"
		leadOut = "\033[0m"
	}

	if c.MessageIsErr {
		leadIn = "\033[31;1m"
		leadOut = "\033[0m"

		if c.Options["theme"] == "normal" {
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

	go c.Visit(u)
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
		c.TopBar.url = strings.ReplaceAll(u, "\t", "%09")
	} else {
		c.TopBar.url = ""
	}
}

func (c *client) Visit(url string) {
	c.SetMessage("Loading...", false)
	c.DrawMessage()

	url = strings.ReplaceAll(url, "%09", "\t")
	u, err := MakeUrl(url)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}

	switch u.Scheme {
	case "gopher":
		if u.DownloadOnly {
			nameSplit := strings.Split(u.Resource, "/")
			filename := nameSplit[len(nameSplit) - 1]
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
			pg.WrapContent(c.Width - 1)
			c.PageState.Add(pg)
			c.SetPercentRead()
			c.ClearMessage()
			c.SetHeaderUrl()
			c.Draw()
		}
	case "gemini":
		capsule, err := gemini.Visit(u.Host, u.Port, u.Resource, &c.Certs)
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			return
		}
		go saveConfig()
		switch capsule.Status {
		case 1:
			c.search("", u.Full, capsule.Content)
		case 2:
			if capsule.MimeMaj == "text" {
				pg := MakePage(u, capsule.Content, capsule.Links)
				pg.WrapContent(c.Width - 1)
				c.PageState.Add(pg)
				c.SetPercentRead()
				c.ClearMessage()
				c.SetHeaderUrl()
				c.Draw()
			} else {
				c.SetMessage("The file is non-text: (o)pen or (w)rite to disk", false)
				c.DrawMessage()
				var ch rune
				for {
					ch = cui.Getch()
					if ch == 'o' || ch == 'w' {
						break
					}
				}
				switch ch {
				case 'o':
					mime := fmt.Sprintf("%s/%s", capsule.MimeMaj, capsule.MimeMin)
					var term bool
					if c.Options["terminalonly"] == "true" {
						term = true
					} else {
						term = false
					}
					mcEntry, err := mc.FindMatch(mime, "view", term)
					if err != nil {
						c.SetMessage(err.Error(), true)
						c.DrawMessage()
						return
					}
					file, err := ioutil.TempFile("/tmp/", "bombadillo-*.tmp")
					if err != nil {
						c.SetMessage("Unable to create temporary file for opening, aborting file open", true)
						c.DrawMessage()
						return
					}
					// defer os.Remove(file.Name())
					file.Write([]byte(capsule.Content))
					com, e := mcEntry.Command(file.Name())
					if e != nil {
						c.SetMessage(e.Error(), true)
						c.DrawMessage()
						return
					}
					com.Stdin = os.Stdin
					com.Stdout = os.Stdout
					com.Stderr = os.Stderr
					if c.Options["terminalonly"] == "true" {
						cui.Clear("screen")
					}
					com.Run()
					c.SetMessage("File opened by an appropriate program", true)
					c.DrawMessage()
					c.Draw()
				case 'w':
					nameSplit := strings.Split(u.Resource, "/")
					filename := nameSplit[len(nameSplit) - 1]
					c.saveFileFromData(capsule.Content, filename)
				}
			}
		case 3:
			c.SetMessage("[3] Redirect. Follow redirect? y or any other key for no", false)
			c.DrawMessage()
			ch := cui.Getch()
			if ch == 'y' || ch == 'Y' {
				c.Visit(capsule.Content)
			} else {
				c.SetMessage("Redirect aborted", false)
				c.DrawMessage()
			}
		}
	case "telnet":
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
	case "http", "https":
		c.SetMessage("Attempting to open in web browser", false)
		c.DrawMessage()
		if strings.ToUpper(c.Options["openhttp"]) == "TRUE"  {
			msg, err := http.OpenInBrowser(u.Full)
			if err != nil {
				c.SetMessage(err.Error(), true)
			} else {
				c.SetMessage(msg, false)
			}
			c.DrawMessage()
		} else {
			c.SetMessage("'openhttp' is not set to true, cannot open web link", false)
			c.DrawMessage()
		}
	default:
		c.SetMessage(fmt.Sprintf("%q is not a supported protocol", u.Scheme), true)
		c.DrawMessage()
	}
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func MakeClient(name string) *client {
	c := client{0, 0, defaultOptions, "", false, MakePages(), MakeBookmarks(), MakeHeadbar(name), MakeFootbar(), gemini.MakeTofuDigest()}
	return &c
}

