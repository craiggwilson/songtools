package songtools

const (
	// The directive name for a title.
	TitleDirectiveName = "title"
	// The directive name for a key.
	KeyDirectiveName = "key"
	// The directive name for a comment.
	CommentDirectiveName = "comment"
)

// Song is a set of nodes.
type Song struct {
	Nodes []SongNode
}

// Key gets the key the song is in. It will autodetect if necessary.
func (s *Song) Key() (Key, bool) {
	for _, n := range s.Nodes {
		switch typedN := n.(type) {
		case *Directive:
			if typedN.Name == KeyDirectiveName {
				return Key(typedN.Value), true
			}
		}
	}

	return "", false
}

// Chords gets all the chords present in the song.
func (s *Song) Chords() []*Chord {
	chords := []*Chord{}

	for _, n := range s.Nodes {
		switch typedN := n.(type) {
		case *Section:
			chords = append(chords, typedN.Chords()...)
		}
	}

	return chords
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

// Chords gets all the chords present in the section.
func (s *Section) Chords() []*Chord {
	chords := []*Chord{}
	for _, n := range s.Nodes {
		switch typedN := n.(type) {
		case *Line:
			chords = append(chords, typedN.Chords...)
		}
	}

	return chords
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
