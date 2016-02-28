package pdf

import (
	"io"

	"github.com/songtools/songtools"
	"github.com/songtools/songtools/format"
)

func init() {
	rw := &pdfWriter{}
	f := &format.Format{
		Name:   "pdf",
		Writer: rw,
	}

	format.Register(f)
}

type pdfWriter struct{}

func (pw *pdfWriter) Write(w io.Writer, s *songtools.Song) error {
	return WriteSong(w, s)
}
