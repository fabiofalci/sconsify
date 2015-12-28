package sconsify

type Album struct {
	ID     int
	URI    string

	Name       string
	Artists    []*Artist
}

