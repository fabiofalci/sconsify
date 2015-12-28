package sconsify

type Artist struct {
	ID     int
	URI    string
	SID    string

	Name       string
	Albums     []*Album
}

func InitArtist(SID string, URI string, name string) *Artist {
	return &Artist{
		URI:  URI,
		Name: name,
		SID: SID,
	}
}

