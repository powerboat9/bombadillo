package local

import (
	"fmt"
	"io/ioutil"
	"os"
  "path/filepath"
  "sort"
	"strings"
)

func Open(address string) (string, []string, error) {
  links := make([]string, 0, 10)

	if !pathExists(address) {
		return "", links, fmt.Errorf("Invalid system path: %s", address)
	}

	file, err := os.Open(address)
	if err != nil {
		return "", links, fmt.Errorf("Unable to open file: %s", address)
	}
	defer file.Close()

	if pathIsDir(address) {
		fileList, err := file.Readdir(0)
		if err != nil {
			return "", links, fmt.Errorf("Unable to read from directory: %s", address)
		}
		var out strings.Builder
		out.WriteString(fmt.Sprintf("Current directory: %s\n\n", address))

    sort.Slice(fileList, func(i, j int) bool {
      return fileList[i].Name() < fileList[j].Name()
    })

		for i, obj := range fileList {
      linkNum := fmt.Sprintf("[%d]", i+1)
      out.WriteString(fmt.Sprintf("%-5s ", linkNum))
      out.WriteString(fmt.Sprintf("%s    ", obj.Mode().String()))
			out.WriteString(obj.Name())
			out.WriteString("\n")
      fp := filepath.Join(address, obj.Name())
      links = append(links, fp)
		}
		return out.String(), links, nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", links ,fmt.Errorf("Unable to read file: %s", address)
	}
	return string(bytes), links, nil
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
