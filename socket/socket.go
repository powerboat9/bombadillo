package socket

import (
	"net"
	"io/ioutil"
	"gsock/gopher"
	"errors"
)



func Retrieve(u gopher.Url) ([]byte, error) {
  nullRes := make([]byte, 0)
  if u.Host == "" || u.Port == "" {
		return nullRes, errors.New("Incomplete request url")
  }

  addr := u.Host + ":" + u.Port
  tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
  if err != nil {
		errortext := "Could not find host: " + u.Full
		return nullRes, errors.New(errortext)
  }

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nullRes, err
	}

	send := u.Resource + "\n"
	if u.Scheme == "http" || u.Scheme == "https" {
		send = u.Gophertype
	}

	_, err = conn.Write([]byte(send))
	if err != nil {
		return nullRes, err
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nullRes, err
	}

	return result, err
}


