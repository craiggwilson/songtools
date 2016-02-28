package songtools

// TransposeSong transposes a Song.
func TransposeSong(s *Song, interval int, names *NoteNames) (*Song, error) {
	newNodes := []SongNode{}
	for _, n := range s.Nodes {
		switch typedN := n.(type) {
		case *Section:
			newSection, err := TransposeSection(typedN, interval, names)
			if err != nil {
				return nil, err
			}

			newNodes = append(newNodes, newSection)
		default:
			newNodes = append(newNodes, n)
		}
	}

	key := s.Key
	if key != "" {
		if c, ok := ParseChord(string(key)); ok {
			key = Key(c.Interval(interval, names).Name)
		}
	}

	return &Song{
		Title:     s.Title,
		Subtitles: s.Subtitles,
		Authors:   s.Authors,
		Key:       key,
		Nodes:     newNodes,
	}, nil
}

// TransposeSection transposes a Section.
func TransposeSection(s *Section, interval int, names *NoteNames) (*Section, error) {
	newNodes := []SectionNode{}
	for _, n := range s.Nodes {
		switch typedN := n.(type) {
		case *Line:
			newLine, err := TransposeLine(typedN, interval, names)
			if err != nil {
				return nil, err
			}

			newNodes = append(newNodes, newLine)
		default:
			newNodes = append(newNodes, n)
		}
	}

	return &Section{s.Kind, newNodes}, nil
}

// TransposeLine transposes a Line.
func TransposeLine(l *Line, interval int, names *NoteNames) (*Line, error) {
	if l.Chords == nil {
		return l, nil
	}

	newChords := []*Chord{}
	for _, c := range l.Chords {
		newChord := c.Interval(interval, names)
		newChords = append(newChords, newChord)
	}

	return &Line{l.Text, newChords, l.ChordPositions}, nil
}
