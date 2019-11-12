package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Open(address string) (string, error) {
	if !pathExists(address) {
		return "", fmt.Errorf("Invalid system path: %s", address)
	}

	file, err := os.Open(address)
	if err != nil {
		return "", fmt.Errorf("Unable to open file: %s", address)
	}
	defer file.Close()

	if pathIsDir(address) {
		fileList, err := file.Readdirnames(0)
		if err != nil {
			return "", fmt.Errorf("Unable to read from directory: %s", address)
		}
		var out strings.Builder
		out.WriteString(fmt.Sprintf("Current directory: %s\n\n", address))
		for _, obj := range fileList {
			out.WriteString(obj)
			out.WriteString("\n")
		}
		return out.String(), nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Unable to read file: %s", address)
	}
	return string(bytes), nil
}

func pathExists(p string) bool {
	exists := true

	if _, err := os.Stat(p); os.IsNotExist(err) {
		exists = false
	}

	return exists
}

func pathIsDir(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return info.IsDir()
}
