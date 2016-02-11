package songtools

import "strings"

// Song is a set of attributes, chords, and lyrics.
type Song struct {
	Attributes map[string]string
	Sections   []*Section
}

func NewSong() *Song {
	return &Song{
		Attributes: make(map[string]string),
		Sections:   []*Section{},
	}
}

func (s *Song) AddSection(kind SectionKind) *Section {
	sec := newSection(kind)
	s.Sections = append(s.Sections, sec)
	return sec
}

func (s *Song) SetAttribute(name, value string) {
	s.Attributes[strings.ToLower(name)] = value
}

type SectionKind string

type Section struct {
	Kind  SectionKind
	Lines []*Line
}

func newSection(kind SectionKind) *Section {
	return &Section{
		Kind:  kind,
		Lines: []*Line{},
	}
}

func (s *Section) AddLine(text string) *Line {
	line := newLine(text)
	s.Lines = append(s.Lines, line)
	return line
}

type Line struct {
	Text           string
	Chords         []*Chord
	ChordPositions []int
}

func newLine(text string) *Line {
	return &Line{
		Text: text,
	}
}

// Chord is a named set of notes.
type Chord struct {
	Name   string
	Root   Note
	Suffix string
}

func newChord(name string, root Note, suffix string) *Chord {
	return &Chord{
		Name:   name,
		Root:   root,
		Suffix: suffix,
	}
}

func (c *Chord) String() string {
	return c.Name
}

// Note is a single note on a scale.
type Note int

const (
	noteCount = 12
	A         = Note(0)
	ASharp    = Note(1)
	B         = Note(2)
	C         = Note(3)
	CSharp    = Note(4)
	D         = Note(5)
	DSharp    = Note(6)
	E         = Note(7)
	F         = Note(8)
	FSharp    = Note(9)
	G         = Note(10)
	GSharp    = Note(11)
)

type ChordNames [noteCount]string

var (
	SharpChordNames = ChordNames{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
	FlatChordNames  = ChordNames{"A", "Bb", "B", "C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab"}
)

func (n Note) String() string {
	return n.StringFromNames(SharpChordNames)
}

func (n Note) StringFromNames(names ChordNames) string {
	return names[n]
}

func (n Note) Interval(interval int) Note {
	return Note((int(n) + interval) % noteCount)
}

func ParseChord(text string) (*Chord, bool) {
	if text == "" {
		return nil, false
	}

	note, ok := parseNote(text)
	if !ok {
		return nil, false
	}

	suffix := ""
	offset := len(note.String())
	if len(text) > offset {
		suffix = text[offset:]
	}

	return newChord(text, note, suffix), true
}

func parseNote(text string) (Note, bool) {
	if text == "" {
		return -1, false
	}

	// this is pretty slow... We can do better.
	for len(text) > 0 {
		for i := 0; i < noteCount; i++ {
			if text == SharpChordNames[i] || text == FlatChordNames[i] {
				return Note(i), true
			}
		}

		text = text[0 : len(text)-1]
	}

	return -1, false
}
