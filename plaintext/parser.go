package plaintext

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtools"
)

//https://bitbucket.org/llg/songbook/src/0ba011f0a3112dd45a075d09f088df6a29981a58/song/chordpro.go?at=default&fileviewer=file-view-default
//http://www.vromans.org/johan/projects/Chordii/chordpro/index.html

// Parse the src to create a songtools.SongSet.
func Parse(src io.Reader) (*songtools.SongSet, error) {

	scanner, err := newScanner(src)
	if err != nil {
		return nil, fmt.Errorf("failed to create a scanner: %v", err)
	}

	parser := &parser{
		scanner: scanner,
	}

	return parser.parse()
}

// Parser produces songs from text.
type parser struct {

	// current token and text
	scanner *scanner
}

func (p *parser) parse() (*songtools.SongSet, error) {
	token, _, err := p.scanner.peek()
	if err != nil {
		return nil, fmt.Errorf("failed to consume initial token: %v", err)
	}
	if token == eofToken {
		return nil, nil
	}

	set := &songtools.SongSet{}
	//for {
	song, err := p.parseSong()
	if err != nil {
		return nil, fmt.Errorf("error parsing song: %v", err)
	}
	// if song == nil {
	// 	break
	// }
	set.Songs = append(set.Songs, song)
	//}

	return set, nil
}

func (p *parser) parseSong() (*songtools.Song, error) {

	song := &songtools.Song{}

	token, text, err := p.scanner.next()
	if err != nil {
		return nil, err
	}

	var section *songtools.Section
	var line *songtools.Line
	numNewLines := 0

	for token != eofToken {
		switch token {
		case directiveToken:
			d, err := parseDirective(text)
			if err != nil {
				return nil, err
			}

			if numNewLines == 2 {
				section = nil
			}

			if section != nil {
				section.Nodes = append(section.Nodes, d)
			} else {
				song.Nodes = append(song.Nodes, d)
			}
			line = nil
			numNewLines = 0
		case sectionToken:
			section = &songtools.Section{
				Kind: songtools.SectionKind(text),
			}

			song.Nodes = append(song.Nodes, section)
			line = nil
			numNewLines = 0
		case textToken:
			if section == nil {
				section = &songtools.Section{}
				song.Nodes = append(song.Nodes, section)
			}

			chords, positions, isChordLine := songtools.ParseTextForChords(text)
			if !isChordLine {
				if line != nil && line.Text == "" {
					line.Text = text
				} else {
					line = &songtools.Line{
						Text: text,
					}
					section.Nodes = append(section.Nodes, line)
				}
			} else {
				line = &songtools.Line{
					Chords:         chords,
					ChordPositions: positions,
				}
				section.Nodes = append(section.Nodes, line)
			}
			numNewLines = 0
		case newLineToken:
			numNewLines++
			if numNewLines > 2 {
				section = nil
				line = nil
				numNewLines = 0
			}
			break
		}

		token, text, err = p.scanner.next()
		if err != nil {
			return nil, err
		}
	}

	return song, nil
}

func parseDirective(text string) (*songtools.Directive, error) {
	parts := strings.SplitN(text, "=", 2)

	if len(parts) != 2 {
		return nil, fmt.Errorf("directives must be a key value pair with an '=' as the seperator: %v", text)
	}

	return &songtools.Directive{parts[0], parts[1]}, nil
}
