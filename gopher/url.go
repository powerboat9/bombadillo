package gopher

import (
	"errors"
	"regexp"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

// The url struct represents a URL for the rest of the system.
// It includes component parts as well as a full URL string.
type Url struct {
	Scheme     string
	Host       string
	Port       string
	Gophertype string
	Resource   string
	Full       string
	IsBinary   bool
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

// MakeUrl is a Url constructor that takes in a string
// representation of a url and returns a Url struct and
// an error (or nil).
func MakeUrl(u string) (Url, error) {
	var out Url
	re := regexp.MustCompile(`^((?P<scheme>gopher|http|https|ftp|telnet):\/\/)?(?P<host>[\w\-\.\d]+)(?::(?P<port>\d+)?)?(?:/(?P<type>[01345679gIhisp])?)?(?P<resource>.*)?$`)
	match := re.FindStringSubmatch(u)

	if valid := re.MatchString(u); !valid {
		return out, errors.New("Invalid URL or command character")
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
			out.Gophertype = match[i]
		case "resource":
			out.Resource = match[i]
		}
	}

	if out.Scheme == "" {
		out.Scheme = "gopher"
	}

	if out.Host == "" {
		return out, errors.New("no host")
	}

	if out.Scheme == "gopher" && out.Port == "" {
		out.Port = "70"
	} else if out.Scheme == "http" && out.Port == "" {
		out.Port = "80"
	} else if out.Scheme == "https" && out.Port == "" {
		out.Port = "443"
	}

	if out.Gophertype == "" && (out.Resource == "" || out.Resource == "/") {
		out.Gophertype = "1"
	}

	if out.Scheme == "gopher" && out.Gophertype == "" {
		out.Gophertype = "0"
	}

	switch out.Gophertype {
	case "1", "0", "h", "7":
		out.IsBinary = false
	default:
		out.IsBinary = true
	}


	out.Full = out.Scheme + "://" + out.Host + ":" + out.Port + "/" + out.Gophertype + out.Resource

	return out, nil
}
