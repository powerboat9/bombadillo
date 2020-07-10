// This will build for osx without a build tag based on the filename

package http

import "os/exec"

func OpenInBrowser(url string) (string, error) {
	err := exec.Command("open", url).Start()
	if err != nil {
		return "", err
	}
	return "Opened in system default web browser", nil
}
