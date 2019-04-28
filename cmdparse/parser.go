package cmdparse

import (
	"fmt"
	"io"
)

//------------------------------------------------\\
// + + +             T Y P E S               + + + \\
//--------------------------------------------------\\

type Parser struct {
	s      *scanner
	buffer struct {
		token Token
		size  int
	}
}

type Command struct {
	Action string
	Target string
	Value  []string
	Type   Comtype
}

type Comtype int

const (
	GOURL Comtype = iota
	GOLINK
	SIMPLE
	DOLINK
	DOLINKAS
	DOAS
)

//------------------------------------------------\\
// + + +           R E C E I V E R S         + + + \\
//--------------------------------------------------\\

func (p *Parser) scan() (current Token) {
	if p.buffer.size != 0 {
		p.buffer.size = 0
		return p.buffer.token
	}

	current = p.s.scan()
	for {
		if current.kind != Whitespace {
			break
		}
		current = p.s.scan()
	}

	p.buffer.token = current
	return
}

func (p *Parser) unscan() { p.buffer.size = 1 }

func (p *Parser) parseNonAction() (*Command, error) {
	p.unscan()
	t := p.scan()
	cm := &Command{}

	if t.kind == Value {
		cm.Target = t.val
		cm.Type = GOLINK
	} else if t.kind == Word {
		cm.Target = t.val
		cm.Type = GOURL
	} else {
		return nil, fmt.Errorf("Found %q, expected action, url, or link number", t.val)
	}

	if u := p.scan(); u.kind != End {
		return nil, fmt.Errorf("Found %q, expected EOF", u.val)
	}
	return cm, nil
}

func (p *Parser) parseAction() (*Command, error) {
	p.unscan()
	t := p.scan()
	cm := &Command{}
	cm.Action = t.val
	t = p.scan()
	switch t.kind {
	case End:
		cm.Type = SIMPLE
		return cm, nil
	case Value:
		cm.Target = t.val
		cm.Type = DOLINK
	case Word:
		cm.Value = append(cm.Value, t.val)
		cm.Type = DOAS
	case Action, Whitespace:
		return nil, fmt.Errorf("Found %q (%d), expected value", t.val, t.kind)
	}
	t = p.scan()
	if t.kind == End {
		return cm, nil
	} else {
		if cm.Type == DOLINK {
			cm.Type = DOLINKAS
		} else {
			cm.Type = DOAS
		}
		cm.Value = append(cm.Value, t.val)

		for {
			token := p.scan()
			if token.kind == End {
				break
			} else if token.kind == Whitespace {
				continue
			}
			cm.Value = append(cm.Value, token.val)
		}
	}
	return cm, nil
}

func (p *Parser) Parse() (*Command, error) {
	if t := p.scan(); t.kind != Action {
		return p.parseNonAction()
	} else {
		return p.parseAction()
	}
}

//------------------------------------------------\\
// + + +          F U N C T I O N S          + + + \\
//--------------------------------------------------\\

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}
