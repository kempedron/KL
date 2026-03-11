package main

import (
	"fmt"
	. "kl/tokens"
)

type ILexer interface {
	NextToken() Token
	GetNextToken() Token
}

type Parser struct {
	Lexer        ILexer
	CurrentToken Token
	Variables    map[string]Value
	Functions    map[string]Function
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{Lexer: lexer, Variables: make(map[string]Value, 10), Functions: make(map[string]Function)}
	p.CurrentToken = p.Lexer.NextToken()
	return p
}

func (p *Parser) expect(t TOKEN_TYPE) Token {
	if p.CurrentToken.Type != t {
		p.langError(p.CurrentToken, fmt.Sprintf("expected: '%v', got: '%v'", t.String(), p.CurrentToken.Type.String()))
	}
	token := p.CurrentToken
	p.advanceParser()

	return token
}

func (parser *Parser) advanceParser() {
	if parser.CurrentToken.Type == TOKEN_EOF {
		return
	}
	parser.CurrentToken = parser.Lexer.NextToken()
}
