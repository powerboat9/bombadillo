// +build linux

package http

import "os/exec"

func OpenInBrowser(url string) (string, error) {
  // Check for a local display server, this is
  // not a silver bullet but should help ssh
  // connected users on many systems get accurate
  // messaging and not spin off processes needlessly
	err := exec.Command("type", "Xorg").Run()
	if err != nil {
		return "", fmt.Errorf("No gui is available, check 'webmode' setting")
	}

  // Use start rather than run or output in order
  // to release the process and not block
	err := exec.Command("xdg-open", url).Start()
	if err != nil {
		return "", err
	}
	return "Opened in system default web browser", nil
}
