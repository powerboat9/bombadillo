package cmdparse

import (
	"bufio"
	"strings"
	"io"
	"bytes"
)



//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Token struct {
	kind	tok
	val	string
}

type scanner struct {
	r	*bufio.Reader
}

type tok int


//------------------------------------------------\\
// + + +         V A R I A B L E S           + + + \\
//--------------------------------------------------\\

var eof rune = rune(0)
const (
	Word tok = iota
	Action
	Value
	End
	Whitespace

	number
	letter
	ws
	illegal
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

func (s *scanner) scanText() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			s.unread()
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	capInput := strings.ToUpper(buf.String())
	switch capInput {
		case "DELETE", "ADD", "WRITE", "SET", "RECALL", "R", "SEARCH",
			"W", "A", "D", "S", "Q", "QUIT", "B", "BOOKMARKS", "H", "HOME":
			return Token{Action, capInput}
	}

	return Token{Word, buf.String()}
}

func (s *scanner) scanWhitespace() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			s.unread()
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			_,_ = buf.WriteRune(ch)
		}
	}

	return Token{Whitespace, buf.String()}

}

func (s *scanner) scanNumber() Token {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			_,_ = buf.WriteRune(ch)
		}
	}
	
	return Token{Value, buf.String()}
}

func (s *scanner) unread() { _ = s.r.UnreadRune() }

func (s *scanner) scan() Token {
	char := s.read()

	if isWhitespace(char) {
		s.unread()
		return s.scanWhitespace()
	} else if isDigit(char) {
		s.unread()
		return s.scanNumber()
	} else if isLetter(char) {
		s.unread()
		return s.scanText()
	}	

	if char == eof {
		return Token{End, ""}
	}

	return Token{illegal, string(char)}
}


//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func NewScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return ch >= '!' && ch <= '~'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isEOF(ch rune) bool {
	return ch == rune(0)
}

