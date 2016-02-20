package plain

import (
	"io"

	"github.com/songtools/songtools"
)

func init() {
	rw := &songReaderWriter{}
	f := &songtools.Format{
		Name:   "plain",
		Reader: rw,
		Writer: rw,
	}

	songtools.RegisterFormat(f)
}

type songReaderWriter struct{}

func (srw *songReaderWriter) Read(r io.Reader) (*songtools.SongSet, error) {
	return ParseSongSet(r)
}

func (srw *songReaderWriter) Write(w io.Writer, set *songtools.SongSet) error {
	return WriteSongSet(w, set)
}
