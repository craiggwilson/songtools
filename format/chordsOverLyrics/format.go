package chordsOverLyrics

import (
	"io"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
)

func init() {
	rw := &plainReaderWriter{}
	f := &format.Format{
		Name:       "chordsOverLyrics",
		Reader:     rw,
		Writer:     rw,
		Extensions: []string{"txt"},
	}

	format.Register(f)
}

type plainReaderWriter struct{}

func (prw *plainReaderWriter) Read(r io.Reader) (*songtools.Song, error) {
	return ParseSong(r)
}

func (prw *plainReaderWriter) Write(w io.Writer, s *songtools.Song) error {
	return WriteSong(w, s)
}

const (
	titleDirectiveName    = "title"
	subtitleDirectiveName = "subtitle"
	keyDirectiveName      = "key"
	authorDirectiveName   = "author"
)
