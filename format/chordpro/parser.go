package chordpro

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
				Hidden: true,
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

			switch d.Name {
			case startOfChorusDirectiveName:
				section = &songtools.Section{
					Kind: songtools.SectionKind("Chorus"),
				}
				song.Nodes = append(song.Nodes, section)
			case endOfChorusDirectiveName:
				section = nil
			case startOfBridgeDirectiveName:
				section = &songtools.Section{
					Kind: songtools.SectionKind("Bridge"),
				}
				song.Nodes = append(song.Nodes, section)
			case endOfBridgeDirectiveName:
				section = nil
			case commentDirectiveName:
				comment := &songtools.Comment{
					Text:   d.Value,
					Hidden: false,
				}
				if section == nil {
					song.Nodes = append(song.Nodes, comment)
				} else {
					section.Nodes = append(section.Nodes, comment)
				}
			case titleDirectiveName:
				song.Title = d.Value
			case subtitleDirectiveName:
				song.Subtitles = append(song.Subtitles, d.Value)
			case authorDirectiveName:
				song.Authors = append(song.Authors, d.Value)
			case keyDirectiveName:
				song.Key = songtools.Key(d.Value)
			default:
				if section != nil && numNewLines == 2 && section.Kind != "Chorus" && section.Kind != "Bridge" {
					section = nil
				}

				if section != nil {
					section.Nodes = append(section.Nodes, d)
				} else {
					song.Nodes = append(song.Nodes, d)
				}
				line = nil
				numNewLines = 0
			}
		case chordToken:
			if section == nil {
				section = &songtools.Section{}
				song.Nodes = append(song.Nodes, section)
			}

			if line == nil {
				line = &songtools.Line{}
			}

			chord, ok := songtools.ParseChord(text)
			if !ok {
				return nil, fmt.Errorf("The text '%v' is not a chord.", text)
			}

			line.Chords = append(line.Chords, chord)
			line.ChordPositions = append(line.ChordPositions, len(text))

		case textToken:
			if section == nil {
				section = &songtools.Section{}
				song.Nodes = append(song.Nodes, section)
			}

			if line == nil {
				line = &songtools.Line{}
			}

			line.Text += text

			section.Nodes = append(section.Nodes, line)
			numNewLines = 0
		case newLineToken:
			line = nil
			numNewLines++
			if section != nil && numNewLines == 2 && section.Kind != "Chorus" && section.Kind != "Bridge" {
				section = nil
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
	parts := strings.SplitN(text, ":", 2)

	if len(parts) > 2 {
		return nil, fmt.Errorf("directives must either have no value or have a value separated by a ':': %v", text)
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
	case "soc":
		name = "start_of_chorus"
	case "eoc":
		name = "end_of_chorus"
	case "sob":
		name = "start_of_bridge"
	case "eob":
		name = "end_of_bridge"
	}

	value := ""
	if len(parts) > 1 {
		value = parts[1]
	}

	return &songtools.Directive{
		Name:  name,
		Value: value,
	}, nil
}
