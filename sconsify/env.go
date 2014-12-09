package sconsify

import (
	"os"
	"strings"

	"github.com/fabiofalci/flagrc"
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
	flagrc.ProcessFlagrc(basePath + "/sconsifyrc")
}
