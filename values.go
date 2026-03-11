package main

import (
	"fmt"
	. "kl/tokens"
)

type Value struct {
	Type  string
	Int   int
	Float float64
	Str   string
}

type Param struct {
	Name string
	Type string
}

func (p *Parser) checkTypeAssignment(varName string, varType string, tok Token) {
	switch varType {
	case "float":
		if tok.Type != TOKEN_FLOAT {
			switch tok.Type {
			case TOKEN_STRING:
				p.langError(tok, fmt.Sprintf("cannot use %q (type string) as type %s in assignment to '%s'", tok.Val, varType, varName))
			case TOKEN_INT:
				p.langError(tok, fmt.Sprintf("cannot use %q (type int) as type %s in assignment to '%s'", tok.Val, varType, varName))
			}
		}
	case "int":
		if tok.Type != TOKEN_INT {
			switch tok.Type {
			case TOKEN_STRING:
				p.langError(tok, fmt.Sprintf("cannot use %q (type string) as type %s in assignment to '%s'", tok.Val, varType, varName))
			case TOKEN_FLOAT:
				p.langError(tok, fmt.Sprintf("cannot use %q (type float) as type %s in assignment to '%s'", tok.Val, varType, varName))
			}
		}
	case "string":
		if tok.Type != TOKEN_STRING {
			switch tok.Type {
			case TOKEN_INT:
				p.langError(tok, fmt.Sprintf("cannot use %q (type int) as type %s in assignment to '%s'", tok.Val, varType, varName))
			case TOKEN_FLOAT:
				p.langError(tok, fmt.Sprintf("cannot use %q (type float) as type %s in assignment to '%s'", tok.Val, varType, varName))
			}
		}
	}
}

func (p *Parser) GetVar(varname string) Value {
	if val, ok := p.Variables[varname]; ok {
		return val
	}
	panic(fmt.Sprintf("undeclared var: %v", varname))
}

func (p *Parser) setVar(varname string, val Value) {
	p.Variables[varname] = val
}

func (p *Parser) isVar(varName string) bool {
	if _, ok := p.Variables[varName]; ok {
		return true
	}
	return false
}
