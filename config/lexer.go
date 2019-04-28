package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Token struct {
	kind TokenType
	val  string
}

type scanner struct {
	r *bufio.Reader
}

type TokenType int

//------------------------------------------------\\
// + + +         V A R I A B L E S           + + + \\
//--------------------------------------------------\\

var eof rune = rune(0)
var l_brace rune = '['
var r_brace rune = ']'
var newline rune = '\n'
var equal rune = '='

const (
	TOK_SECTION TokenType = iota
	TOK_KEY
	TOK_VALUE

	TOK_EOF
	TOK_NEWLINE
	TOK_ERROR
	TOK_WHITESPACE
)

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *scanner) skipWhitespace() {
	for {
		if ch := s.read(); ch == eof {
			s.unread()
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}
}

func (s *scanner) skipToEndOfLine() {
	for {
		if ch := s.read(); ch == eof {
			s.unread()
			break
		} else if ch == newline {
			s.unread()
			break
		}
	}
}

func (s *scanner) scanSection() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			s.unread()
			return Token{TOK_ERROR, "Reached end of feed without closing section"}
		} else if ch == r_brace {
			break
		} else if ch == newline {
			s.unread()
			return Token{TOK_ERROR, "No closing brace for section before newline"}
		} else if ch == equal {
			return Token{TOK_ERROR, "Illegal character"}
		} else if ch == l_brace {
			s.skipToEndOfLine()
			return Token{TOK_ERROR, "Second left brace encountered before closing right brace in section"}
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Token{TOK_SECTION, buf.String()}
}

func (s *scanner) scanKey() Token {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			s.unread()
			return Token{TOK_ERROR, "Reached end of feed without assigning value to key"}
		} else if ch == equal {
			s.unread()
			break
		} else if ch == newline {
			s.unread()
			return Token{TOK_ERROR, "No value assigned to key"}
		} else if ch == r_brace || ch == l_brace {
			s.skipToEndOfLine()
			return Token{TOK_ERROR, "Illegal brace character in key"}
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Token{TOK_KEY, buf.String()}
}

func (s *scanner) scanValue() Token {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			s.unread()
			break
		} else if ch == equal {
			_, _ = buf.WriteRune(ch)
		} else if ch == newline {
			s.unread()
			break
		} else if ch == r_brace || ch == l_brace {
			s.skipToEndOfLine()
			return Token{TOK_ERROR, "Illegal brace character in key"}
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Token{TOK_VALUE, buf.String()}
}

func (s *scanner) unread() { _ = s.r.UnreadRune() }

func (s *scanner) scan() Token {
	char := s.read()

	if isWhitespace(char) {
		s.skipWhitespace()
	}

	if char == l_brace {
		return s.scanSection()
	} else if isText(char) {
		s.unread()
		return s.scanKey()
	} else if char == equal {
		return s.scanValue()
	} else if char == newline {
		return Token{TOK_NEWLINE, "New line"}
	}

	if char == eof {
		return Token{TOK_EOF, "Reached end of feed"}
	}

	return Token{TOK_ERROR, fmt.Sprintf("Error on character %q", char)}
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func NewScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isText(ch rune) bool {
	return ch >= '!' && ch <= '~' && ch != equal && ch != l_brace && ch != r_brace
}

func isEOF(ch rune) bool {
	return ch == eof
}
