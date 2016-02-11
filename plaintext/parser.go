package plaintext

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtools"
)

//https://bitbucket.org/llg/songbook/src/0ba011f0a3112dd45a075d09f088df6a29981a58/song/chordpro.go?at=default&fileviewer=file-view-default
//http://www.vromans.org/johan/projects/Chordii/chordpro/index.html

func Parse(src io.Reader) ([]*songtools.Song, error) {

	scanner, err := NewScanner(src)
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
	scanner *Scanner
}

func (p *parser) parse() ([]*songtools.Song, error) {
	token, _, err := p.scanner.Peek()
	if err != nil {
		return nil, fmt.Errorf("failed to consume initial token: %v", err)
	}
	if token == EofToken {
		return nil, nil
	}

	songs := []*songtools.Song{}
	//for {
	song, err := p.parseSong()
	if err != nil {
		return nil, fmt.Errorf("error parsing song: %v", err)
	}
	// if song == nil {
	// 	break
	// }
	songs = append(songs, song)
	//}

	return songs, nil
}

func (p *parser) parseSong() (*songtools.Song, error) {
	song := songtools.NewSong()

	var lastToken Token
	token, text, err := p.scanner.Next()
	if err != nil {
		return nil, err
	}

	var section *songtools.Section
	var line *songtools.Line

	for token != EofToken {
		switch token {
		case DirectiveToken:
			parts, err := parseDirective(text)
			if err != nil {
				return nil, err
			}
			song.SetAttribute(parts[0], parts[1])
		case SectionToken:
			section = song.AddSection(songtools.SectionKind(text))
		case TextToken:
			if section == nil {
				section = song.AddSection(songtools.SectionKind("Verse"))
			}

			chords, positions, isChordLine := parseLine(text)
			if !isChordLine {
				println(text)
				if line != nil && line.Text == "" {
					line.Text = text
				} else {
					section.AddLine(text)
				}
			} else {
				line = section.AddLine("")
				line.Chords = chords
				line.ChordPositions = positions
			}
		case NewLineToken:
			if lastToken == NewLineToken {
				section = nil
			}
			break
		}

		lastToken = token
		token, text, err = p.scanner.Next()
		if err != nil {
			return nil, err
		}
	}

	return song, nil
}

func parseDirective(text string) ([]string, error) {
	parts := strings.SplitN(text, "=", 2)

	if len(parts) < 2 {
		return nil, fmt.Errorf("directives must be a key value pair with an '=' as the seperator: %v", text)
	}

	return parts, nil
}

func parseLine(text string) ([]*songtools.Chord, []int, bool) {

	chords := []*songtools.Chord{}
	positions := []int{}

	i := 0
	for i < len(text) {
		for i < len(text) && text[i] == ' ' {
			i++
		}

		name := ""
		pos := i
		for i < len(text) && text[i] != ' ' {
			name += string(text[i])
			i++
		}

		chord, ok := songtools.ParseChord(name)
		if !ok {
			println(name)
			// we aren't a chord line
			return nil, nil, false
		}

		chords = append(chords, chord)
		positions = append(positions, pos)
	}

	return chords, positions, true
}
