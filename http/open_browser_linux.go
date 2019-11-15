// +build linux

package http

import (
	"fmt"
	"os"
	"os/exec"
)

// OpenInBrowser checks for the presence of a display server
// and environment variables indicating a gui is present. If found
// then xdg-open is called on a url to open said url in the default
// gui web browser for the system
func OpenInBrowser(url string) (string, error) {
	disp := os.Getenv("DISPLAY")
	wayland := os.Getenv("WAYLAND_DISPLAY")
	_, err := exec.LookPath("Xorg")
	if disp == "" && wayland == "" && err != nil {
		return "", fmt.Errorf("No gui is available, check 'webmode' setting")
	}

	// Use start rather than run or output in order
	// to release the process and not block
	err = exec.Command("xdg-open", url).Start()
	if err != nil {
		return "", err
	}
	return "Opened in system default web browser", nil
}
