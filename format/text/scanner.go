package text

import (
	"io"
	"strings"
)

//go:generate goyacc -o parser.go parser.y

func NewScannerString(s string) *Scanner {
	return &Scanner{
		input: strings.NewReader(s),
	}
}

type Scanner struct {
	input    io.RuneReader
	pos      uint
	lastRune rune
}

func (s *Scanner) Lex(lval *yySymType) int {
	return 0
}

func (s *Scanner) Error(e string) {
}

func (s *Scanner) scan() (int, []rune) {
	if s.lastRune == "" {
		s.next()
	}

	switch r := s.lastRune; {
	case isWhitespace(r):
		break
	}
}

func (s *Scanner) next() {
	r, _, err := s.input.ReadRune()
	if err != nil {

	}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}
