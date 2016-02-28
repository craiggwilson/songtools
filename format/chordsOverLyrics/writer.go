package chordsOverLyrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/songtools/songtools"
)

// WriteSong writes a single song to the writer.
func WriteSong(w io.Writer, s *songtools.Song) error {

	if s.Title != "" {
		_, err := fmt.Fprintln(w, "#"+titleDirectiveName+"="+s.Title)
		if err != nil {
			return err
		}
	}
	for _, st := range s.Subtitles {
		_, err := fmt.Fprintln(w, "#"+subtitleDirectiveName+"="+st)
		if err != nil {
			return err
		}
	}
	for _, a := range s.Authors {
		_, err := fmt.Fprintln(w, "#"+authorDirectiveName+"="+a)
		if err != nil {
			return err
		}
	}
	if s.Key != "" {
		_, err := fmt.Fprintln(w, "#"+keyDirectiveName+"="+string(s.Key))
		if err != nil {
			return err
		}
	}

	for _, n := range s.Nodes {
		err := writeSongNode(w, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeSongNode(w io.Writer, n songtools.SongNode) error {
	switch typedN := n.(type) {
	case *songtools.Comment:
		return writeComment(w, typedN)
	case *songtools.Directive:
		return writeDirective(w, typedN)
	case *songtools.Section:
		if typedN.Kind != "" {
			_, err := fmt.Fprint(w, fmt.Sprintf("[%v]", typedN.Kind))
			if err != nil {
				return err
			}
		} else {
			// need to make sure we have an extra line here cause nothing
			// other than space is separating this from the previous
			// section
			_, err := fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}

		anyChords := len(typedN.Chords()) > 0

		for i, sn := range typedN.Nodes {
			if i == 0 {
				if c, ok := sn.(*songtools.Comment); ok {
					_, err := fmt.Fprintln(w, c.Text)
					if err != nil {
						return err
					}
					continue
				} else {
					_, err := fmt.Fprintln(w)
					if err != nil {
						return err
					}
				}
			}
			err := writeSectionNode(w, sn, anyChords)
			if err != nil {
				return err
			}
		}

		_, err := fmt.Fprintln(w)
		if err != nil {
			return err
		}
	default:
		panic("Unknown node")
	}

	return nil
}

func writeSectionNode(w io.Writer, n songtools.SectionNode, blankLineForNoChords bool) error {
	switch typedN := n.(type) {
	case *songtools.Comment:
		return writeComment(w, typedN)
	case *songtools.Directive:
		return writeDirective(w, typedN)
	case *songtools.Line:
		return writeLine(w, typedN, blankLineForNoChords)
	default:
		panic("Unknown node")
	}
}

func writeComment(w io.Writer, c *songtools.Comment) error {
	_, err := fmt.Fprintln(w, fmt.Sprintf("{%v}", c.Text))
	return err
}

func writeDirective(w io.Writer, d *songtools.Directive) error {
	_, err := fmt.Fprintln(w, fmt.Sprintf("#%v=%v", d.Name, d.Value))
	return err
}

func writeLine(w io.Writer, l *songtools.Line, blankLineForNoChords bool) error {
	if l.Chords != nil {
		pos := 0
		for i := 0; i < len(l.Chords); i++ {
			_, err := fmt.Fprint(w, strings.Repeat(" ", l.ChordPositions[i]-pos))
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(w, l.Chords[i].Name)
			if err != nil {
				return err
			}
			pos = l.ChordPositions[i] + len(l.Chords[i].Name)
		}
		_, err := fmt.Fprintln(w)
		if err != nil {
			return err
		}
	} else if blankLineForNoChords {
		_, err := fmt.Fprintln(w)
		if err != nil {
			return err
		}
	}

	_, err := fmt.Fprintln(w, l.Text)
	return err
}
