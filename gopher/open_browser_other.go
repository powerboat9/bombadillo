// +build !linux
// +build !darwin
// +build !windows

package gopher

import "fmt"

func OpenBrowser(url string) error {
	return fmt.Errorf("Unsupported os for browser detection")
}
