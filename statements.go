package main

import (
	"fmt"
	. "kl/tokens"
	"strconv"
)

func (p *Parser) Parse() {
	switch p.CurrentToken.Type {
	case TOKEN_INT_KEYWORD:
		p.parseInt()
	case TOKEN_PRINT:
		p.printToken()
	case TOKEN_SEMICOLON:
		p.advanceParser()
	case TOKEN_STRING_KEYWORD:
		p.parseStr()
	case TOKEN_IDENT:
		if p.isFuncCall() {
			p.callFunc()
		} else {
			p.parseReAssignVar()
		}
	case TOKEN_FLOAT_KEYWORD:
		p.parseFloat()
	case TOKEN_INPUT_KEYWORD:
		p.parseInput()
	case TOKEN_PRINTF:
		p.ParsePrintf()
	case TOKEN_FUNC_KEYWORD:
		p.ParseFunc()
	case TOKEN_RPAREN:
		p.langError(p.CurrentToken, "unexpected ')'")
	case TOKEN_IF_KEYWORD:
		p.parseIf()
	default:
		panic(fmt.Sprintf("unexpected token: %v", p.CurrentToken.Type))
	}

}

func (p *Parser) parseInput() {
	p.expect(TOKEN_INPUT_KEYWORD)
	p.expect(TOKEN_LPAREN)

	ident := p.expect(TOKEN_IDENT).Val
	p.expect(TOKEN_COMMA)
	text := p.expect(TOKEN_STRING).Val
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)

	if !p.isVar(ident) {
		p.langError(p.CurrentToken, fmt.Sprintf("undefined: %s", ident))
	}

	fmt.Print(text)
	var input string
	fmt.Scan(&input)
	existing := p.Variables[ident]
	switch existing.Type {
	case "int":
		val, err := strconv.Atoi(input)
		if err != nil {
			p.langError(p.CurrentToken, fmt.Sprintf("cannot convert %q to int for '%s'", input, ident))
		}
		p.Variables[ident] = Value{Type: "int", Int: val}

	case "float":
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			p.langError(p.CurrentToken, fmt.Sprintf("cannot convert %q to float for '%s'", input, ident))
		}
		p.Variables[ident] = Value{Type: "float", Float: val}
	default:
		p.Variables[ident] = Value{Type: "string", Str: input}
	}
}

func (p *Parser) parseIf() {
	p.expect(TOKEN_IF_KEYWORD)
	p.expect(TOKEN_LPAREN)
	condition := p.parseExpression()
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_LBRACE)

	ifBody := p.collectBody()

	type elifBranch struct {
		condition float64
		body      []Token
	}

	var elifBranches []elifBranch

	for p.CurrentToken.Type == TOKEN_ELIF_KEYWORD {
		p.expect(TOKEN_ELIF_KEYWORD)
		p.expect(TOKEN_LPAREN)
		elifCondition := p.parseExpression()
		p.expect(TOKEN_RPAREN)
		p.expect(TOKEN_LBRACE)
		elifBody := p.collectBody()
		elifBranches = append(elifBranches, elifBranch{elifCondition, elifBody})
	}

	var elseBody []Token
	if p.CurrentToken.Type == TOKEN_ELSE_KEYWORD {
		p.expect(TOKEN_ELSE_KEYWORD)
		p.expect(TOKEN_LBRACE)
		elseBody = p.collectBody()
	}
	if condition != 0 {
		p.execBody(ifBody)
		return
	}

	for _, branch := range elifBranches {
		if branch.condition != 0 {
			p.execBody(branch.body)
			return
		}
	}
	if elseBody != nil {
		p.execBody(elseBody)
	}

}

func (p *Parser) collectBody() []Token {
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
	return body
}

func (p *Parser) execBody(body []Token) {
	bodyLexer := &TokenSliceLexer{Tokens: body, Pos: 0}
	subParser := &Parser{
		Lexer:     bodyLexer,
		Variables: p.Variables,
		Functions: p.Functions,
	}
	subParser.CurrentToken = bodyLexer.NextToken()
	for subParser.CurrentToken.Type != TOKEN_EOF && subParser.CurrentToken.Type != TOKEN_RBRACE {
		subParser.Parse()
	}
	for k, v := range subParser.Variables {
		p.Variables[k] = v
	}
	for k, v := range subParser.Functions {
		p.Functions[k] = v
	}
}

func (p *Parser) parseInt() {
	p.expect(TOKEN_INT_KEYWORD)
	ident := p.expect(TOKEN_IDENT)

	if p.CurrentToken.Type == TOKEN_SEMICOLON {

		p.setVar(ident.Val, Value{Type: "int", Int: 0})
		p.expect(TOKEN_SEMICOLON)
		return
	}
	p.expect(TOKEN_ASSIGN)
	p.checkTypeAssignment(ident.Val, "int", p.CurrentToken)
	result := p.parseExpression()

	p.expect(TOKEN_SEMICOLON)

	p.setVar(ident.Val, Value{Type: "int", Int: int(result)})
}

func (p *Parser) parseFloat() {
	p.expect(TOKEN_FLOAT_KEYWORD)
	ident := p.expect(TOKEN_IDENT)
	if p.CurrentToken.Type == TOKEN_SEMICOLON {
		p.setVar(ident.Val, Value{Type: "float", Float: 0})
		p.expect(TOKEN_SEMICOLON)
		return
	}
	p.expect(TOKEN_ASSIGN)
	p.checkTypeAssignment(ident.Val, "float", p.CurrentToken)
	result := p.parseExpression()
	p.expect(TOKEN_SEMICOLON)
	p.setVar(ident.Val, Value{Type: "float", Float: float64(result)})
}

func (p *Parser) parseStr() {
	p.expect(TOKEN_STRING_KEYWORD)
	ident := p.expect(TOKEN_IDENT)
	if p.CurrentToken.Type == TOKEN_SEMICOLON {
		p.setVar(ident.Val, Value{Type: "string", Str: ""})
		p.expect(TOKEN_SEMICOLON)
		return
	}
	p.expect(TOKEN_ASSIGN)

	if p.CurrentToken.Type == TOKEN_STRING_KEYWORD {
		val := p.parseStringCast()
		p.expect(TOKEN_SEMICOLON)
		p.setVar(ident.Val, val)
		return
	}

	p.checkTypeAssignment(ident.Val, "string", p.CurrentToken)

	t := p.expect(TOKEN_STRING)
	p.expect(TOKEN_SEMICOLON)
	p.setVar(ident.Val, Value{Type: "string", Str: t.Val})

}

func (p *Parser) parseReAssignVar() {
	varName := p.expect(TOKEN_IDENT).Val
	p.expect(TOKEN_ASSIGN)

	if p.CurrentToken.Type == TOKEN_STRING_KEYWORD {
		val := p.parseStringCast()
		p.expect(TOKEN_SEMICOLON)
		p.setVar(varName, val)
		return
	}
	if existing, ok := p.Variables[varName]; ok {
		p.checkTypeAssignment(varName, existing.Type, p.CurrentToken)
	}
	var varVal Value
	switch p.CurrentToken.Type {
	case TOKEN_INT:
		integer, err := strconv.Atoi(p.expect(TOKEN_INT).Val)
		if err != nil {
			panic(fmt.Sprintf("failed atoi: %s", err))
		}
		varVal = Value{Type: "int", Int: integer}
	case TOKEN_STRING:
		varVal = Value{Type: "string", Str: p.expect(TOKEN_STRING).Val}
	case TOKEN_FLOAT:
		f, err := strconv.ParseFloat(p.expect(TOKEN_FLOAT).Val, 64)
		if err != nil {
			panic(fmt.Sprintf("failed atoi: %s", err))
		}
		varVal = Value{Type: "float", Float: f}
	default:
		result := p.parseExpression()
		varVal = Value{Type: "float", Float: result}
	}

	p.setVar(varName, varVal)
	p.expect(TOKEN_SEMICOLON)
}
