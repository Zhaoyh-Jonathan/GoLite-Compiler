package token

import (
	"fmt"
)

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "eof"
	INT     = "int"
	NUM     = "number" // an integer containing one or more digits (0-9)
	BOOL    = "bool"
	TRUE    = "true"
	FALSE   = "false"
	ID      = "id"
	NIL     = "nil"

	LET     = "let"
	PRINT   = "Print"
	PRINTLN = "Println"
	RETURN  = "return"
	PACK    = "package"
	IMPORT  = "import"
	FMT     = "fmt"
	TYPE    = "type"
	STRUCT  = "struct"
	SCAN    = "Scan"
	IF      = "if"
	ELSE    = "else"
	FOR     = "for"
	FUNC    = "func"
	VAR     = "var"

	DOT     = "."
	COMMA   = ","
	QTDMARK = "\""
	LBRACE  = "{"
	RBRACE  = "}"
	LPAREN  = "("
	RPAREN  = ")"

	ASSIGN    = "="
	AMPERS    = "&" // for getting address
	SEMICOLON = ";"
	ADD       = "+"
	MINUS     = "-"
	MULTIPLY  = "*"
	DIVIDE    = "/"
	OR        = "||"
	AND       = "&&"
	NOT       = "!"
	EQUAL     = "=="
	NEQUAL    = "!="
	GT        = ">"
	GE        = ">="
	LT        = "<"
	LE        = "<="
	COMMENT   = "//"
)

type Token struct {
	Type        TokenType
	Literal     string
	type_lookup string
}

func (tok Token) String() string {
	return fmt.Sprintf("%T, %v", tok.Type, tok.Literal)
}
