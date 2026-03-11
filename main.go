package main

import (
	"fmt"
	. "kl/tokens"
	"os"
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

func ParseProgramm(parser *Parser) {
	for parser.CurrentToken.Type != TOKEN_EOF {
		parser.Parse()
	}
}
