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
