package main

import (
	"fmt"
	. "kl/tokens"
)

type Function struct {
	Name       string
	Params     []Param
	Body       []Token
	ReturnType string
}

func (p *Parser) ParseFunc() {
	p.expect(TOKEN_FUNC_KEYWORD)
	var returnType string

	funcName := p.expect(TOKEN_IDENT).Val
	p.expect(TOKEN_LPAREN)
	params := make([]Param, 0)
	for p.CurrentToken.Type != TOKEN_RPAREN {
		typeParam := p.CurrentToken.Val
		p.advanceParser()
		paramName := p.expect(TOKEN_IDENT).Val
		params = append(params, Param{Name: paramName, Type: typeParam})
		if p.CurrentToken.Type != TOKEN_RPAREN {
			p.expect(TOKEN_COMMA)
		}

	}
	p.expect(TOKEN_RPAREN)

	if p.CurrentToken.Type != TOKEN_LBRACE {
		returnType = p.CurrentToken.Val
		p.advanceParser()
	}
	p.expect(TOKEN_LBRACE)

	var body []Token
	depth := 1
	for depth > 0 {
		tok := p.CurrentToken
		if tok.Type == TOKEN_LBRACE {
			depth++
		}
		if tok.Type == TOKEN_RBRACE {
			depth--
		}
		if depth > 0 {
			body = append(body, tok)
		}
		p.advanceParser()
	}
	p.Functions[funcName] = Function{Name: funcName, Params: params, Body: body, ReturnType: returnType}

}

func (p *Parser) isFuncCall() bool {
	name := p.CurrentToken.Val
	_, isFn := p.Functions[name]
	return isFn
}

func (p *Parser) callFunc() Value {
	name := p.expect(TOKEN_IDENT).Val
	if _, isFn := p.Functions[name]; !isFn {
		panic(fmt.Sprintf("Unknow functions: %s", name))
	}
	result := p.callFuncExpr(name)
	p.expect(TOKEN_SEMICOLON)
	return Value{Type: "float", Float: result}

}

func (p *Parser) callFuncExpr(name string) float64 {
	fn, ok := p.Functions[name]
	if !ok {
		p.langError(p.CurrentToken, fmt.Sprintf("undefined: %s", name))
	}
	p.expect(TOKEN_LPAREN)

	var args []Value
	for p.CurrentToken.Type != TOKEN_RPAREN {
		args = append(args, p.parseArgValue())
		if p.CurrentToken.Type == TOKEN_COMMA {
			p.expect(TOKEN_COMMA)
		}
	}
	p.expect(TOKEN_RPAREN)

	if len(args) != len(fn.Params) {
		p.langError(p.CurrentToken, fmt.Sprintf("too many arguments in call to %s: have: %d, want: %d", name, len(args), len(fn.Params)))
	}

	for i, param := range fn.Params {
		arg := args[i]
		if arg.Type != param.Type {
			p.langError(p.CurrentToken, fmt.Sprintf("cannot use %s(type %s) as type %s in argument '%s' to %s", arg.Type, arg.Type, param.Type, param.Name, name))
		}
	}

	bodyLexer := &TokenSliceLexer{Tokens: fn.Body, Pos: 0}
	subParser := &Parser{
		Lexer:     bodyLexer,
		Variables: make(map[string]Value),
		Functions: p.Functions,
	}
	for i, param := range fn.Params {
		subParser.Variables[param.Name] = args[i]
	}
	subParser.CurrentToken = bodyLexer.NextToken()

	for subParser.CurrentToken.Type != TOKEN_EOF {
		if subParser.CurrentToken.Type == TOKEN_RETURN_KEYWORD {
			subParser.advanceParser()
			result := subParser.parseExpression()

			if fn.ReturnType == "int" {
				return float64(int(result))
			}

			if fn.ReturnType == "float" {
				return result
			}

			if fn.ReturnType == "string" {
				p.langError(p.CurrentToken, "cannot return numeric value from string function")
			}
			return result
		}
		subParser.Parse()
	}
	if fn.ReturnType != "" {
		p.langError(p.CurrentToken, fmt.Sprintf("missing return statement in function %s", fn.Name))
	}
	return 0
}
