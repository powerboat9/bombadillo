// +build windows

package http

import "os/exec"

func OpenInBrowser(url string) (string, error) {
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	if err != nil {
		return "", err
	}
	return "Opened in system default web browser", nil
}
