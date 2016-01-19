package sconsify

import "strings"

type Artist struct {
	URI string

	Name   string
	Albums []*Album
}

func InitArtist(URI string, name string) *Artist {
	return &Artist{
		URI:  URI,
		Name: name,
	}
}

func (artist *Artist) GetSpotifyID() string {
	return artist.URI[strings.LastIndex(artist.URI, ":")+1 : len(artist.URI)]
}
