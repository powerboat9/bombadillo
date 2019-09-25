package gemini

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)


type Capsule struct {
	MimeMaj	string
	MimeMin string
	Status	int
	Content	string
	Links []string
}

type TofuDigest struct {
	db  map[string][]map[string]string
}

//------------------------------------------------\\
// + + +          R E C E I V E R S          + + + \\
//--------------------------------------------------\\

func (t *TofuDigest) Remove(host string, indexToRemove int) error {
	if _, ok := t.db[host]; ok {
		if indexToRemove < 0 || indexToRemove >= len(t.db[host]) {
			return fmt.Errorf("Invalid index")
		} else if len(t.db[host]) > indexToRemove {
			t.db[host] = append(t.db[host][:indexToRemove], t.db[host][indexToRemove+1:]...)
		} else if len(t.db[host]) - 1 == indexToRemove {
			t.db[host] = t.db[host][:indexToRemove]
		}
		return nil
	}
	return fmt.Errorf("Invalid host")
}

func (t *TofuDigest) Add(host, hash string, start, end int64) {
	s := strconv.FormatInt(start, 10)
	e := strconv.FormatInt(end, 10)
	added := strconv.FormatInt(time.Now().Unix(), 10)
	entry := map[string]string{"hash": hash, "start": s, "end": e, "added": added}
	t.db[host] = append(t.db[host], entry)
}

// Removes all entries that are expired
func (t *TofuDigest) Clean() {
	now := time.Now()
	for host, slice := range t.db {
		for index, entry := range slice {
			intFromStringTime, err := strconv.ParseInt(entry["end"], 10, 64)
			if err != nil || now.After(time.Unix(intFromStringTime, 0)) {
				t.Remove(host, index)
			}
		}
	}
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

	// Verify that the handshake ahs completed and that
	// the hostname on the certificate(s) from the server
	// is the hostname we have requested
	connState := conn.ConnectionState()
	if connState.HandshakeComplete {
		if len(connState.PeerCertificates) > 0 {
			for _, cert := range connState.PeerCertificates {
				if err = cert.VerifyHostname(host); err == nil {
					break
				} 
			}
			if err != nil {
				return "", err
			}
		}
	}


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

func Fetch(host, port, resource string) ([]byte, error) {
	rawResp, err := Retrieve(host, port, resource)
	if err != nil {
		return make([]byte, 0), err
	} 

	resp := strings.SplitN(rawResp, "\r\n", 2)
	if len(resp) != 2 {
		if err != nil {
			return make([]byte, 0), fmt.Errorf("Invalid response from server")
		} 
	}
	header := strings.SplitN(resp[0], " ", 2)
	if len([]rune(header[0])) != 2 {
		header = strings.SplitN(resp[0], "\t", 2)
		if len([]rune(header[0])) != 2 {
			return make([]byte,0), fmt.Errorf("Invalid response format from server")
		} 
	} 

	// Get status code single digit form
	status, err := strconv.Atoi(string(header[0][0]))
	if err != nil {
		return make([]byte, 0), fmt.Errorf("Invalid status response from server")
	}

	if status != 2 {
		switch status {
		case 1:
			return make([]byte, 0), fmt.Errorf("[1] Queries cannot be saved.")
		case 3:
			return make([]byte, 0), fmt.Errorf("[3] Redirects cannot be saved.")
		case 4:
			return make([]byte, 0), fmt.Errorf("[4] Temporary Failure.")
		case 5:
			return make([]byte, 0), fmt.Errorf("[5] Permanent Failure.")
		case 6:
			return make([]byte, 0), fmt.Errorf("[6] Client Certificate Required (Not supported by Bombadillo)")
		default:
			return make([]byte, 0), fmt.Errorf("Invalid response status from server")
		}
	}

	return []byte(resp[1]), nil

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
		capsule.Content = header[1]
		return capsule, nil
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
			if len(resource) > 0 && resource[0] != '/' {
				resource = fmt.Sprintf("/%s", resource)
			} else if resource == "" {
				resource = "/"
			}
			currentUrl := fmt.Sprintf("gemini://%s:%s%s", host, port, resource)
			rootUrl := fmt.Sprintf("gemini://%s:%s", host, port)
			capsule.Content, capsule.Links = parseGemini(body, rootUrl, currentUrl)
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

func parseGemini(b, rootUrl, currentUrl string) (string, []string) {
	splitContent := strings.Split(b, "\n")
	links := make([]string, 0, 10)

	for i, ln := range splitContent {
		splitContent[i] = strings.Trim(ln, "\r\n")
		if len([]rune(ln)) > 3 && ln[:2] == "=>" {
			var link, decorator string
			subLn := strings.Trim(ln[2:], "\r\n\t \a")
			splitPoint := strings.IndexAny(subLn, " \t")

			if splitPoint < 0 || len([]rune(subLn)) - 1 <= splitPoint {
				link = subLn
				decorator = subLn
			} else {
				link = strings.Trim(subLn[:splitPoint], "\t\n\r \a")
				decorator = strings.Trim(subLn[splitPoint:], "\t\n\r \a")
			}

			if strings.Index(link, "://") < 0  {
				link = handleRelativeUrl(link, rootUrl, currentUrl)
			}

			links = append(links, link)
			linknum := fmt.Sprintf("[%d]", len(links))
			splitContent[i] = fmt.Sprintf("%-5s %s", linknum, decorator)
		}
	}
	return strings.Join(splitContent, "\n"), links
}

func handleRelativeUrl(u, root, current string) string {
	if len(u) < 1 {
		return u
	}

	if u[0] == '/' {
		return fmt.Sprintf("%s%s", root, u)
	}

	ind := strings.LastIndex(current, "/")
	if ind < 10 {
		return fmt.Sprintf("%s/%s", root, u)
	}

	current = current[:ind + 1]
	return fmt.Sprintf("%s%s", current, u)
}


func MakeCapsule() Capsule {
	return Capsule{"", "", 0, "", make([]string, 0, 5)}
}

