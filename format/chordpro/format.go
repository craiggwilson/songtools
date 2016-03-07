package chordpro

import (
	"io"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
)

func init() {
	rw := &cpReaderWriter{}
	f := &format.Format{
		Name:       "chordpro",
		Reader:     rw,
		Writer:     rw,
		Extensions: []string{".cho", ".chordpro", ".chopro"},
	}

	format.Register(f)
}

type cpReaderWriter struct{}

func (cprw *cpReaderWriter) Read(r io.Reader) (*songtools.Song, error) {
	return ParseSong(r)
}

func (cprw *cpReaderWriter) Write(w io.Writer, s *songtools.Song) error {
	return WriteSong(w, s)
}

const (
	titleDirectiveName         = "title"
	subtitleDirectiveName      = "subtitle"
	keyDirectiveName           = "key"
	authorDirectiveName        = "author"
	startOfChorusDirectiveName = "start_of_chorus"
	endOfChorusDirectiveName   = "end_of_chorus"
	startOfBridgeDirectiveName = "start_of_bridge"
	endOfBridgeDirectiveName   = "end_of_bridge"
	commentDirectiveName       = "comment"
)
