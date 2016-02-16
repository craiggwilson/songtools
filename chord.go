package songtools

// Chord is a named set of notes.
type Chord struct {
	Name   string
	Root   Note
	Suffix string
}

func (c *Chord) String() string {
	return c.Name
}

// Interval returns a new chord at the specified interval.
func (c *Chord) Interval(interval int, names NoteNames) *Chord {
	newRoot := c.Root.Interval(interval)
	newName := newRoot.StringFromNames(names) + c.Suffix

	return &Chord{newName, newRoot, c.Suffix}
}

// Note is a single note on a scale.
type Note int

const (
	noteCount = 12
)

// NoteNames is a dictionary for looking up a note name from it's chromatic number.
type NoteNames [noteCount]string

var (
	sharpNoteNames = NoteNames{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
	flatNoteNames  = NoteNames{"A", "Bb", "B", "C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab"}
)

func (n Note) String() string {
	return n.StringFromNames(sharpNoteNames)
}

// StringFromNames returns the name of a note.
func (n Note) StringFromNames(names NoteNames) string {
	return names[n]
}

// Interval returns a new note at the specified interval.
func (n Note) Interval(interval int) Note {
	return Note((int(n) + interval) % noteCount)
}

// ParseChord parses some text and returns a chord.
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
			if text == sharpNoteNames[i] || text == flatNoteNames[i] {
				return Note(i), true
			}
		}

		text = text[0 : len(text)-1]
	}

	return -1, false
}
