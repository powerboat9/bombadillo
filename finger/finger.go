package finger

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func Finger(host, port, resource string) (string, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	timeOut := time.Duration(3) * time.Second
	conn, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	_, err = conn.Write([]byte(resource + "\r\n"))
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
