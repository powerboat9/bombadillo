package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	if len(u) < 1 {
		return Url{}, fmt.Errorf("Invalid url, unable to parse")
	}
	if strings.HasPrefix(u, "finger://") {
		return parseFinger(u)
	}

	var out Url
	if local := strings.HasPrefix(u, "local://"); u[0] == '/' || u[0] == '.' || u[0] == '~' || local {
		if local && len(u) > 8 {
			u = u[8:]
		}
		home, err := os.UserHomeDir()
		if err != nil {
			home = ""
		}
		u = strings.Replace(u, "~", home, 1)
		res, err := filepath.Abs(u)
		if err != nil {
			return out, fmt.Errorf("Invalid path, unable to parse")
		}
		out.Scheme = "local"
		out.Host = ""
		out.Port = ""
		out.Mime = ""
		out.Resource = res
		out.Full = out.Scheme + "://" + out.Resource
		return out, nil
	}

	re := regexp.MustCompile(`^((?P<scheme>[a-zA-Z]+):\/\/)?(?P<host>[\w\-\.\d/]+)(?::(?P<port>\d+)?)?(?:/(?P<type>[01345679gIhisp])?)?(?P<resource>.*)?$`)
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

func parseFinger(u string) (Url, error) {
	var out Url
	out.Scheme = "finger"
	if len(u) < 10 {
		return out, fmt.Errorf("Invalid finger address")
	}
	u = u[9:]
	userPlusAddress := strings.Split(u, "@")
	if len(userPlusAddress) > 1 {
		out.Resource = userPlusAddress[0]
		u = userPlusAddress[1]
	} 
	hostPort := strings.Split(u, ":")
	if len(hostPort) < 2 {
		out.Port = "79"
	} else {
		out.Port = hostPort[1]
	}
	out.Host = hostPort[0]
	resource := ""
	if out.Resource != "" {
		resource = out.Resource + "@"
	}
	out.Full = fmt.Sprintf("%s://%s%s:%s", out.Scheme, resource, out.Host, out.Port)
	return out, nil
}


