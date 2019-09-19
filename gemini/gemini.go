package gemini

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	// "tildegit.org/sloum/mailcap"
)

type Capsule struct {
	MimeMaj	string
	MimeMin string
	Status	int
	Content	string
	Links []string
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func Retrieve(host, port, resource string) (string, error) {
	if host == "" || port == "" {
		return "", fmt.Errorf("Incomplete request url")
	}

	addr := host + ":" + port

	conf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	send := "gemini://" + addr + "/" + resource + "\r\n"

	_, err = conn.Write([]byte(send))
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func Visit(host, port, resource string) (Capsule, error) {
	capsule := MakeCapsule()
	rawResp, err := Retrieve(host, port, resource)
	if err != nil {
		return capsule, err
	} 

	resp := strings.SplitN(rawResp, "\r\n", 2)
	if len(resp) != 2 {
		if err != nil {
			return capsule, fmt.Errorf("Invalid response from server")
		} 
	}
	header := strings.SplitN(resp[0], " ", 2)
	if len([]rune(header[0])) != 2 {
		header = strings.SplitN(resp[0], "\t", 2)
		if len([]rune(header[0])) != 2 {
			return capsule, fmt.Errorf("Invalid response format from server")
		} 
	} 

	body := resp[1]
	
	// Get status code single digit form
	capsule.Status, err = strconv.Atoi(string(header[0][0]))
	if err != nil {
		return capsule, fmt.Errorf("Invalid status response from server")
	}

	// Parse the meta as needed
	var meta string

	switch capsule.Status {
	case 1:
		// handle search
		return capsule, fmt.Errorf("Gemini input not yet supported")
	case 2:
		mimeAndCharset := strings.Split(header[1], ";")
		meta = mimeAndCharset[0]
		minMajMime := strings.Split(meta, "/")
		if len(minMajMime) < 2 {
			return capsule, fmt.Errorf("Improperly formatted mimetype received from server")
		}
		capsule.MimeMaj = minMajMime[0]
		capsule.MimeMin = minMajMime[1]
		if capsule.MimeMaj == "text" && capsule.MimeMin == "gemini" {
			rootUrl := fmt.Sprintf("gemini://%s:%s", host, port)
			capsule.Content, capsule.Links = parseGemini(body, rootUrl)
		} else {
			capsule.Content = body
		}
		return capsule, nil
	case 3:
		// The client will handle informing the user of a redirect
		// and then request the new url
		capsule.Content = header[1]
		return capsule, nil
	case 4:
		return capsule, fmt.Errorf("[4] Temporary Failure. %s", header[1])
	case 5:
		return capsule, fmt.Errorf("[5] Permanent Failure. %s", header[1])
	case 6:
		return capsule, fmt.Errorf("[6] Client Certificate Required (Not supported by Bombadillo)")
	default:
		return capsule, fmt.Errorf("Invalid response status from server")
	}
}

func parseGemini(b, rootUrl string) (string, []string) {
	splitContent := strings.Split(b, "\n")
	links := make([]string, 0, 10)

	for i, ln := range splitContent {
		splitContent[i] = strings.Trim(ln, "\r\n")
		if len([]rune(ln)) > 3 && ln[:2] == "=>" {
			trimmedSubLn := strings.Trim(ln[2:], "\r\n\t \a")
			lineSplit := strings.SplitN(trimmedSubLn, " ", 2)
			if len(lineSplit) != 2 {
				lineSplit = append(lineSplit, lineSplit[0])
			}
			lineSplit[0] = strings.Trim(lineSplit[0], "\t\n\r \a")
			lineSplit[1] = strings.Trim(lineSplit[1], "\t\n\r \a")
			if len(lineSplit[0]) > 0 && lineSplit[0][0] == '/' {
				lineSplit[0] = fmt.Sprintf("%s%s", rootUrl, lineSplit[0])
			}
			links = append(links, lineSplit[0])
			linknum := fmt.Sprintf("[%d]", len(links))
			splitContent[i] = fmt.Sprintf("%-5s %s", linknum, lineSplit[1])
		}
	}
	return strings.Join(splitContent, "\n"), links
}


func MakeCapsule() Capsule {
	return Capsule{"", "", 0, "", make([]string, 0, 5)}
}

