package main

import (
	"fmt"
	"regexp"
	"strings"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Url struct {
	Scheme string
	Host string
	Port string
	Resource string
	Full string
	Mime string
	DownloadOnly bool
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

// There are currently no receivers for the Url struct


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\


// MakeUrl is a Url constructor that takes in a string
// representation of a url and returns a Url struct and
// an error (or nil).
func MakeUrl(u string) (Url, error) {
	var out Url
	re := regexp.MustCompile(`^((?P<scheme>[a-zA-Z]+):\/\/)?(?P<host>[\w\-\.\d]+)(?::(?P<port>\d+)?)?(?:/(?P<type>[01345679gIhisp])?)?(?P<resource>.*)?$`)
	match := re.FindStringSubmatch(u)

	if valid := re.MatchString(u); !valid {
		return out, fmt.Errorf("Invalid url, unable to parse")
	}

	for i, name := range re.SubexpNames() {
		switch name {
		case "scheme":
			out.Scheme = match[i]
		case "host":
			out.Host = match[i]
		case "port":
			out.Port = match[i]
		case "type":
			out.Mime = match[i]
		case "resource":
			out.Resource = match[i]
		}
	}

	if out.Host == "" {
		return out, fmt.Errorf("no host")
	}

	out.Scheme = strings.ToLower(out.Scheme)

	if out.Scheme == "" {
		out.Scheme = "gopher"
	}

	if out.Scheme == "gopher" && out.Port == "" {
		out.Port = "70"
	} else if out.Scheme == "http" && out.Port == "" {
		out.Port = "80"
	} else if out.Scheme == "https" && out.Port == "" {
		out.Port = "443"
	} else if out.Scheme == "gemini" && out.Port == "" {
		out.Port = "1965"
	} else if out.Scheme == "telnet" && out.Port == "" {
		out.Port = "23"
	}

	if out.Scheme == "gopher" {
		if out.Mime == "" {
			out.Mime = "1"
		}
		if out.Resource == "" || out.Resource == "/" {
			out.Mime = "1"
		}
		if out.Mime == "7" && strings.Contains(out.Resource, "\t") {
			out.Mime = "1"
		}
		switch out.Mime {
		case "1", "0", "h", "7":
			out.DownloadOnly = false
		default:
			out.DownloadOnly = true
		}
	} else {
		out.Resource = fmt.Sprintf("%s%s", out.Mime, out.Resource)
		out.Mime = ""
	}

	out.Full = out.Scheme + "://" + out.Host + ":" + out.Port + "/" + out.Mime + out.Resource

	return out, nil
}
