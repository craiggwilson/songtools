package plaintext

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtools"
)

// WriteSongSet writes multiple songs the same writer.
func WriteSongSet(w io.Writer, ss *songtools.SongSet) error {
	for _, s := range ss.Songs {
		err := WriteSong(w, s)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteSong writes a single song to the writer.
func WriteSong(w io.Writer, s *songtools.Song) error {
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
		if typedN.Kind == "" {
			_, err := io.WriteString(w, "\n")
			if err != nil {
				return err
			}
		} else {
			_, err := io.WriteString(w, fmt.Sprintf("\n[%v]\n", typedN.Kind))
			if err != nil {
				return err
			}
		}

		for _, sn := range typedN.Nodes {
			err := writeSectionNode(w, sn)
			if err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, "\n")
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
		return writeDirective(w, typedN)
	case *songtools.Line:
		return writeLine(w, typedN)
	default:
		panic("Unknown node")
	}
}

func writeComment(w io.Writer, c *songtools.Comment) error {
	_, err := io.WriteString(w, fmt.Sprintf("{%v}\n", c.Text))
	return err
}

func writeDirective(w io.Writer, d *songtools.Directive) error {

	if strings.ToLower(d.Name) == "comment" || strings.ToLower(d.Name) == "c" {
		_, err := io.WriteString(w, fmt.Sprintf("{%v}\n", d.Value))
		return err
	}
	_, err := io.WriteString(w, fmt.Sprintf("#%v=%v\n", d.Name, d.Value))
	return err
}

func writeLine(w io.Writer, l *songtools.Line) error {
	if l.Chords != nil {
		pos := 0
		for i := 0; i < len(l.Chords); i++ {
			_, err := io.WriteString(w, strings.Repeat(" ", l.ChordPositions[i]-pos))
			if err != nil {
				return err
			}
			_, err = io.WriteString(w, l.Chords[i].Name)
			if err != nil {
				return err
			}
			pos = l.ChordPositions[i] + len(l.Chords[i].Name)
		}
		_, err := io.WriteString(w, "\n")
		if err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, l.Text)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "\n")
	return err
}
