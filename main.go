package main

import (
	"fmt"
	. "kl/tokens"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", r)
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("использование ./kl <файл .kl>")
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
	Line int
	Col  int
}

type Lexer struct {
	CurrentPos  int
	InputString string
	Line        int
	Col         int
}

type Value struct {
	Type  string
	Int   int
	Float float64
	Str   string
}

type Function struct {
	Name   string
	Params []string
	Body   []Token
}

type LangError struct {
	Line    int
	Col     int
	Message string
}
type ILexer interface {
	NextToken() Token
	GetNextToken() Token
}

func (p *Parser) langError(tok Token, msg string) {
	panic(fmt.Sprintf("%d:%d: %s", tok.Line+1, tok.Col+1, msg))
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
		return Token{TOKEN_IDENT, value, line, col}
	}

	ch := l.InputString[l.CurrentPos]
	line, col := l.Line, l.Col

	l.advance()

	switch ch {
	case '=':
		return Token{TOKEN_ASSIGN, "=", line, col}
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

func (p *Parser) ParsePrintf() {
	p.expect(TOKEN_PRINTF)
	p.expect(TOKEN_LPAREN)
	var args []any
	strToFormatting := p.expect(TOKEN_STRING).Val
	p.expect(TOKEN_COMMA)
	arg := p.expect(TOKEN_IDENT).Val
	varVal := p.GetVar(arg)
	var argVal any
	switch varVal.Type {
	case "string":
		argVal = varVal.Str
	case "int":
		argVal = varVal.Int
	case "float":
		argVal = varVal.Float
	default:
		argVal = varVal
	}
	args = append(args, argVal)

	for p.CurrentToken.Type != TOKEN_RPAREN {
		p.expect(TOKEN_COMMA)
		arg := p.expect(TOKEN_IDENT).Val
		varVal := p.GetVar(arg)
		var argVal any
		switch varVal.Type {
		case "string":
			argVal = varVal.Str
		case "int":
			argVal = varVal.Int
		case "float":
			argVal = varVal.Float
		default:
			argVal = varVal
		}
		args = append(args, argVal)
	}

	fmt.Printf(strToFormatting, args...)
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)

}

//здесь точка нужная сюда (для поиска)

func (p *Parser) Parse() {
	//println(p.CurrentToken.Type)
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
	default:
		panic(fmt.Sprintf("неизвестный токен: %v", p.CurrentToken.Type))
	}

}

func (p *Parser) parseReAssignVar() {
	varName := p.expect(TOKEN_IDENT).Val
	p.expect(TOKEN_ASSIGN)
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

	// if p.CurrentToken.Type == TOKEN_COMMA {
	// 	fmt.Println("")
	// 	p.expect(TOKEN_COMMA)
	// 	panic("запятая")
	// }

	p.expect(TOKEN_QUOTATION)
	p.expect(TOKEN_RPAREN)
	p.expect(TOKEN_SEMICOLON)
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
		token := p.expect(TOKEN_STRING)
		fmt.Print(token.Val)
	case TOKEN_IDENT:
		if !p.isVar(p.CurrentToken.Val) {
			panic(fmt.Sprintf("undeclarated var: %v", p.CurrentToken.Val))
		}
		val := p.Variables[p.CurrentToken.Val]
		switch val.Type {
		case "int":
			fmt.Print(val.Int)
		case "float":
			fmt.Println(val.Float)
		case "string":
			fmt.Print(val.Str)
		}
		p.advanceParser()
	case TOKEN_INT:
		result := p.parseExpression()
		fmt.Print(result)
	case TOKEN_FLOAT:
		result := p.parseExpression()
		fmt.Println(result)
	case TOKEN_NEW_LINE:
		fmt.Println()
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
	p.checkTypeAssignment(ident.Val, "string", p.CurrentToken)
	t := p.expect(TOKEN_STRING)

	p.expect(TOKEN_SEMICOLON)
	p.setVar(ident.Val, Value{Type: "string", Str: t.Val})
}

func (p *Parser) setVar(varname string, val Value) {
	p.Variables[varname] = val
}

func ParseProgramm(parser *Parser) {
	for parser.CurrentToken.Type != TOKEN_EOF {
		parser.Parse()
	}
}

type Number interface {
	int | float64
}

func (p *Parser) parseExpression() float64 {
	return p.parseAddSub()
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
			panic(fmt.Sprintf("невалидное число: %s", token.Val))
		}
		panic(fmt.Sprintf("undeclarated var: %s", token.Val))

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

	default:
		panic(fmt.Sprintf("неожиданный токен в выражении: %v", p.CurrentToken.Type))
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

	if p.isVar(ident) {
		var input string
		fmt.Println(text)
		fmt.Scan(&input)
		p.Variables[ident] = Value{Type: "string", Str: input}
	} else {
		panic(fmt.Sprintf("undeclared var: %v", ident))
	}
}

func (p *Parser) GetVar(varname string) Value {
	if val, ok := p.Variables[varname]; ok {
		return val
	}
	panic(fmt.Sprintf("undeclared var: %v", varname))
}

func (p *Parser) ParseFunc() {
	p.expect(TOKEN_FUNC_KEYWORD)
	funcName := p.expect(TOKEN_IDENT).Val
	p.expect(TOKEN_LPAREN)
	var params []string
	for p.CurrentToken.Type != TOKEN_RPAREN {
		param := p.expect(TOKEN_IDENT).Val
		params = append(params, param)
		if p.CurrentToken.Type != TOKEN_RPAREN {
			p.expect(TOKEN_COMMA)
		}

	}
	p.expect(TOKEN_RPAREN)
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
	p.Functions[funcName] = Function{Name: funcName, Params: params, Body: body}

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

func (p *Parser) callFuncExpr(name string) float64 {
	fn, ok := p.Functions[name]
	if !ok {
		panic(fmt.Sprintf("Unknow functions: %s", name))
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
		panic(fmt.Sprintf("wrong number of arguments: %d != %d", len(args), len(fn.Params)))
	}

	bodyLexer := &TokenSliceLexer{Tokens: fn.Body, Pos: 0}
	subParser := &Parser{
		Lexer:     bodyLexer,
		Variables: make(map[string]Value),
		Functions: p.Functions,
	}
	for i, param := range fn.Params {
		subParser.Variables[param] = args[i]
	}
	subParser.CurrentToken = bodyLexer.NextToken()

	for subParser.CurrentToken.Type != TOKEN_EOF {
		if subParser.CurrentToken.Type == TOKEN_RETURN_KEYWORD {
			subParser.advanceParser()
			return subParser.parseExpression()
		}
		subParser.Parse()
	}
	return 0
}

func (p *Parser) parseArgValue() Value {
	switch p.CurrentToken.Type {
	case TOKEN_STRING:
		return Value{Type: "string", Str: p.expect(TOKEN_STRING).Val}
	case TOKEN_IDENT:
		if val, ok := p.Variables[p.CurrentToken.Val]; ok && val.Type == "string" {
			p.advanceParser()
			return val
		}
		return Value{Type: "float", Float: p.parseExpression()}
	default:
		return Value{Type: "float", Float: p.parseExpression()}
	}
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
