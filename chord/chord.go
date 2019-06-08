package chord

import (
	"github.com/songtools/songtools/note"
)

func Parse(s string) (Chord, bool) {
	return Chord{}, false
}

type Chord struct {
	Name      string
	Root      note.Note
	Quality   Quality
	Intervals uint8
}
