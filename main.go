package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type TOKEN_TYPE int

const (
	TOKEN_INT   TOKEN_TYPE = iota
	TOKEN_PRINT            // print
	TOKEN_EOF
	TOKEN_IDENT       // имя переменной
	TOKEN_SEMICOLON   // ;
	TOKEN_ASSIGN      // =
	TOKEN_LPAREN      // (
	TOKEN_RPAREN      // )
	TOKEN_INT_KEYWORD // int
	TOKEN_COMMA       // ,
	TOKEN_SPACE       // пробел
	TOKEN_STRING
	TOKEN_STRING_KEYWORD
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULTIPLY
	TOKEN_DIVIDE
	TOKEN_MODULO
	TOKEN_QUOTATION
	TOKEN_FLOAT
	TOKEN_FLOAT_KEYWORD
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("использование ./kepr <файл с кодом>")
		os.Exit(1)
	}

	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("error read file %s: %v", os.Args[1], err)
		panic(err)
	}

	lexer := NewLexer(string(source))

	parser := NewParser(lexer)

	ParseProgramm(parser)
}

type Token struct {
	Type TOKEN_TYPE
	Val  string
}

type Lexer struct {
	CurrentPos  int
	InputString string
}

func NewLexer(inpStr string) *Lexer {
	return &Lexer{
		InputString: inpStr,
		CurrentPos:  0,
	}

}

func (l *Lexer) NextToken() Token {
	for l.CurrentPos < len(l.InputString) && unicode.IsSpace(rune(l.InputString[l.CurrentPos])) {
		l.CurrentPos++
	}

	if l.CurrentPos >= len(l.InputString) {
		return Token{Type: TOKEN_EOF}
	}

	if unicode.IsDigit(rune(l.InputString[l.CurrentPos])) {
		start := l.CurrentPos
		for l.CurrentPos < len(l.InputString) && unicode.IsDigit(rune(l.InputString[l.CurrentPos])) {
			l.CurrentPos++
		}
		return Token{TOKEN_INT, l.InputString[start:l.CurrentPos]}
	}

	if unicode.IsLetter(rune(l.InputString[l.CurrentPos])) {
		start := l.CurrentPos

		for l.CurrentPos < len(l.InputString) && (unicode.IsLetter(rune(l.InputString[l.CurrentPos]))) || unicode.IsDigit(rune(l.InputString[l.CurrentPos])) {
			l.CurrentPos++
		}
		value := l.InputString[start:l.CurrentPos]

		if value == "int" {
			return Token{TOKEN_INT_KEYWORD, value}
		}

		if value == "float" {
			return Token{TOKEN_FLOAT_KEYWORD, value}
		}

		if value == "print" {
			return Token{TOKEN_PRINT, value}
		}

		if value == "string" {
			return Token{TOKEN_STRING_KEYWORD, value}
		}
		return Token{TOKEN_IDENT, value}
	}

	ch := l.InputString[l.CurrentPos]

	l.CurrentPos++

	switch ch {
	case '=':
		return Token{TOKEN_ASSIGN, "="}
	case ';':
		return Token{TOKEN_SEMICOLON, ";"}
	case '(':
		return Token{TOKEN_LPAREN, "("}
	case ')':
		return Token{TOKEN_RPAREN, ")"}
	case '"':
		start := l.CurrentPos
		l.CurrentPos++

		for l.CurrentPos < len(l.InputString) && l.InputString[l.CurrentPos] != '"' {
			l.CurrentPos++
		}

		stringContent := l.InputString[start:l.CurrentPos]
		l.CurrentPos++
		return Token{TOKEN_STRING, stringContent}
	case ',':
		return Token{TOKEN_COMMA, ","}
	case '+':
		return Token{TOKEN_PLUS, "+"}
	case '-':
		return Token{TOKEN_MINUS, "-"}
	case '*':
		return Token{TOKEN_MULTIPLY, "*"}
	case '/':
		return Token{TOKEN_DIVIDE, "/"}
	case '%':
		return Token{TOKEN_MODULO, "%"}
	case ' ':
		fmt.Println("найден пробел")
		return Token{TOKEN_SPACE, " "}

	}

	panic(fmt.Sprintf("неизвестный символ: %c", ch))
}

type Parser struct {
	Lexer        *Lexer
	CurrentToken Token
	Variables    map[string]string
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{Lexer: lexer, Variables: make(map[string]string, 10)}
	p.CurrentToken = p.Lexer.NextToken()
	return p
}

func (p *Parser) expect(t TOKEN_TYPE) Token {
	if p.CurrentToken.Type != t {
		panic(fmt.Sprintf("ошибка: токены не совпадают( %v(have) != %v(want))", p.CurrentToken.Type, t))
	}
	token := p.CurrentToken
	p.advanceParser()

	return token
}

func (parser *Parser) advanceParser() {
	if parser.CurrentToken.Type == TOKEN_EOF {
		os.Exit(0)
	}
	parser.CurrentToken = parser.Lexer.NextToken()
}

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
		p.parseReAssignVar()
	default:
		panic(fmt.Sprintf("неизвестный токен: %v", p.CurrentToken.Type))
	}

}

func (p *Parser) parseReAssignVar() {
	varName := p.expect(TOKEN_IDENT).Val
	p.expect(TOKEN_ASSIGN)
	var varVal string
	switch p.CurrentToken.Type {
	case TOKEN_INT:
		varVal = p.expect(TOKEN_INT).Val
	case TOKEN_STRING:
		varVal = p.expect(TOKEN_STRING).Val
	case TOKEN_FLOAT:
		varVal = p.expect(TOKEN_FLOAT).Val
	}

	p.setVar(varName, varVal)
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
		fmt.Println()
		return

	}

	for {
		p.printArgument()
		if p.CurrentToken.Type == TOKEN_COMMA {
			p.expect(TOKEN_COMMA)
			fmt.Print(" ")
			continue
		}
		break
	}
	p.advanceParser()
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)

	fmt.Println()
}

func (p *Parser) printString() {
	p.expect(TOKEN_QUOTATION)

	for {
		fmt.Println(p.CurrentToken)
		p.printArgument()

		if p.CurrentToken.Type == TOKEN_COMMA {
			p.expect(TOKEN_COMMA)
			fmt.Print(" ")
			continue
		}
		break
	}

	p.expect(TOKEN_QUOTATION)
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)
	fmt.Println("закончил выводить строку")
}

func (p *Parser) isVar(varName string) bool {
	if _, ok := p.Variables[varName]; ok {
		return true
	}
	return false
}

func (p *Parser) printArgument() {
	switch p.CurrentToken.Type {
	case TOKEN_STRING:
		fmt.Println("это string")
		token := p.expect(TOKEN_STRING)
		fmt.Print(token.Val)
	case TOKEN_IDENT:
		if !p.isVar(p.CurrentToken.Val) {
			panic(fmt.Sprintf("undeclarated var: %v", p.CurrentToken.Val))
		}
		val := p.Variables[p.CurrentToken.Val]
		fmt.Print(val)
	case TOKEN_INT:
		result := p.parseExpression()
		fmt.Print(result)
	}
}

func (p *Parser) parseInt() {
	p.expect(TOKEN_INT_KEYWORD)

	ident := p.expect(TOKEN_IDENT)
	p.expect(TOKEN_ASSIGN)
	result := p.parseExpression()
	p.expect(TOKEN_SEMICOLON)
	p.setVar(ident.Val, strconv.Itoa(result))
}

func (p *Parser) parseFloat() {
	p.expect(TOKEN_FLOAT_KEYWORD)

	ident := p.expect(TOKEN_IDENT)
	p.expect(TOKEN_ASSIGN)
	result := p.parseExpression()
	p.expect(TOKEN_SEMICOLON)
	p.setVar(ident.Val, strconv.Itoa(result))
}

func (p *Parser) parseStr() {
	p.expect(TOKEN_STRING_KEYWORD)
	ident := p.expect(TOKEN_IDENT)
	p.expect(TOKEN_ASSIGN)
	t := p.expect(TOKEN_STRING)

	p.expect(TOKEN_SEMICOLON)
	p.setVar(ident.Val, t.Val)
}

func (p *Parser) setVar(varname string, val string) {
	p.Variables[varname] = val
}

func ParseProgramm(parser *Parser) {
	for parser.CurrentToken.Type != TOKEN_EOF {
		parser.Parse()
	}
}

func (p *Parser) parseExpression() int {
	return p.parseAddSub()
}
func (p *Parser) parseAddSub() int {
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

func (p *Parser) parseMulDiv() int {
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
			left = left % right
		}
	}
	return left
}

func (p *Parser) parsePrimary() int {
	switch p.CurrentToken.Type {
	case TOKEN_INT:
		token := p.expect(TOKEN_INT)
		val, err := strconv.Atoi(token.Val)
		if err != nil {
			panic(fmt.Sprintf("невалидное число: %s", token.Val))
		}
		return val
	case TOKEN_IDENT:
		token := p.expect(TOKEN_IDENT)
		if val, ok := p.Variables[token.Val]; ok {
			if valI, err := strconv.Atoi(val); err == nil {
				return valI
			}
			panic(fmt.Sprintf("невалидное число: %s", token.Val))
		}
		panic(fmt.Sprintf("undeclarated var: %s", token.Val))

	case TOKEN_LPAREN:
		p.expect(TOKEN_LPAREN)
		result := p.parseExpression()
		p.expect(TOKEN_RPAREN)
		return result

	case TOKEN_MINUS:
		p.expect(TOKEN_MINUS)
		return -p.parsePrimary()

	case TOKEN_PLUS:
		p.expect(TOKEN_PLUS)
		return p.parsePrimary()

	default:
		panic(fmt.Sprintf("неожиданный токен в выражении: %v", p.CurrentToken.Type))
	}

}
