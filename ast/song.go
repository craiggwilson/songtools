package ast

// Span is a representation of a span of text.
type Span struct {
	Start  uint
	Length uint
}

// Song is the root of the tree.
type Song struct {
	Span       Span
	Components []SongComponent
	Props      map[string]string
}

// SongComponent is a marker interface for components
// that are able to used in a song.
type SongComponent interface {
	songComponent()
}

func (s *Section) songComponent() {}
func (c *Comment) songComponent() {}
func (l *Line) songComponent()    {}

// Comment is a textual comment.
type Comment struct {
	Span Span
	Text []byte
}

// Section is a container for parts of a song.
type Section struct {
	Span       Span
	Components []SectionComponent
	Props      map[string]string
}

// SectionComponent is a marker interface for components
// that are able to used in a section.
type SectionComponent interface {
	sectionComponent()
}

func (c *Comment) sectionComponent() {}
func (l *Line) sectionComponent()    {}

// Line is a container for parts of a section or song.
type Line struct {
	Span       Span
	Components []LineComponent
}

// LineComponent is a marker interface for components
// that are able to be used in a line.
type LineComponent interface {
	lineComponent()
}

func (c *Comment) lineComponent() {}
func (t *Text) lineComponent()    {}
func (c *Chord) lineComponent()   {}

// Text represents raw text, generally lyrics.
type Text struct {
	Span Span
	Text []byte
}

// Chord represents a chord.
type Chord struct {
	Span Span
	Text []byte
}
