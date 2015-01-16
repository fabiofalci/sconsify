package spotify

import "io/ioutil"

func getKey() ([]byte, error) {
	if len(key) == 0 {
		appKey, err := ioutil.ReadFile("spotify_appkey.key")
		if err != nil {
			return nil, err
		}
		return appKey, nil
	}
	return key, nil
}

// when building the last line will get replaced by deploying key
var key = []byte{}
