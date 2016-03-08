package chordpro

import (
	"fmt"
	"io"
	"io/ioutil"
)

type token int

const (
	eofToken token = -(iota + 1)
	newLineToken
	textToken
	commentToken
	directiveToken
	chordToken
)

func (t token) String() string {
	switch t {
	case eofToken:
		return "<eof>"
	case newLineToken:
		return "<newline>"
	case textToken:
		return "<text>"
	case commentToken:
		return "<comment>"
	case directiveToken:
		return "<directive>"
	case chordToken:
		return "<chordToken>"
	default:
		return "<unknown>"
	}
}

// Scanner produces tokens from text.
type scanner struct {
	src string
	pos int

	peeks []peek
}

type peek struct {
	pos   int
	token token
	text  string
	err   error
}

// NewScanner creates a new Scanner from a Reader.
func newScanner(src io.Reader) (*scanner, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new scanner: %v", err)
	}

	s := &scanner{src: string(b), pos: 0}
	return s, nil
}

// Peek returns the next token, but doesn't advance the scanner.
func (s *scanner) la(count int) (token, string, error) {
	for len(s.peeks) <= count {
		nextPos := s.pos
		nextToken, nextText, nextErr := s._internalNext()
		s.peeks = append(s.peeks, peek{nextPos, nextToken, nextText, nextErr})
	}

	peek := s.peeks[count]

	return peek.token, peek.text, peek.err
}

// Next scans for the next token and the associated string.
func (s *scanner) next() (token, string, error) {
	if len(s.peeks) > 0 {
		peek := s.peeks[0]
		s.peeks = s.peeks[1:]
		return peek.token, peek.text, peek.err
	}

	return s._internalNext()
}

func (s *scanner) _internalNext() (token, string, error) {
	if s.pos < len(s.src) {
		switch s.src[s.pos] {
		case '\r':
			return s.scanNewLine()
		case '\n':
			return s.scanNewLine()
		case '#':
			return s.scanComment()
		case '[':
			return s.scanChord()
		case '{':
			return s.scanDirective()
		default:
			return s.scanText()
		}
	}

	return eofToken, "", nil
}

func (s *scanner) scanChord() (token, string, error) {
	start := s.pos + 1
	for s.pos < len(s.src) {
		s.pos++
		if s.src[s.pos] == ']' {
			s.pos++
			return chordToken, s.src[start : s.pos-1], nil
		}
	}
	if s.pos == len(s.src) {
		return eofToken, "", fmt.Errorf("Expected ']', but found Eof")
	}
	return eofToken, "", nil
}

func (s *scanner) scanComment() (token, string, error) {
	start := s.pos + 1
	for s.pos < len(s.src) {
		if s.src[s.pos] == '\r' || s.src[s.pos] == '\n' {
			break
		}
		s.pos++
	}

	return commentToken, s.src[start:s.pos], nil
}

func (s *scanner) scanDirective() (token, string, error) {
	start := s.pos + 1
	for s.pos < len(s.src) {
		s.pos++
		if s.src[s.pos] == '}' {
			s.pos++
			return directiveToken, s.src[start : s.pos-1], nil
		}
	}
	if s.pos == len(s.src) {
		return eofToken, "", fmt.Errorf("Expected '}', but found Eof")
	}
	return eofToken, "", nil
}

func (s *scanner) scanNewLine() (token, string, error) {
	for s.pos < len(s.src) {
		if s.src[s.pos] == '\n' {
			s.pos++
			break
		}
		if s.src[s.pos] == '\r' {
			s.pos++
		} else {
			break
		}
	}

	return newLineToken, "", nil
}

func (s *scanner) scanText() (token, string, error) {
	start := s.pos
	for s.pos < len(s.src) {
		if s.src[s.pos] == '\r' || s.src[s.pos] == '\n' || s.src[s.pos] == '[' {
			break
		}
		s.pos++
	}

	return textToken, s.src[start:s.pos], nil
}
