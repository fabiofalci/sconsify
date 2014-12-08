package sconsify

import (
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const SCONSIFY_CONF_LOCATION = "/.sconsify"

func GetCacheLocation() *string {
	dir, err := homedir.Dir()
	if err == nil {
		dir, err = homedir.Expand(dir)
		if err == nil && dir != "" {
			path := dir + SCONSIFY_CONF_LOCATION + "/cache/"
			return &path
		}
	}
	return nil
}

func DeleteCache(cacheLocation *string) {
	if strings.HasSuffix(*cacheLocation, SCONSIFY_CONF_LOCATION) {
		os.RemoveAll(*cacheLocation)
	}
}
