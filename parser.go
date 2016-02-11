package songtools

import "io"

// SongParser is an interface representing a parser for a song.
type SongParser interface {
	// Parse returns a song or an error after consuming the reader.
	Parse(io.Reader) ([]Song, error)
}
