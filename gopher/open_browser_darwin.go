// +build darwin

package gopher

import "os/exec"

func OpenBrowser(url string) error {
	return exec.Command("open", url).Start()
}
