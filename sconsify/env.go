package sconsify

import (
	"bufio"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const SCONSIFY_CONF_LOCATION = "/.sconsify"

func GetCacheLocation() string {
	basePath := getConfLocation()
	if basePath != "" {
		return basePath + "/cache"
	}
	return ""
}

func DeleteCache(cacheLocation string) {
	if strings.HasSuffix(cacheLocation, SCONSIFY_CONF_LOCATION+"/cache") {
		os.RemoveAll(cacheLocation)
	}
}

func getConfLocation() string {
	dir, err := homedir.Dir()
	if err == nil {
		dir, err = homedir.Expand(dir)
		if err == nil && dir != "" {
			return dir + SCONSIFY_CONF_LOCATION
		}
	}
	return ""
}

func ProcessSconsifyrc() {
	basePath := getConfLocation()
	if basePath == "" {
		return
	}

	file, err := os.Open(basePath + "/sconsifyrc")
	if err != nil {
		return

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		key := line
		if index := strings.Index(line, "="); index > 0 {
			key = line[0:index]
		}
		contains := false
		for i, value := range os.Args {
			if i > 0 {
				if strings.HasPrefix(value, key) {
					contains = true
					break
				}
			}
		}

		if !contains {
			os.Args = append(os.Args, line)
		}
	}
}
