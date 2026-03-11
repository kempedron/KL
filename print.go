package main

import (
	"fmt"
	. "kl/tokens"
)

func (p *Parser) ParsePrintf() {
	p.expect(TOKEN_PRINTF)
	p.expect(TOKEN_LPAREN)
	var args []any
	strToFormatting := p.expect(TOKEN_STRING).Val
	for p.CurrentToken.Type != TOKEN_RPAREN {
		p.expect(TOKEN_COMMA)

		if p.CurrentToken.Type == TOKEN_STRING {
			args = append(args, p.expect(TOKEN_STRING).Val)
			continue
		}
		if p.CurrentToken.Type == TOKEN_IDENT {
			if val, ok := p.Variables[p.CurrentToken.Val]; ok && val.Type == "string" {
				args = append(args, val.Str)
				p.advanceParser()
				continue
			}
		}
		result := p.parseExpression()
		if result == float64(int(result)) {
			args = append(args, int(result))
		} else {
			args = append(args, result)
		}
	}
	fmt.Printf(strToFormatting, args...)
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)

}

func (p *Parser) printToken() {
	p.expect(TOKEN_PRINT)
	p.expect(TOKEN_LPAREN)
	if p.CurrentToken.Type == TOKEN_QUOTATION {
		p.printString()
	} else {
		p.printVar()
	}

}

func (p *Parser) printVar() {
	if p.CurrentToken.Type == TOKEN_RPAREN {
		p.expect(TOKEN_RPAREN)
		p.expect(TOKEN_SEMICOLON)
		return

	}

	for {

		p.printArgument()
		if p.CurrentToken.Type == TOKEN_COMMA {
			p.expect(TOKEN_COMMA)
			continue
		}
		break
	}

	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)

	fmt.Println()
}

func (p *Parser) printString() {
	p.expect(TOKEN_QUOTATION)

	p.printArgument()

	p.expect(TOKEN_QUOTATION)
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)
}

func (p *Parser) printArgument() {
	switch p.CurrentToken.Type {
	case TOKEN_STRING:
		token := p.expect(TOKEN_STRING)
		fmt.Print(token.Val)
	case TOKEN_IDENT:
		if _, isFn := p.Functions[p.CurrentToken.Val]; isFn {
			result := p.parseExpression()
			fmt.Print(result)
			return
		}
		if !p.isVar(p.CurrentToken.Val) {
			p.langError(p.CurrentToken, fmt.Sprintf("undefined: %s", p.CurrentToken.Val))
		}
		val := p.Variables[p.CurrentToken.Val]
		if val.Type == "string" {
			fmt.Print(val.Str)
			p.advanceParser()
			return
		}
		fmt.Print(p.parseExpression())
	case TOKEN_INT:
		result := p.parseExpression()
		fmt.Print(result)
	case TOKEN_FLOAT:
		result := p.parseExpression()
		fmt.Println(result)
	case TOKEN_NEW_LINE:
		fmt.Println()
	default:
		fmt.Print(p.parseExpression())
	}
}
