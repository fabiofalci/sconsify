package spotify

import (
	"github.com/fabiofalci/sconsify/infrastructure"
	"io/ioutil"
	"encoding/json"
	webspotify "github.com/zmb3/spotify"
)

type WebApiCache struct {
	Albums      []webspotify.SavedAlbum
	Songs       []webspotify.SavedTrack
	NewReleases []webspotify.FullPlaylist
}

func (spotify *Spotify) loadWebApiCache() *WebApiCache {
	if fileLocation := infrastructure.GetWebApiCacheFileLocation(); fileLocation != "" {
		if b, err := ioutil.ReadFile(fileLocation); err == nil {
			var webApiCache WebApiCache
			if err := json.Unmarshal(b, &webApiCache); err == nil {
				return &webApiCache
			}
		}
	}
	return &WebApiCache{}
}

func (spotify *Spotify) persistWebApiCache(webApiCache *WebApiCache) {
	if b, err := json.Marshal(webApiCache); err == nil {
		if fileLocation := infrastructure.GetWebApiCacheFileLocation(); fileLocation != "" {
			infrastructure.SaveFile(fileLocation, b)
		}
	}
}

