package note

type EnglishLanguage struct{}

var englishNames = []string{"Ab", "A", "Bb", "B", "C", "Db", "D", "Eb", "E", "F", "Gb", "G"}

func (EnglishLanguage) Name(n Note) (string, error) {
	if n >= 0 && n < 12 {
		return englishNames[n], nil
	}
	return "", ErrInvalidNote
}

func (EnglishLanguage) Parse(s string) (Note, uint8, error) {
	if len(s) > 1 {
		switch s[:2] {
		case "G#", "Ab":
			return 0, 2, nil
		case "A#", "Bb":
			return 2, 2, nil
		case "C#", "Db":
			return 5, 2, nil
		case "D#", "Eb":
			return 7, 2, nil
		case "F#", "Gb":
			return 10, 2, nil
		}
	} else if len(s) > 1 {
		switch s[:1] {
		case "A":
			return 1, 1, nil
		case "B":
			return 3, 1, nil
		case "C":
			return 4, 1, nil
		case "D":
			return 6, 1, nil
		case "E":
			return 8, 1, nil
		case "F":
			return 9, 1, nil
		case "G":
			return 11, 1, nil
		}
	}

	return 0, 0, ErrInvalidNote
}
