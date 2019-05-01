package config

import (
	"fmt"
	"io"
	"strings"
	"tildegit.org/sloum/bombadillo/gopher"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Parser struct {
	s      *scanner
	row    int
	buffer struct {
		token Token
		size  int
	}
}

type Config struct {
	Bookmarks gopher.Bookmarks
	Colors    []KeyValue
	Settings  []KeyValue
}

type KeyValue struct {
	Key   string
	Value string
}

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (p *Parser) scan() (current Token) {
	if p.buffer.size != 0 {
		p.buffer.size = 0
		return p.buffer.token
	}

	current = p.s.scan()
	p.buffer.token = current
	return
}

func (p *Parser) parseKeyValue() (KeyValue, error) {
	kv := KeyValue{}
	t1 := p.scan()
	kv.Key = strings.TrimSpace(t1.val)

	if t := p.scan(); t.kind == TOK_VALUE {
		kv.Value = strings.TrimSpace(t.val)
	} else {
		return kv, fmt.Errorf("Got non-value expected VALUE on row %d", p.row)
	}

	if t := p.scan(); t.kind != TOK_NEWLINE {
		return kv, fmt.Errorf("Expected NEWLINE, got %q on row %d", t.kind, p.row)
	}

	return kv, nil
}

func (p *Parser) unscan() { p.buffer.size = 1 }

func (p *Parser) Parse() (Config, error) {
	p.row = 1
	section := ""
	c := Config{}

	for {
		if t := p.scan(); t.kind == TOK_NEWLINE {
			p.row++
		} else if t.kind == TOK_SECTION {
			section = strings.ToUpper(t.val)
		} else if t.kind == TOK_EOF {
			break
		} else if t.kind == TOK_KEY {
			p.unscan()
			keyval, err := p.parseKeyValue()
			if err != nil {
				return Config{}, err
			}
			switch section {
			case "BOOKMARKS":
				err := c.Bookmarks.Add([]string{keyval.Value, keyval.Key})
				if err != nil {
					return c, err
				}
			case "COLORS":
				c.Colors = append(c.Colors, keyval)
			case "SETTINGS":
				c.Settings = append(c.Settings, keyval)
			}
		} else if t.kind == TOK_ERROR {
			return Config{}, fmt.Errorf("Error on row %d: %s", p.row, t.val)
		}
	}

	return c, nil
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}
