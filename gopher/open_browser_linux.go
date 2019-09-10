// +build linux

package gopher

import "os/exec"

func OpenBrowser(url string) error {
	return exec.Command("xdg-open", url).Start()
}
