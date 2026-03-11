package main

import "fmt"

type LangError struct {
	Line    int
	Col     int
	Message string
}

func (p *Parser) langError(tok Token, msg string) {
	panic(fmt.Sprintf("%d:%d: %s", tok.Line+1, tok.Col+1, msg))
}
