package songtools

import (
	"fmt"
	"strings"
)

// Chord is a named set of notes.
type Chord struct {
	Name   string
	Root   Note
	Base   Note
	Suffix string
}

func (c *Chord) String() string {
	return c.Name
}

// Interval returns a new chord at the specified interval.
func (c *Chord) Interval(interval int, names *NoteNames) *Chord {
	newRoot := c.Root.Interval(interval)
	newBase := c.Base.Interval(interval)
	newName := newRoot.StringFromNames(names) + c.Suffix
	if newBase != newRoot {
		newName += "/" + newBase.StringFromNames(names)
	}

	return &Chord{
		Name:   newName,
		Root:   newRoot,
		Base:   newBase,
		Suffix: c.Suffix,
	}
}

// Key is...
type Key string

var (
	sharpKeys = []Key{"A", "F#m", "B", "G#m", "C", "Am", "C#", "A#m", "D", "Bm", "E", "C#m", "F#", "D#m", "G", "Em"}
	flatKeys  = []Key{"Ab", "Fm", "Bb", "Gm", "Cb", "Abm", "Db", "Bbm", "Eb", "Cm", "F", "Dm", "Gb", "Ebm"}
)

// Note is a single note on a scale.
type Note int

const (
	noteCount = 12
)

// NoteNames is a dictionary for looking up a note name from it's chromatic number.
type NoteNames [noteCount]string

var (
	sharpNoteNames = &NoteNames{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
	flatNoteNames  = &NoteNames{"A", "Bb", "B", "C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab"}
)

// NoteNamesAndIntervalFromKeyToKey returns the NoteNames and the interval in order to transposed
// the original key to the transposed key.
func NoteNamesAndIntervalFromKeyToKey(original, transposed Key) (*NoteNames, int, error) {
	oChord, ok := ParseChord(string(original))
	if !ok {
		return nil, 0, fmt.Errorf("not a key: %v", original)
	}
	tChord, ok := ParseChord(string(transposed))
	if !ok {
		return nil, 0, fmt.Errorf("not a key: %v", original)
	}

	names, err := NoteNamesFromKey(transposed)
	if err != nil {
		return nil, 0, err
	}

	return names, int(tChord.Root) - int(oChord.Root), nil
}

// NoteNamesFromKey gets the correct NoteNames for the given key.
func NoteNamesFromKey(key Key) (*NoteNames, error) {
	for i := 0; i < len(sharpKeys); i++ {
		if sharpKeys[i] == key {
			return sharpNoteNames, nil
		}

		if i < len(flatKeys) && flatKeys[i] == key {
			return flatNoteNames, nil
		}
	}

	return nil, fmt.Errorf("invalid key name: %v", key)
}

func (n Note) String() string {
	return n.StringFromNames(sharpNoteNames)
}

// StringFromNames returns the name of a note.
func (n Note) StringFromNames(names *NoteNames) string {
	return names[n]
}

// Interval returns a new note at the specified interval.
func (n Note) Interval(interval int) Note {
	if interval < noteCount {
		interval += noteCount
	}
	return Note((int(n) + interval) % noteCount)
}

const validSuffixChars = "masd+M-245679("

// ParseChord parses some text and returns a chord.
func ParseChord(text string) (*Chord, bool) {
	if text == "" {
		return nil, false
	}

	name := text

	// root note
	rootNote, ok := parseNote(text)
	if !ok {
		return nil, false
	}

	baseNote := rootNote
	// base note
	idx := strings.Index(text, "/")
	if idx != -1 {
		baseText := text[idx+1:]
		baseNote, ok = parseNote(baseText)
		if !ok {
			return nil, false
		}

		text = text[:idx-1]
	}

	// suffix
	suffix := ""
	offset := len(rootNote.String())
	if len(text) > offset {
		suffix = text[offset:]
		found := false
		for _, r := range validSuffixChars {
			if r == rune(suffix[0]) {
				found = true
				break
			}
		}

		if !found {
			return nil, false
		}
	}

	return &Chord{
		Name:   name,
		Root:   rootNote,
		Base:   baseNote,
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
