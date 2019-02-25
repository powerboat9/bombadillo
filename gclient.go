package main

import (
	"fmt"
	"gsock/gopher"
	"os"
	"bufio"
	"regexp"
	"strings"
	"strconv"
)

var history gopher.History = gopher.MakeHistory()

func err_exit(err string, code int) {
	fmt.Println(err)
	os.Exit(code)
}

func getln() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(": ")
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}


func route_input(s string) {
	sl := strings.ToLower(s)
	if sl == "quit" || sl == "exit" || sl == "q" {
		err_exit("Quitting...", 0)
	} else if num, _ := regexp.MatchString(`^\d+$`, s); num && history.Length > 0 {
		linkcount := len(history.Collection[history.Position].Links)
		item, _ := strconv.Atoi(s)
		if item <= linkcount {
			linkurl := history.Collection[history.Position].Links[item - 1]

			v, err := history.Visit(linkurl)
			if err != nil {
				fmt.Println(err.Error())
			}
			if v.Address.IsBinary {
				// Query for download here
				fmt.Println("Would you like to download this file?")
			} else {
				history.Add(v)
				history.DisplayCurrentView()
			}
		} else {
			fmt.Println("Invalid link id")
		}
	} else if sl == "back" || sl == "b" {
		history.GoBack()
	} else if sl == "forward" || sl == "f" {
		history.GoForward()
	} else {
		v, err := history.Visit(s)
		if err != nil {
			fmt.Println(err.Error())
		}
		if v.Address.IsBinary {
			// Query for download here
			fmt.Println("Would you like to download this file?")
		} else {
			history.Add(v)
			history.DisplayCurrentView()
		}
	}
}

func make_request(s string) ([]string, gopher.Url, error) {
	u, _ := gopher.MakeUrl(s)
	text, err := gopher.Retrieve(u)
	if err != nil {
		return []string{}, u, err
	}
	return strings.Split(string(text), "\n"), u, nil
}

func main() {
	history.Position = -1
	var inp string
	if len(os.Args) >= 2 {
		inp = os.Args[1]
		route_input(inp)
	} 

	for {
		inp = getln()
		if inp == "" {
			continue
		}
		route_input(inp)
	}
}
