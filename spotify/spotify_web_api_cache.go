package spotify

import (
	webspotify "github.com/zmb3/spotify"
	"github.com/fabiofalci/sconsify/infrastructure"
	"io/ioutil"
	"encoding/json"
)

type WebApiCache struct {
	Albums      []CachedAlbum
	Songs       []CachedTrack
	NewReleases []webspotify.SimplePlaylist
	Artists     []webspotify.FullArtist

	SharedArtists []CachedArtist
}

type CachedAlbum struct {
	URI    string
	Name   string
	Tracks []CachedTrack
}

type CachedTrack struct {
	URI          string
	Name         string
	TimeDuration string
	ArtistsURI   []string
}

type CachedArtist struct {
	URI  string
	Name string
}

func (webApiCache *WebApiCache) findSharedArtist(URI string) *CachedArtist {
	for _, cachedArtist := range webApiCache.SharedArtists {
		if cachedArtist.URI == URI {
			return &cachedArtist
		}
	}
	return nil
}

func (webApiCache *WebApiCache) addSharedArtist(artist CachedArtist) {
	if webApiCache.SharedArtists == nil {
		webApiCache.SharedArtists = make([]CachedArtist, 0)
	}
	webApiCache.SharedArtists = append(webApiCache.SharedArtists, artist)
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

