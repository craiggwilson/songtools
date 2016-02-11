package plaintext

import (
	"fmt"
	"io"
	"io/ioutil"
)

type Token int

const (
	EofToken Token = -(iota + 1)
	NewLineToken
	TextToken
	DirectiveToken
	SectionToken
)

// Scanner produces tokens from text.
type Scanner struct {
	src string
	pos int

	isPeeked  bool
	peekPos   int
	peekToken Token
	peekText  string
	peekErr   error
}

// NewScanner creates a new Scanner from a Reader.
func NewScanner(src io.Reader) (*Scanner, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("unable to create a new scanner: %v", err)
	}

	s := &Scanner{src: string(b), pos: 0}
	return s, nil
}

// Peek returns the next token, but doesn't advance the scanner.
func (s *Scanner) Peek() (Token, string, error) {
	if !s.isPeeked {
		s.isPeeked = true
		s.peekPos = s.pos
		s.peekToken, s.peekText, s.peekErr = s.Next()
	}

	return s.peekToken, s.peekText, s.peekErr
}

// Next scans for the next token and the associated string.
func (s *Scanner) Next() (Token, string, error) {
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
		default:
			return s.scanText()
		}
	}

	return EofToken, "", nil
}

func (s *Scanner) scanDirective() (Token, string, error) {
	start := s.pos + 1
	for s.pos < len(s.src) {
		if s.src[s.pos] == '\r' || s.src[s.pos] == '\n' {
			break
		}
		s.pos++
	}

	return DirectiveToken, s.src[start:s.pos], nil
}

func (s *Scanner) scanNewLine() (Token, string, error) {
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

	return NewLineToken, "", nil
}

func (s *Scanner) scanSectionHeader() (Token, string, error) {
	start := s.pos + 1
	for s.pos < len(s.src) {
		s.pos++
		if s.src[s.pos] == ']' {
			s.pos++
			return SectionToken, s.src[start : s.pos-1], nil
		}
	}
	if s.pos == len(s.src) {
		return EofToken, "", fmt.Errorf("Expected ']', but found Eof")
	}
	return EofToken, "", nil
}

func (s *Scanner) scanText() (Token, string, error) {
	start := s.pos
	for s.pos < len(s.src) {
		if s.src[s.pos] == '\r' || s.src[s.pos] == '\n' {
			break
		}
		s.pos++
	}

	return TextToken, s.src[start:s.pos], nil
}
