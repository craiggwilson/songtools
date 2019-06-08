package chord

// Quality is the quality of a chord.
type Quality byte

// Quality constants
const (
	Major Quality = iota
	Minor
	Augmented
	Diminished
	HalfDiminished
)
