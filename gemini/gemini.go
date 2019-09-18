package gemini

import (
	"crypto/tls"
	"fmt"
	"net"
	"io/ioutil"
	// "strings"
	"time"

	// "tildegit.org/sloum/mailcap"
)


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func Retrieve(host, port, resource string) ([]byte, error) {
	nullRes := make([]byte, 0)
	timeOut := time.Duration(5) * time.Second

	if host == "" || port == "" {
		return nullRes, fmt.Errorf("Incomplete request url")
	}

	addr := host + ":" + port

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		return nullRes, err
	}

	secureConn := tls.Client(conn, conf)

	send := resource + "\n"

	_, err = secureConn.Write([]byte(send))
	if err != nil {
		return nullRes, err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nullRes, err
	}

	return result, nil
}

func Visit(host, port, resource string) (string, []string, error) {
	resp, err := Retrieve(host, port, resource)
	if err != nil {
		return "", []string{}, err
	} 
	
	// TODO break out the header
	// header := ""
	mime := ""
	mimeMaj := mime
	mimeMin := mime
	// status := ""
	content := string(resp)

	if mimeMaj == "text" &&  mimeMin == "gemini" {
		// text := string(resp)
		// links := []string{}

		// TODO parse geminimap from 'content'
	} else if mimeMaj == "text" {
		// TODO just return the text
	} else {
		// TODO use mailcap to try and open the file
	}


	return content, []string{}, nil
}

