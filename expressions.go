package main

import (
	"fmt"
	. "kl/tokens"
	"strconv"
)

func (p *Parser) parseArgValue() Value {
	switch p.CurrentToken.Type {
	case TOKEN_STRING:
		return Value{Type: "string", Str: p.expect(TOKEN_STRING).Val}
	case TOKEN_INT:
		val, _ := strconv.Atoi(p.CurrentToken.Val)
		p.advanceParser()
		return Value{Type: "int", Int: val}
	case TOKEN_FLOAT:
		val, _ := strconv.ParseFloat(p.CurrentToken.Val, 64)
		p.advanceParser()
		return Value{Type: "float", Float: val}
	case TOKEN_INT_KEYWORD:
		result := p.parsePrimary()
		return Value{Type: "int", Int: int(result)}
	case TOKEN_FLOAT_KEYWORD:
		result := p.parsePrimary()
		return Value{Type: "float", Float: result}
	case TOKEN_IDENT:
		if val, ok := p.Variables[p.CurrentToken.Val]; ok { //&& val.Type == "string"
			p.advanceParser()
			return val
		}
		return Value{Type: "float", Float: p.parseExpression()}

	default:
		return Value{Type: "float", Float: p.parseExpression()}
	}
}

type Number interface {
	int | float64
}

func (p *Parser) parseExpression() float64 {
	return p.parseOr()
}

func (p *Parser) parseOr() float64 {
	left := p.parseAnd()
	for p.CurrentToken.Type == TOKEN_OR {
		p.advanceParser()
		right := p.parseAnd()
		if left != 0 || right != 0 {
			left = 1
		} else {
			left = 0
		}
	}
	return left
}

func (p *Parser) parseAnd() float64 {
	left := p.parseEquality()
	for p.CurrentToken.Type == TOKEN_AND {
		p.advanceParser()
		right := p.parseEquality()
		if left != 0 && right != 0 {
			left = 1
		} else {
			left = 0
		}
	}
	return left
}

// ==, !=
func (p *Parser) parseEquality() float64 {
	left := p.parseComparison()
	for p.CurrentToken.Type == TOKEN_EQ || p.CurrentToken.Type == TOKEN_NEQ {
		op := p.CurrentToken.Type
		p.advanceParser()
		right := p.parseComparison()

		if op == TOKEN_EQ {
			if left == right {
				left = 1
			} else {
				left = 0
			}
		} else {
			if left != right {
				left = 1
			} else {
				left = 0
			}
		}
	}
	return left
}

// >, <, >=, <=
func (p *Parser) parseComparison() float64 {
	left := p.parseAddSub()

	for p.CurrentToken.Type == TOKEN_EQ || p.CurrentToken.Type == TOKEN_NEQ || p.CurrentToken.Type == TOKEN_LT || p.CurrentToken.Type == TOKEN_GT || p.CurrentToken.Type == TOKEN_LTE || p.CurrentToken.Type == TOKEN_GTE {
		op := p.CurrentToken.Type
		p.advanceParser()
		right := p.parseAddSub()

		switch op {
		case TOKEN_EQ:
			if left == right {
				left = 1
			} else {
				left = 0
			}
		case TOKEN_NEQ:
			if left != right {
				left = 1
			} else {
				left = 0
			}
		case TOKEN_LT:
			if left < right {
				left = 1
			} else {
				left = 0
			}
		case TOKEN_GT:
			if left > right {
				left = 1
			} else {
				left = 0
			}
		case TOKEN_LTE:
			if left <= right {
				left = 1
			} else {
				left = 0
			}
		case TOKEN_GTE:
			if left >= right {
				left = 1
			} else {
				left = 0
			}
		}
	}
	return left
}

func (p *Parser) parseAddSub() float64 {
	left := p.parseMulDiv()

	for p.CurrentToken.Type == TOKEN_PLUS || p.CurrentToken.Type == TOKEN_MINUS {

		op := p.CurrentToken.Type
		p.advanceParser()
		right := p.parseMulDiv()

		if op == TOKEN_PLUS {
			left = left + right
		} else {

			left = left - right
		}

	}
	return left
}

func (p *Parser) parseMulDiv() float64 {
	left := p.parsePrimary()

	for p.CurrentToken.Type == TOKEN_MULTIPLY || p.CurrentToken.Type == TOKEN_DIVIDE || p.CurrentToken.Type == TOKEN_MODULO {
		op := p.CurrentToken.Type
		p.advanceParser()
		right := p.parsePrimary()

		switch op {
		case TOKEN_MULTIPLY:
			left = left * right
		case TOKEN_DIVIDE:
			if right == 0 {
				panic("деление на 0")
			}
			left = left / right
		case TOKEN_MODULO:
			if right == 0 {
				panic("остаток от деления на 0")
			}
			left = float64(int(left) % int(right))

		}
	}
	return left
}

func (p *Parser) parsePrimary() float64 {
	switch p.CurrentToken.Type {
	case TOKEN_INT:
		token := p.expect(TOKEN_INT)
		val, err := strconv.Atoi(token.Val)
		if err != nil {
			panic(fmt.Sprintf("невалидное число для int: %s", token.Val))
		}
		return float64(val)

	case TOKEN_FLOAT:
		token := p.expect(TOKEN_FLOAT)
		val, err := strconv.ParseFloat(token.Val, 64)
		if err != nil {
			panic(fmt.Sprintf("невалидное число для float: %s", token.Val))
		}
		return val
	case TOKEN_IDENT:
		token := p.expect(TOKEN_IDENT)
		if _, isFn := p.Functions[token.Val]; isFn {
			return p.callFuncExpr(token.Val)
		}
		if val, ok := p.Variables[token.Val]; ok {
			if val.Type == "float" {
				return val.Float
			}
			if val.Type == "int" {
				return float64(val.Int)
			}
			p.langError(token, fmt.Sprintf("cannot use '%s' (type string) as numeric value", token.Val))
		}
		p.langError(token, fmt.Sprintf("undefined: %s", token.Val))

	case TOKEN_LPAREN:
		p.expect(TOKEN_LPAREN)
		result := p.parseExpression()
		p.expect(TOKEN_RPAREN)
		return float64(result)

	case TOKEN_MINUS:
		p.expect(TOKEN_MINUS)
		return -p.parsePrimary()

	case TOKEN_PLUS:
		p.expect(TOKEN_PLUS)
		return p.parsePrimary()

	case TOKEN_INT_KEYWORD:
		p.advanceParser()
		p.expect(TOKEN_LPAREN)
		if p.CurrentToken.Type == TOKEN_IDENT {
			if val, ok := p.Variables[p.CurrentToken.Val]; ok && val.Type == "string" {
				tok := p.CurrentToken
				p.advanceParser()
				p.expect(TOKEN_RPAREN)
				result, err := strconv.Atoi(val.Str)
				if err != nil {
					p.langError(tok, fmt.Sprintf("cannot convert %q to int", val.Str))
				}
				return float64(result)
			}
		}
		val := p.parseExpression()
		p.expect(TOKEN_RPAREN)
		return float64(int(val))

	case TOKEN_FLOAT_KEYWORD:
		p.advanceParser()
		p.expect(TOKEN_LPAREN)
		if p.CurrentToken.Type == TOKEN_IDENT {
			if val, ok := p.Variables[p.CurrentToken.Val]; ok && val.Type == "string" {
				tok := p.CurrentToken
				p.advanceParser()
				p.expect(TOKEN_RPAREN)
				result, err := strconv.ParseFloat(val.Str, 64)
				if err != nil {
					p.langError(tok, fmt.Sprintf("cannot convert %q to float", val.Str))
				}
				return result
			}
		}
		val := p.parseExpression()
		p.expect(TOKEN_RPAREN)
		return val

	default:
		panic(fmt.Sprintf("неожиданный токен в выражении: %v", p.CurrentToken.Type))
	}
	return 0
}
