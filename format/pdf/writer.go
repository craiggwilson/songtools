package pdf

import (
	"io"

	"github.com/jung-kurt/gofpdf"
	"github.com/songtools/songtools"
)

const (
	titlePointSize = 24
)

// WriteSong writes a single song to the writer.
func WriteSong(w io.Writer, s *songtools.Song) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetTitle(s.Title, true)

	writeHeader(pdf, s)

	err := pdf.OutputFileAndClose("C:\\projects\\go\\src\\github.com\\songtools\\songtools\\test.pdf")
	if err != nil {
		return err
	}
	pdf.Close()
	return nil
}

func writeHeader(pdf *gofpdf.Fpdf, s *songtools.Song) {
	pdf.SetFont("Arial", "B", titlePointSize)
	ht := pdf.PointConvert(titlePointSize)
	pdf.CellFormat(0, ht, s.Title, "B", 1, "", false, 0, "")
	pdf.Ln(ht)
}
