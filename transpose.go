package songtools

// TransposeSongSet transposes a SongSet.
func TransposeSongSet(ss *SongSet, interval int, names *NoteNames) (*SongSet, error) {
	newSongs := []*Song{}
	for _, s := range ss.Songs {
		newSong, err := TransposeSong(s, interval, names)
		if err != nil {
			return nil, err
		}
		newSongs = append(newSongs, newSong)
	}

	return &SongSet{newSongs}, nil
}

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

	return &Song{newNodes}, nil
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
