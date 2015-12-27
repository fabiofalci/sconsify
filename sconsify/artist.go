package sconsify

type Artist struct {
	Id         string
	Uri        string
	Name       string
}

func InitArtist(id string, uri string, name string) *Artist {
	return &Artist{
		Id:   id,
		Uri:  uri,
		Name: name,
	}
}

