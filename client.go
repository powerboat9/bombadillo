package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tildegit.org/sloum/bombadillo/cmdparse"
	"tildegit.org/sloum/bombadillo/cui"
	// "tildegit.org/sloum/bombadillo/gemini"
	// "tildegit.org/sloum/bombadillo/gopher"
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
	PageState Pages
	BookMarks Bookmarks
	TopBar Headbar
	FootBar Footbar
}


//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (c *client) GetSize() {
	for {
		redraw := false
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
			redraw = true
		}

		c.Height = h
		c.Width = w

		if redraw {
			c.Draw()
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *client) Draw() {
	// TODO build this out.
	// It should call all of the renders
	// and add them to the a string buffer
	// It should then print the buffer
}

func (c *client) TakeControlInput() {
	input := cui.Getch()

	switch input {
	case 'j', 'J':
		// scroll down one line
		c.Scroll(1)
	case 'k', 'K':
		// scroll up one line
		c.Scroll(-1)
	case 'q', 'Q':
		// quite bombadillo
		cui.Exit()
	case 'g':
		// scroll to top
		c.Scroll(-len(c.PageState.History[c.PageState.Position].WrappedContent))
	case 'G':
		// scroll to bottom
		c.Scroll(len(c.PageState.History[c.PageState.Position].WrappedContent))
	case 'd':
		// scroll down 75%
		distance := c.Height - c.Height / 4
		c.Scroll(distance)
	case 'u':
		// scroll up 75%
		distance := c.Height - c.Height / 4
		c.Scroll(-distance)
	case 'b':
		// go back
		err := c.PageState.NavigateHistory(-1)
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
	case 'f', 'F':
		// go forward
		err := c.PageState.NavigateHistory(1)
		if err != nil {
			c.SetMessage(err.Error(), false)
			c.DrawMessage()
		} else {
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
		entry, err := cui.GetLine()
		c.ClearMessageLine()
		if err != nil {
			c.SetMessage(err.Error(), true)
			c.DrawMessage()
			break
		} else if strings.TrimSpace(entry) == "" {
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
		// err = doLinkCommand(com.Action, com.Target)
	case cmdparse.DOAS:
		c.doCommandAs(com.Action, com.Value)
	case cmdparse.DOLINKAS:
		// err = doLinkCommandAs(com.Action, com.Target, com.Value)
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
	case "SEARCH":
		c.search()
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

	case "WRITE", "W":
		// TODO figure out how best to handle file
		// writing... it will depend on request model
		// using fetch would be best
		// - - - - - - - - - - - - - - - - - - - - -
		// var data []byte
		// if values[0] == "." {
			// d, err := c.getCurrentPageRawData()
			// if err != nil {
				// c.SetMessage(err.Error(), true)
				// c.DrawMessage()
				// return
			// }
			// data = []byte(d)
		// }
		// fp, err := c.saveFile(data, strings.Join(values[1:], " "))
		// if err != nil {
			// c.SetMessage(err.Error(), true)
			// c.DrawMessage()
			// return
		// }
		// c.SetMessage(fmt.Sprintf("File saved to: %s", fp), false)
		// c.DrawMessage()

	case "SET", "S":
		if _, ok := c.Options[values[0]]; ok {
			c.Options[values[0]] = strings.Join(values[1:], " ")
			err := saveConfig()
			if err != nil {
				c.SetMessage("Value set, but error saving config to file", true)
				c.DrawMessage()
			} else {
				c.SetMessage(fmt.Sprintf("%s is now set to %q", values[0], c.Options[values[0]]), true)
				c.DrawMessage()
			}
			return
		}
		c.SetMessage(fmt.Sprintf("Unable to set %s, it does not exist", values[0]), true)
		c.DrawMessage()
		return
	}
	c.SetMessage(fmt.Sprintf("Unknown command structure"), true)
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

func (c *client) saveFile(data []byte, name string) (string, error) {
	savePath := c.Options["savelocation"] + name
	err := ioutil.WriteFile(savePath, data, 0644)
	if err != nil {
		return "", err
	}

	return savePath, nil
}

func (c *client) search() {
	c.ClearMessage()
	c.ClearMessageLine()
	fmt.Print("?")
	entry, err := cui.GetLine()
	c.ClearMessageLine()
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	} else if strings.TrimSpace(entry) == "" {
		return
	}
	u, err := MakeUrl(c.Options["searchurl"])
	if err != nil {
		c.SetMessage("'searchurl' is not set to a valid url", true)
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
	page := c.PageState.History[c.PageState.Position]
	bottom := len(page.WrappedContent) - c.Height
	if amount < 0 && page.ScrollPosition == 0 {
		c.SetMessage("You are already at the top", false)
		c.DrawMessage()
		fmt.Print("\a")
		return
	} else if amount > 0 && page.ScrollPosition == bottom || bottom < 0 {
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

	page.ScrollPosition = newScrollPosition
	c.Draw()
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
	leadIn, leadOut := "", ""
	if isError {
		leadIn = "\033[31m"
		leadOut = "\033[0m"
	}

	c.Message = fmt.Sprintf("%s%s%s", leadIn, msg, leadOut)
}

func (c *client) DrawMessage() {
	c.ClearMessageLine()
	cui.MoveCursorTo(c.Height-1, 0)
	fmt.Print(c.Message)
}

func (c *client) ClearMessage() {
	c.Message = ""
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

	c.SetMessage(fmt.Sprintf("Invalid link id: %s", l), true)
	c.DrawMessage()
}

func (c *client) Visit(url string) {
	// TODO both gemini and gopher should return a string
	// The wrap lines function in cui needs to be rewritten
	u, err := MakeUrl(url)
	if err != nil {
		c.SetMessage(err.Error(), true)
		c.DrawMessage()
		return
	}

	switch u.Scheme {
	case "gopher":
		// TODO send over to gopher request
	case "gemini":
		// TODO send over to gemini request
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
	var userinfo, _ = user.Current()
	var options = map[string]string{
		"homeurl":      "gopher://colorfield.space:70/1/bombadillo-info",
		"savelocation": userinfo.HomeDir,
		"searchengine": "gopher://gopher.floodgap.com:70/7/v2/vs",
		"openhttp":     "false",
		"httpbrowser":  "lynx",
		"configlocation": userinfo.HomeDir,
	}
	c := client{0, 0, options, "", MakePages(), MakeBookmarks(), MakeHeadbar(name), MakeFootbar()}
	c.GetSize()
	return &c
}

// Retrieve a byte slice of raw response dataa 
// from a url string
func Fetch(url string) ([]byte, error) {
	u, err := MakeUrl(url)
	if err != nil {
		return []byte(""), err
	}

	timeOut := time.Duration(5) * time.Second

	if u.Host == "" || u.Port == "" {
		return []byte(""), fmt.Errorf("Incomplete request url")
	}

	addr := u.Host + ":" + u.Port

	conn, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		return []byte(""), err
	}

	send := u.Resource + "\n"

	_, err = conn.Write([]byte(send))
	if err != nil {
		return []byte(""), err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return []byte(""), err
	}

	return result, err
}
