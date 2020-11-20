package gemini

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Capsule struct {
	MimeMaj string
	MimeMin string
	Status  int
	Content string
	Links   []string
}

type TofuDigest struct {
	certs      map[string]string
}

var BlockBehavior string = "block"
var TlsTimeout time.Duration = time.Duration(15) * time.Second

//------------------------------------------------\\
// + + +          R E C E I V E R S          + + + \\
//--------------------------------------------------\\

func (t *TofuDigest) Purge(host string) error {
	host = strings.ToLower(host)
	if host == "*" {
		t.certs = make(map[string]string)
		return nil
	} else if _, ok := t.certs[strings.ToLower(host)]; ok {
		delete(t.certs, host)
		return nil
	}
	return fmt.Errorf("Invalid host %q", host)
}

func (t *TofuDigest) Add(host, hash string, time int64) {
	t.certs[strings.ToLower(host)] = fmt.Sprintf("%s|%d", hash, time)
}

func (t *TofuDigest) Exists(host string) bool {
	if _, ok := t.certs[strings.ToLower(host)]; ok {
		return true
	}
	return false
}

func (t *TofuDigest) Find(host string) (string, error) {
	if hash, ok := t.certs[strings.ToLower(host)]; ok {
		return hash, nil
	}
	return "", fmt.Errorf("Invalid hostname, no key saved")
}

func (t *TofuDigest) Match(host, localCert string, cState *tls.ConnectionState) error {
	now := time.Now()

	for _, cert := range cState.PeerCertificates {
		if localCert != hashCert(cert.Raw) {
			continue
		}

		if now.Before(cert.NotBefore) {
			return fmt.Errorf("Certificate is not valid yet")
		}

		if now.After(cert.NotAfter) {
			return fmt.Errorf("EXP")
		}

		if err := cert.VerifyHostname(host); err != nil && cert.Subject.CommonName != host {
			return fmt.Errorf("Certificate error: %s", err)
		}

		return nil
	}

	return fmt.Errorf("No matching certificate was found for host %q", host)
}

func (t *TofuDigest) newCert(host string, cState *tls.ConnectionState) error {
	host = strings.ToLower(host)
	now := time.Now()
	var reasons strings.Builder

	for index, cert := range cState.PeerCertificates {
		if index > 0 {
			reasons.WriteString("; ")
		}
		if now.Before(cert.NotBefore) {
			reasons.WriteString(fmt.Sprintf("Cert [%d] is not valid yet", index+1))
			continue
		}

		if now.After(cert.NotAfter) {
			reasons.WriteString(fmt.Sprintf("Cert [%d] is expired", index+1))
			continue
		}

		if err := cert.VerifyHostname(host); err != nil && cert.Subject.CommonName != host {
			reasons.WriteString(fmt.Sprintf("Cert [%d] hostname does not match", index+1))
			continue
		}

		t.Add(host, hashCert(cert.Raw), cert.NotAfter.Unix())
		return nil
	}

	return fmt.Errorf(reasons.String())
}

func (t *TofuDigest) GetCertAndTimestamp(host string) (string, int64, error) {
	certTs, err := t.Find(host)
	if err != nil {
		return "", -1, err
	}
	certTsSplit := strings.SplitN(certTs, "|", -1)
	if len(certTsSplit) < 2 {
		_ = t.Purge(host)
		return certTsSplit[0], -1, fmt.Errorf("Invalid certstring, no delimiter")
	}
	ts, err := strconv.ParseInt(certTsSplit[1], 10, 64)
	if err != nil {
		_ = t.Purge(host)
		return certTsSplit[0], -1, err
	}
	now := time.Now()
	if ts < now.Unix() {
		// Ignore error return here since an error would indicate
		// the host does not exist and we have already checked for
		// that and the desired outcome of the action is that the
		// host will no longer exist, so we are good either way
		_ = t.Purge(host)
		return "", -1, fmt.Errorf("Expired cert")
	}
	return certTsSplit[0], ts, nil
}

func (t *TofuDigest) IniDump() string {
	if len(t.certs) < 1 {
		return ""
	}
	var out strings.Builder
	out.WriteString("[CERTS]\n")
	for k, v := range t.certs {
		out.WriteString(k)
		out.WriteString("=")
		out.WriteString(v)
		out.WriteString("\n")
	}
	return out.String()
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func Retrieve(host, port, resource string, td *TofuDigest) (string, error) {
	if host == "" || port == "" {
		return "", fmt.Errorf("Incomplete request url")
	}

	addr := host + ":" + port

	conf := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: TlsTimeout}, "tcp", addr, conf)
	if err != nil {
		return "", fmt.Errorf("TLS Dial Error: %s", err.Error())
	}

	defer conn.Close()

	connState := conn.ConnectionState()

	// Begin TOFU screening...

	// If no certificates are offered, bail out
	if len(connState.PeerCertificates) < 1 {
		return "", fmt.Errorf("Insecure, no certificates offered by server")
	}

	localCert, localTs, err := td.GetCertAndTimestamp(host)

	if localTs > 0 {
		// See if we have a matching cert
		err := td.Match(host, localCert, &connState)
		if err != nil && err.Error() != "EXP" {
			// If there is no match and it isnt because of an expiration
			// just return the error
			return "", err
		} else if err != nil {
			// The cert expired, see if they are offering one that is valid...
			err := td.newCert(host, &connState)
			if err != nil {
				// If there are no valid certs to offer, let the client know
				return "", err
			}
		}
	} else {
		err = td.newCert(host, &connState)
		if err != nil {
			// If there are no valid certs to offer, let the client know
			return "", err
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

func Fetch(host, port, resource string, td *TofuDigest) ([]byte, error) {
	rawResp, err := Retrieve(host, port, resource, td)
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
			return make([]byte, 0), fmt.Errorf("Invalid response format from server")
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
			return make([]byte, 0), fmt.Errorf("[6] Client Certificate Required (Unsupported)")
		default:
			return make([]byte, 0), fmt.Errorf("Invalid response status from server")
		}
	}

	return []byte(resp[1]), nil

}

func Visit(host, port, resource string, td *TofuDigest) (Capsule, error) {
	capsule := MakeCapsule()
	rawResp, err := Retrieve(host, port, resource, td)
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
		if meta == "" {
			meta = "text/gemini"
		}
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
			capsule.Content, capsule.Links = parseGemini(body, currentUrl)
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
		return capsule, fmt.Errorf("[6] Client Certificate Required (Unsupported)")
	default:
		return capsule, fmt.Errorf("Invalid response status from server")
	}
}

func parseGemini(b, currentUrl string) (string, []string) {
	splitContent := strings.Split(b, "\n")
	links := make([]string, 0, 10)

	inPreBlock := false
	spacer := "      "

	outputIndex := 0
	for i, ln := range splitContent {
		splitContent[i] = strings.Trim(ln, "\r\n")
		isPreBlockDeclaration := strings.HasPrefix(ln, "```")
		if isPreBlockDeclaration && !inPreBlock && (BlockBehavior == "both" || BlockBehavior == "alt") {
			inPreBlock = !inPreBlock
			alt := strings.TrimSpace(ln)
			if len(alt) > 3 {
				alt = strings.TrimSpace(alt[3:])
				splitContent[outputIndex] = fmt.Sprintf("%s[ALT][ %s ]", spacer, alt)
				outputIndex++
			}
		} else if isPreBlockDeclaration {
			inPreBlock = !inPreBlock
		} else if len([]rune(ln)) > 3 && ln[:2] == "=>" && !inPreBlock {
			var link, decorator string
			subLn := strings.Trim(ln[2:], "\r\n\t \a")
			splitPoint := strings.IndexAny(subLn, " \t")

			if splitPoint < 0 || len([]rune(subLn))-1 <= splitPoint {
				link = subLn
				decorator = subLn
			} else {
				link = strings.Trim(subLn[:splitPoint], "\t\n\r \a")
				decorator = strings.Trim(subLn[splitPoint:], "\t\n\r \a")
			}

			if strings.Index(link, "://") < 0 {
				link, _ = HandleRelativeUrl(link, currentUrl)
			}

			links = append(links, link)
			linknum := fmt.Sprintf("[%d]", len(links))
			splitContent[outputIndex] = fmt.Sprintf("%-5s %s", linknum, decorator)
			outputIndex++
		} else {
			if inPreBlock && (BlockBehavior == "alt" || BlockBehavior == "neither") {
				continue
			}
			var leader, tail string = "", ""
			if len(ln) > 0 && ln[0] == '#' {
				leader = "\033[1m"
				tail = "\033[0m"
			}
			splitContent[outputIndex] = fmt.Sprintf("%s%s%s%s", spacer, leader, ln, tail)
			outputIndex++
		}
	}
	return strings.Join(splitContent[:outputIndex], "\n"), links
}

// handleRelativeUrl provides link completion
func HandleRelativeUrl(relLink, current string) (string, error) {
	base, err := url.Parse(current)
	if err != nil {
		return relLink, err
	}
	rel, err := url.Parse(relLink)
	if err != nil {
		return relLink, err
	}
	return base.ResolveReference(rel).String(), nil
}

func hashCert(cert []byte) string {
	hash := sha1.Sum(cert)
	hex := make([][]byte, len(hash))
	for i, data := range hash {
		hex[i] = []byte(fmt.Sprintf("%02X", data))
	}
	return fmt.Sprintf("%s", string(bytes.Join(hex, []byte(":"))))
}

func MakeCapsule() Capsule {
	return Capsule{"", "", 0, "", make([]string, 0, 5)}
}

func MakeTofuDigest() TofuDigest {
	return TofuDigest{make(map[string]string)}
}
