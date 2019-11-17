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
		offset := 1
		fileList, err := file.Readdir(0)
		if err != nil {
			return "", links, fmt.Errorf("Unable to read from directory: %s", address)
		}
		var out strings.Builder
		out.WriteString(fmt.Sprintf("Current directory: %s\n\n", address))

		if address != "/" {
			offset = 2
			upFp := filepath.Join(address, "..")
			upOneLevel, _ := filepath.Abs(upFp)
			info, err := os.Stat(upOneLevel)
			if err == nil {
				out.WriteString("[1]   ")
				out.WriteString(fmt.Sprintf("%-12s   ", info.Mode().String()))
				out.WriteString("../\n")
				links = append(links, upOneLevel)
			}
		}

		sort.Slice(fileList, func(i, j int) bool {
			return fileList[i].Name() < fileList[j].Name()
		})

		for i, obj := range fileList {
			linkNum := fmt.Sprintf("[%d]", i+offset)
			out.WriteString(fmt.Sprintf("%-5s ", linkNum))
			out.WriteString(fmt.Sprintf("%-12s   ", obj.Mode().String()))
			out.WriteString(obj.Name())
			if obj.IsDir() {
				out.WriteString("/")
			}
			out.WriteString("\n")
			fp := filepath.Join(address, obj.Name())
			links = append(links, fp)
		}
		return out.String(), links, nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", links, fmt.Errorf("Unable to read file: %s", address)
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
