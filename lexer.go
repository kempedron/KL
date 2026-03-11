package main

import (
	"fmt"
	. "kl/tokens"
	"strings"
	"unicode"
)

type Token struct {
	Type TOKEN_TYPE
	Val  string
	Line int
	Col  int
}

type Lexer struct {
	CurrentPos  int
	InputString string
	Line        int
	Col         int
}

func NewLexer(inpStr string) *Lexer {
	return &Lexer{
		InputString: inpStr,
		CurrentPos:  0,
	}

}

func (l Lexer) GetNextToken() Token {
	return l.NextToken()
}

func (l *Lexer) advance() {
	if l.CurrentPos < len(l.InputString) {
		if l.InputString[l.CurrentPos] == '\n' {
			l.Line++
			l.Col = 0
		} else {
			l.Col++
		}
		l.CurrentPos++
	}
}

func (l *Lexer) NextToken() Token {
	for l.CurrentPos < len(l.InputString) && unicode.IsSpace(rune(l.InputString[l.CurrentPos])) {
		l.advance()
	}

	if l.CurrentPos >= len(l.InputString) {
		return Token{Type: TOKEN_EOF}
	}

	if unicode.IsDigit(rune(l.InputString[l.CurrentPos])) {
		line, col := l.Line, l.Col
		start := l.CurrentPos
		hasDecimal := false
		for l.CurrentPos < len(l.InputString) && (unicode.IsDigit(rune(l.InputString[l.CurrentPos])) || l.InputString[l.CurrentPos] == '.') {
			if l.InputString[l.CurrentPos] == '.' {
				if hasDecimal {
					break
				}
				hasDecimal = true
			}
			l.advance()

		}
		value := l.InputString[start:l.CurrentPos]
		if hasDecimal {
			return Token{TOKEN_FLOAT, value, line, col}
		}

		return Token{TOKEN_INT, value, line, col}
	}

	if unicode.IsLetter(rune(l.InputString[l.CurrentPos])) {
		start := l.CurrentPos
		line, col := l.Line, l.Col

		for l.CurrentPos < len(l.InputString) && (unicode.IsLetter(rune(l.InputString[l.CurrentPos]))) || unicode.IsDigit(rune(l.InputString[l.CurrentPos])) {
			l.advance()
		}
		value := l.InputString[start:l.CurrentPos]

		if value == "int" {
			return Token{TOKEN_INT_KEYWORD, value, line, col}
		}

		if value == "float" {
			return Token{TOKEN_FLOAT_KEYWORD, value, line, col}
		}

		if value == "print" {
			return Token{TOKEN_PRINT, value, line, col}
		}

		if value == "string" {
			return Token{TOKEN_STRING_KEYWORD, value, line, col}
		}
		if value == "input" {
			return Token{TOKEN_INPUT_KEYWORD, value, line, col}
		}
		if value == "printf" {
			return Token{TOKEN_PRINTF, value, line, col}
		}
		if value == "func" {
			return Token{TOKEN_FUNC_KEYWORD, value, line, col}
		}
		if value == "return" {
			return Token{TOKEN_RETURN_KEYWORD, value, line, col}
		}
		if value == "if" {
			return Token{TOKEN_IF_KEYWORD, value, line, col}
		}
		if value == "else" {
			return Token{TOKEN_ELSE_KEYWORD, value, line, col}
		}
		if value == "elif" {
			return Token{TOKEN_ELIF_KEYWORD, value, line, col}
		}

		return Token{TOKEN_IDENT, value, line, col}
	}

	ch := l.InputString[l.CurrentPos]
	line, col := l.Line, l.Col

	l.advance()

	switch ch {
	case '=':
		if l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] == '=' {
			l.advance()
			return Token{TOKEN_EQ, "==", line, col}
		}
		return Token{TOKEN_ASSIGN, "=", line, col}
	case '!':
		if l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] == '=' {
			l.advance()
			return Token{TOKEN_NEQ, "!=", line, col}
		}
		panic("unexpected character: '!'")

	case '&':
		if l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] == '&' {
			l.advance()
			return Token{TOKEN_AND, "&&", line, col}
		}
		panic("unexpected character: '&'")
	case '|':
		if l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] == '|' {
			l.advance()
			return Token{TOKEN_OR, "||", line, col}
		}
		panic("unexpected character: '|'")
	case '>':
		if l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] == '=' {
			l.advance()
			return Token{TOKEN_GTE, ">=", line, col}
		}
		return Token{TOKEN_GT, ">", line, col}
	case '<':
		if l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] == '=' {
			l.advance()
			return Token{TOKEN_LTE, "<=", line, col}
		}
		return Token{TOKEN_LT, "<", line, col}
	case ';':
		return Token{TOKEN_SEMICOLON, ";", line, col}
	case '(':
		return Token{TOKEN_LPAREN, "(", line, col}
	case ')':
		return Token{TOKEN_RPAREN, ")", line, col}
	case '{':
		return Token{TOKEN_LBRACE, "{", line, col}
	case '}':
		return Token{TOKEN_RBRACE, "}", line, col}
	case '"':

		var sb strings.Builder

		for l.CurrentPos < len(l.InputString) {

			ch := l.InputString[l.CurrentPos]
			if ch == '"' {
				l.advance()
				break
			}

			if ch == '\\' {
				l.CurrentPos++
				if l.CurrentPos >= len(l.InputString) {
					panic("незакрытая escape последовательность")
				}
				switch l.InputString[l.CurrentPos] {
				case 'n':
					sb.WriteRune('\n')
				case 't':
					sb.WriteRune('\t')
				case '\\':
					sb.WriteRune('\\')
				case '"':
					sb.WriteRune('"')
				default:
					sb.WriteByte('\\')
					sb.WriteByte(l.InputString[l.CurrentPos])
				}
				l.advance()
			} else {
				sb.WriteByte(ch)
				l.advance()
			}

		}

		return Token{TOKEN_STRING, sb.String(), line, col}

	case ',':
		return Token{TOKEN_COMMA, ",", line, col}
	case '+':
		return Token{TOKEN_PLUS, "+", line, col}
	case '-':
		return Token{TOKEN_MINUS, "-", line, col}
	case '*':
		return Token{TOKEN_MULTIPLY, "*", line, col}
	case '/':
		return Token{TOKEN_DIVIDE, "/", line, col}
	case '%':
		return Token{TOKEN_MODULO, "%", line, col}
	case ' ':
		return Token{TOKEN_SPACE, " ", line, col}
	case '\n':
		return Token{TOKEN_NEW_LINE, "\n", line, col}

	}

	panic(fmt.Sprintf("неизвестный символ: %c", ch))
}

type TokenSliceLexer struct {
	Tokens []Token
	Pos    int
}

func (t *TokenSliceLexer) NextToken() Token {
	if t.Pos >= len(t.Tokens) {
		return Token{Type: TOKEN_EOF}
	}
	tok := t.Tokens[t.Pos]
	t.Pos++
	return tok
}

func (t *TokenSliceLexer) GetNextToken() Token {
	if t.Pos >= len(t.Tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return t.Tokens[t.Pos]
}
