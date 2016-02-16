package songtools

// SongSet is a set of songs.
type SongSet struct {
	Songs []*Song
}

// Song is a set of nodes.
type Song struct {
	Nodes []SongNode
}

// SongNode represents a node that can appear in a song.
type SongNode interface {
	songNode()
}

func (c *Comment) songNode()   {}
func (d *Directive) songNode() {}
func (s *Section) songNode()   {}

// Directive contains a name and value associated with either
// a Song or a Section.
type Directive struct {
	Name  string
	Value string
}

// Comment contains text that represents a comment not intended
// for output in print form.
type Comment struct {
	Text string
}

// SectionKind is the type of section. Examples are Chorus, Verse, and Bridge.
type SectionKind string

// Section contains nodes
type Section struct {
	Kind  SectionKind
	Nodes []SectionNode
}

// SectionNode represents a node that can appear in a section.
type SectionNode interface {
	sectionNode()
}

func (c *Comment) sectionNode()   {}
func (d *Directive) sectionNode() {}
func (l *Line) sectionNode()      {}

// Line represents a lyric line and/or chords.
type Line struct {
	Text           string
	Chords         []*Chord
	ChordPositions []int
}

// Chord is a named set of notes.
type Chord struct {
	Name   string
	Root   Note
	Suffix string
}

func (c *Chord) String() string {
	return c.Name
}

// Note is a single note on a scale.
type Note int

const (
	noteCount = 12
)

// ChordNames is a dictionary for looking up a chord name from it's chromatic number.
type ChordNames [noteCount]string

var (
	sharpChordNames = ChordNames{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
	flatChordNames  = ChordNames{"A", "Bb", "B", "C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab"}
)

func (n Note) String() string {
	return n.StringFromNames(sharpChordNames)
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

	return &Chord{
		Name:   text,
		Root:   note,
		Suffix: suffix,
	}, true
}

func parseNote(text string) (Note, bool) {
	if text == "" {
		return -1, false
	}

	// this is pretty slow... We can do better.
	for len(text) > 0 {
		for i := 0; i < noteCount; i++ {
			if text == sharpChordNames[i] || text == flatChordNames[i] {
				return Note(i), true
			}
		}

		text = text[0 : len(text)-1]
	}

	return -1, false
}
