package songtool

// Song is a set of attributes, chords, and lyrics.
type Song struct {
	Attributes map[string]string
	Sections   []Section
}

type SectionKind string

type Section struct {
	Kind  SectionKind
	Lines []Line
}

type Line struct {
	Text           string
	Chords         []Chord
	ChordPositions []int
}
