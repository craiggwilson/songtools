package html

import (
	"io"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
)

func init() {
	rw := &htmlWriter{}
	f := &format.Format{
		Name:       "html",
		Writer:     rw,
		Extensions: []string{".html", ".htm"},
	}

	format.Register(f)
}

type htmlWriter struct{}

func (hw *htmlWriter) Write(w io.Writer, s *songtools.Song) error {
	return WriteSong(w, s)
}
