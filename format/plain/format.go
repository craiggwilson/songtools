package plain

import (
	"io"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
)

func init() {
	rw := &plainReaderWriter{}
	f := &format.Format{
		Name:   "plain",
		Reader: rw,
		Writer: rw,
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
