package tokens

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

)
