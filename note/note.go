package note

// Note is a single note in a scale.
type Note uint8

// Namer provides the name for note.
type Namer interface {
	Name(Note) (string, error)
}

// Parser parses the string into a note.
type Parser interface {
	Parse(string) (Note, uint8, error)
}

// Language handles naming and parsing notes.
type Language interface {
	Namer
	Parser
}
