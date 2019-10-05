package local

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Open(address string) (string, error) {
	file, err := os.Open(address)
	if err != nil {
		return "", fmt.Errorf("Unable to open file: %s", address)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Unable to read file: %s", address)
	}
	return string(bytes), nil
}
