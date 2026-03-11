package main

import (
	"fmt"
	. "kl/tokens"
	"strconv"
)

func (p *Parser) parseStringCast() Value {
	p.expect(TOKEN_STRING_KEYWORD)
	p.expect(TOKEN_LPAREN)

	switch p.CurrentToken.Type {
	case TOKEN_IDENT:
		if val, ok := p.Variables[p.expect(TOKEN_IDENT).Val]; ok {
			p.expect(TOKEN_RPAREN)
			switch val.Type {
			case "int":
				return Value{Type: "string", Str: strconv.Itoa(val.Int)}
			case "float":
				return Value{Type: "string", Str: strconv.FormatFloat(val.Float, 'f', -1, 64)}
			case "string":
				return Value{Type: "string", Str: val.Str}
			default:
				panic(fmt.Sprintf("unsupported type for string cast: %v", val.Type))
			}
		}

	case TOKEN_INT:
		val := p.expect(TOKEN_INT).Val
		p.expect(TOKEN_RPAREN)
		return Value{Type: "string", Str: val}
	case TOKEN_FLOAT:
		val := p.expect(TOKEN_FLOAT).Val
		p.expect(TOKEN_RPAREN)
		return Value{Type: "string", Str: val}
	}
	p.langError(p.CurrentToken, "invalid argument for string()")
	return Value{}
}
