// +build !linux
// +build !darwin
// +build !windows

package http

import "fmt"

func OpenInBrowser(url string) (string, error) {
	return "", fmt.Errorf("Unsupported os for browser detection")
}
