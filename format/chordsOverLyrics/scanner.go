package chordsOverLyrics

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type token int

const (
	eofToken token = -(iota + 1)
	newLineToken
	textToken
	commentToken
	directiveToken
	sectionToken
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
	case sectionToken:
		return "<section>"
	default:
		return "<unknown>"
	}
}

// Scanner produces tokens from text.
type scanner struct {
	src string
	pos int

	isPeeked  bool
	peekPos   int
	peekToken token
	peekText  string
	peekErr   error
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
func (s *scanner) peek() (token, string, error) {
	if !s.isPeeked {
		s.isPeeked = true
		s.peekPos = s.pos
		s.peekToken, s.peekText, s.peekErr = s.next()
	}

	return s.peekToken, s.peekText, s.peekErr
}

// Next scans for the next token and the associated string.
func (s *scanner) next() (token, string, error) {
	if s.isPeeked {
		s.isPeeked = false
		return s.peekToken, s.peekText, s.peekErr
	}

	if s.pos < len(s.src) {
		switch s.src[s.pos] {
		case '\r':
			return s.scanNewLine()
		case '\n':
			return s.scanNewLine()
		case '#':
			return s.scanDirective()
		case '[':
			return s.scanSectionHeader()
		case '/':
			if s.pos+1 < len(s.src) && s.src[s.pos+1] == '/' {
				return s.scanComment()
			}
			return s.scanText()
		default:
			return s.scanText()
		}
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
		if s.src[s.pos] == '\r' || s.src[s.pos] == '\n' {
			break
		}
		s.pos++
	}

	return directiveToken, s.src[start:s.pos], nil
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

func (s *scanner) scanSectionHeader() (token, string, error) {
	start := s.pos + 1
	for s.pos < len(s.src) {
		s.pos++
		if s.src[s.pos] == ']' {
			s.pos++
			return sectionToken, s.src[start : s.pos-1], nil
		}
	}
	if s.pos == len(s.src) {
		return eofToken, "", fmt.Errorf("Expected ']', but found Eof")
	}
	return eofToken, "", nil
}

func (s *scanner) scanText() (token, string, error) {
	start := s.pos
	for s.pos < len(s.src) {
		if s.src[s.pos] == '\r' || s.src[s.pos] == '\n' {
			if strings.TrimSpace(s.src[start:s.pos]) == "" {
				return s.scanNewLine()
			}
			break
		}
		s.pos++
	}

	return textToken, s.src[start:s.pos], nil
}
