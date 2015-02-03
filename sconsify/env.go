package sconsify

import (
	"errors"
	"os"
	"strings"

	"github.com/fabiofalci/flagrc"
	"github.com/mitchellh/go-homedir"
)

const SCONSIFY_CONF_LOCATION = "/.sconsify"

func GetCacheLocation() string {
	if basePath := getConfLocation(); basePath != "" {
		return basePath + "/cache"
	}
	return ""
}

func DeleteCache(cacheLocation string) error {
	if strings.HasSuffix(cacheLocation, SCONSIFY_CONF_LOCATION+"/cache") {
		return os.RemoveAll(cacheLocation)
	}
	return errors.New("Invalid cache location: " + cacheLocation)
}

func GetLogFileLocation() string {
	if basePath := getConfLocation(); basePath != "" {
		return basePath + "/sconsify.log"
	}
	return ""
}

func DeleteLogFile(logFileLocation string) error {
	if strings.HasSuffix(logFileLocation, SCONSIFY_CONF_LOCATION+"/sconsify.log") {
		return os.Remove(logFileLocation)
	}
	return errors.New("Invalid log location: " + logFileLocation)
}

func GetStateFileLocation() string {
	if basePath := getConfLocation(); basePath != "" {
		return basePath + "/state.json"
	}
	return ""
}

func SaveFile(fileLocation string, content []byte) {
	file, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer file.Close()
		file.Write(content)
	}
}

func getConfLocation() string {
	if dir, err := homedir.Dir(); err == nil {
		if dir, err = homedir.Expand(dir); err == nil && dir != "" {
			return dir + SCONSIFY_CONF_LOCATION
		}
	}
	return ""
}

func ProcessSconsifyrc() {
	if basePath := getConfLocation(); basePath != "" {
		flagrc.ProcessFlagrc(basePath + "/sconsifyrc")
	}
}
