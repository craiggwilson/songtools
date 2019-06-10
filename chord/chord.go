package chord

import (
	"github.com/songtools/songtools/note"
)

func Parse(s string) (Chord, error) {
	return parse(s)
}

type Chord struct {
	Name      string
	Root      note.Note
	Intervals []uint8
}

func parse(s string) (Chord, error) {

}
