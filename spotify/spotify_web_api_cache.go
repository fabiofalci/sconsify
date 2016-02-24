package spotify

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/fabiofalci/sconsify/infrastructure"
	webspotify "github.com/zmb3/spotify"
	"io"
	"io/ioutil"
)

type WebApiCache struct {
	Albums      []webspotify.SavedAlbum
	Songs       []webspotify.SavedTrack
	NewReleases []webspotify.FullPlaylist
}

func (spotify *Spotify) loadWebApiCache() *WebApiCache {
	if spotify.cacheWebApiContent {
		if fileLocation := infrastructure.GetWebApiCacheFileLocation(); fileLocation != "" {
			if b, err := ioutil.ReadFile(fileLocation); err == nil {
				compressed := bytes.NewBuffer(b)
				if r, err := gzip.NewReader(compressed); err == nil {
					var uncompressed bytes.Buffer
					io.Copy(&uncompressed, r)
					r.Close()
					var webApiCache WebApiCache
					if err := json.Unmarshal(uncompressed.Bytes(), &webApiCache); err == nil {
						return &webApiCache
					}
				}
			}
		}
	}
	return &WebApiCache{}
}

func (spotify *Spotify) persistWebApiCache(webApiCache *WebApiCache) {
	if spotify.cacheWebApiContent {
		if b, err := json.Marshal(webApiCache); err == nil {
			var compressed bytes.Buffer
			w := gzip.NewWriter(&compressed)
			w.Write([]byte(b))
			w.Close()
			if fileLocation := infrastructure.GetWebApiCacheFileLocation(); fileLocation != "" {
				infrastructure.SaveFile(fileLocation, compressed.Bytes())
			}
		}
	}
}
