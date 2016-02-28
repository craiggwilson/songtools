package chordsOverLyrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/songtools/songtools"
)

// ParseSong the src to create a songtools.Song.
func ParseSong(src io.Reader) (*songtools.Song, error) {
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

func (p *parser) parse() (*songtools.Song, error) {
	token, _, err := p.scanner.peek()
	if err != nil {
		return nil, fmt.Errorf("failed to consume initial token: %v", err)
	}
	if token == eofToken {
		return nil, nil
	}

	return p.parseSong()
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
		case commentToken:
			comment := &songtools.Comment{
				Text:   text,
				Hidden: false,
			}
			if section != nil {
				section.Nodes = append(section.Nodes, comment)
			} else {
				song.Nodes = append(song.Nodes, comment)
			}
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
				switch d.Name {
				case titleDirectiveName:
					song.Title = d.Value
				case subtitleDirectiveName:
					song.Subtitles = append(song.Subtitles, d.Value)
				case authorDirectiveName:
					song.Authors = append(song.Authors, d.Value)
				case keyDirectiveName:
					song.Key = songtools.Key(d.Value)
				default:
					song.Nodes = append(song.Nodes, d)
				}

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
			if section != nil && numNewLines == 0 {
				// we have text immediately following a section without a newline
				directive := &songtools.Comment{
					Text:   text,
					Hidden: false,
				}
				section.Nodes = append(section.Nodes, directive)
				break
			}

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
				line = nil
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

	name := strings.ToLower(parts[0])
	switch name {
	case "t":
		name = titleDirectiveName
	case "st":
		name = subtitleDirectiveName
	case "a":
		name = authorDirectiveName
	case "k":
		name = keyDirectiveName
	}

	return &songtools.Directive{
		Name:  name,
		Value: parts[1],
	}, nil
}
