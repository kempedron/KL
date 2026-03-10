package tokens

import "fmt"

type TOKEN_TYPE int

const (
	TOKEN_INT   TOKEN_TYPE = iota
	TOKEN_PRINT            // print
	TOKEN_EOF
	TOKEN_IDENT                   // имя переменной
	TOKEN_SEMICOLON               // ;
	TOKEN_ASSIGN                  // =
	TOKEN_LPAREN                  // (
	TOKEN_RPAREN                  // )
	TOKEN_INT_KEYWORD             // int
	TOKEN_COMMA                   // ,
	TOKEN_SPACE                   // пробел
	TOKEN_STRING                  // type string
	TOKEN_STRING_KEYWORD          // 'string'
	TOKEN_PLUS                    // +
	TOKEN_MINUS                   // -
	TOKEN_MULTIPLY                // *
	TOKEN_DIVIDE                  // /
	TOKEN_MODULO                  // %
	TOKEN_QUOTATION               // "
	TOKEN_FLOAT                   // type float
	TOKEN_FLOAT_KEYWORD           // 'float'
	TOKEN_INPUT_KEYWORD           // 'input'
	TOKEN_NEW_LINE                //  \n
	TOKEN_FOR_FORMATTING_TO_STR   // %s
	TOKEN_FOR_FORMATTING_TO_INT   // %d
	TOKEN_FOR_FORMATTING_TO_FLOAT // %f
	TOKEN_PRINTF                  // printf
	TOKEN_FUNCTION_NAME           //func name
	TOKEN_FUNC_KEYWORD            //func
	TOKEN_RETURN_KEYWORD          //return
	TOKEN_LBRACE                  // {
	TOKEN_RBRACE                  // }

)

func (t TOKEN_TYPE) String() string {
	switch t {
	case TOKEN_INT:
		return "int literal"
	case TOKEN_FLOAT:
		return "float literal"
	case TOKEN_STRING:
		return "string literal"
	case TOKEN_IDENT:
		return "identifier"
	case TOKEN_ASSIGN:
		return "'='"
	case TOKEN_SEMICOLON:
		return "';'"
	case TOKEN_LPAREN:
		return "'('"
	case TOKEN_RPAREN:
		return "')'"
	case TOKEN_LBRACE:
		return "'{'"
	case TOKEN_RBRACE:
		return "'}'"
	case TOKEN_PLUS:
		return "'+'"
	case TOKEN_MINUS:
		return "'-'"
	case TOKEN_MULTIPLY:
		return "'*'"
	case TOKEN_DIVIDE:
		return "'/'"
	case TOKEN_MODULO:
		return "'%'"
	case TOKEN_COMMA:
		return "','"
	case TOKEN_EOF:
		return "end of file"
	case TOKEN_INT_KEYWORD:
		return "'int'"
	case TOKEN_FLOAT_KEYWORD:
		return "'float'"
	case TOKEN_STRING_KEYWORD:
		return "'string'"
	case TOKEN_FUNC_KEYWORD:
		return "'func'"
	case TOKEN_RETURN_KEYWORD:
		return "'return'"
	case TOKEN_PRINT:
		return "'print'"
	case TOKEN_PRINTF:
		return "'printf'"
	case TOKEN_INPUT_KEYWORD:
		return "'input'"
	default:
		return fmt.Sprintf("token(%d)", t)
	}
}
