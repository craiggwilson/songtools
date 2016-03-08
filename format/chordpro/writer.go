package chordpro

import (
	"fmt"
	"io"
	"strings"

	"github.com/songtools/songtools"
)

// WriteSong writes a single song to the writer.
func WriteSong(w io.Writer, s *songtools.Song) error {

	if s.Title != "" {
		err := writeDirective(w, titleDirectiveName, s.Title)
		if err != nil {
			return err
		}
	}
	for _, st := range s.Subtitles {
		err := writeDirective(w, subtitleDirectiveName, st)
		if err != nil {
			return err
		}
	}
	for _, a := range s.Authors {
		err := writeDirective(w, authorDirectiveName, a)
		if err != nil {
			return err
		}
	}
	if s.Key != "" {
		err := writeDirective(w, keyDirectiveName, string(s.Key))
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
		return writeDirective(w, typedN.Name, typedN.Value)
	case *songtools.Section:
		if typedN.Kind != "" {
			switch typedN.Kind {
			case "Chorus":
				writeDirective(w, startOfChorusDirectiveName, "")
			case "Bridge":
				writeDirective(w, startOfBridgeDirectiveName, "")
			default:
				writeDirective(w, commentDirectiveName, string(typedN.Kind))
			}
		}

		for _, sn := range typedN.Nodes {
			err := writeSectionNode(w, sn)
			if err != nil {
				return err
			}
		}

		switch typedN.Kind {
		case "Chorus":
			writeDirective(w, endOfChorusDirectiveName, "")
		case "Bridge":
			writeDirective(w, endOfBridgeDirectiveName, "")
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

func writeSectionNode(w io.Writer, n songtools.SectionNode) error {
	switch typedN := n.(type) {
	case *songtools.Comment:
		return writeComment(w, typedN)
	case *songtools.Directive:
		return writeDirective(w, typedN.Name, typedN.Value)
	case *songtools.Line:
		return writeLine(w, typedN)
	default:
		panic("Unknown node")
	}
}

func writeComment(w io.Writer, c *songtools.Comment) error {
	if c.Hidden {
		_, err := fmt.Fprintln(w, fmt.Sprintf("#%v", c.Text))
		return err
	}

	return writeDirective(w, commentDirectiveName, c.Text)
}

func writeDirective(w io.Writer, name, value string) error {
	if value != "" {
		_, err := fmt.Fprintln(w, fmt.Sprintf("{%v:%v}", name, value))
		return err
	}

	_, err := fmt.Fprintln(w, fmt.Sprintf("{%v}", name))
	return err

}

func writeLine(w io.Writer, l *songtools.Line) error {
	if l.Chords != nil {
		pos := 0
		for i := 0; i < len(l.Chords); i++ {
			if pos < l.ChordPositions[i] {
				if pos < len(l.Text) {
					if len(l.Text) < l.ChordPositions[i] {
						_, err := fmt.Fprint(w, l.Text[pos:len(l.Text)])
						if err != nil {
							return err
						}
						pos = len(l.Text)
					} else {
						_, err := fmt.Fprint(w, l.Text[pos:l.ChordPositions[i]])
						if err != nil {
							return err
						}
						pos = l.ChordPositions[i]
					}
				}

				if pos >= l.ChordPositions[i] {
					_, err := fmt.Fprint(w, strings.Repeat(" ", l.ChordPositions[i]-pos))
					if err != nil {
						return err
					}
					pos = l.ChordPositions[i]
				}
			}

			_, err := fmt.Fprint(w, "["+l.Chords[i].Name+"]")
			if err != nil {
				return err
			}
		}

		if pos < len(l.Text) {
			_, err := fmt.Fprint(w, l.Text[pos:])
			if err != nil {
				return err
			}
		}

		_, err := fmt.Fprintln(w)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprintln(w, l.Text)
		if err != nil {
			return err
		}
	}

	return nil
}
