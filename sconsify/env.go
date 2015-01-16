package sconsify

import (
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

func DeleteCache(cacheLocation string) {
	if strings.HasSuffix(cacheLocation, SCONSIFY_CONF_LOCATION+"/cache") {
		os.RemoveAll(cacheLocation)
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
