package scanner

import (
	"bufio"
	"io"
	"log"
	"proj/golite/token"
	"regexp"
)

type Scanner struct {
	reader *bufio.Reader

	// accumulated valid string, everytime the scanner returns a token
	// it will be cleaned up (re-set to an empty string)
	lexeme string

	number_compiled *regexp.Regexp
	id_compiled     *regexp.Regexp
	whitespaces     *regexp.Regexp

	keywords map[string]token.TokenType
	symbols  map[string]token.TokenType

	isComment  bool
	lineNumber int
}

func New(input_reader *bufio.Reader) *Scanner {
	scanner := &Scanner{reader: input_reader, lexeme: ""}
	scanner.number_compiled, _ = regexp.Compile("^[0-9]+$")
	scanner.id_compiled, _ = regexp.Compile("^[a-zA-Z][a-zA-Z0-9]*$")
	scanner.whitespaces, _ = regexp.Compile("\\s+")
	scanner.isComment = false
	scanner.lineNumber = 1

	keywords_map := map[string]token.TokenType{
		"int":    token.INT,
		"number": token.NUM,
		"bool":   token.BOOL,
		"true":   token.TRUE,
		"false":  token.FALSE,
		"id":     token.ID,
		"nil":    token.NIL,

		"let":     token.LET,
		"Print":   token.PRINT,
		"Println": token.PRINTLN,
		"return":  token.RETURN,
		"package": token.PACK,
		"import":  token.IMPORT,
		"fmt":     token.FMT,
		"type":    token.TYPE,
		"struct":  token.STRUCT,
		"Scan":    token.SCAN,
		"if":      token.IF,
		"else":    token.ELSE,
		"for":     token.FOR,
		"func":    token.FUNC,
		"var":     token.VAR,
	}

	symbols_map := map[string]token.TokenType{
		".":  token.DOT,
		",":  token.COMMA,
		"\"": token.QTDMARK,
		"{":  token.LBRACE,
		"}":  token.RBRACE,
		"(":  token.LPAREN,
		")":  token.RPAREN,

		"=":  token.ASSIGN,
		"&":  token.AMPERS,
		";":  token.SEMICOLON,
		"+":  token.ADD,
		"-":  token.MINUS,
		"*":  token.MULTIPLY,
		"/":  token.DIVIDE,
		"||": token.OR,
		"&&": token.AND,
		"!":  token.NOT,
		"==": token.EQUAL,
		"!=": token.NEQUAL,
		">":  token.GT,
		">=": token.GE,
		"<":  token.LT,
		"<=": token.LE,
		"//": token.COMMENT,
	}

	// type_lookup := map[string]string{
	// 	"ILLEGAL": "ILLEGAL",
	// 	"eof":     "EOF",
	// 	"int":     "INT",
	// 	"number":  "NUM",
	// 	"bool":    "BOOL",
	// 	"true":    "TRUE",
	// 	"false":   "FALSE",
	// 	"id":      "ID",
	// 	"nil":     "NIL",

	// 	"let":     "LET",
	// 	"Print":   "PRINT",
	// 	"Println": "PRINTLN",
	// 	"return":  "RETURN",
	// 	"package": "PACK",
	// 	"import":  "IMPORT",
	// 	"fmt":     "FMT",
	// 	"type":    "TYPE",
	// 	"struct":  "STRUCT",
	// 	"Scan":    "SCAN",
	// 	"if":      "IF",
	// 	"else":    "ELSE",
	// 	"for":     "FOR",
	// 	"func":    "FUNC",
	// 	"var":     "VAR",

	// 	".":  "DOT",
	// 	",":  "COMMA",
	// 	"\"": "QTDMARK",
	// 	"{":  "LBRACE",
	// 	"}":  "RBRACE",
	// 	"(":  "LPAREN",
	// 	")":  "RPAREN",

	// 	"=":  "ASSIGN",
	// 	"&":  "AMPERS",
	// 	";":  "SEMICOLON",
	// 	"+":  "ADD",
	// 	"-":  "MINUS",
	// 	"*":  "MULTIPLY",
	// 	"/":  "DIVIDE",
	// 	"||": "OR",
	// 	"&&": "AND",
	// 	"!":  "NOT",
	// 	"==": "EQUAL",
	// 	"!=": "NEQUAL",
	// 	">":  "GT",
	// 	">=": "GE",
	// 	"<":  "LT",
	// 	"<=": "LE",
	// 	"//": "COMMENT",
	// }
	scanner.keywords = keywords_map
	scanner.symbols = symbols_map
	return scanner
}

func (l *Scanner) NextToken() token.Token {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// return 'eof' if we have not processed any chars as current literal (lexeme)
				if len(l.lexeme) == 0 {
					return token.Token{Type: token.EOF, Literal: "eof", LineNum: l.lineNumber}
				}
				if l.number_compiled.MatchString(l.lexeme) {
					temp_lexeme := l.lexeme
					l.lexeme = ""
					return token.Token{Type: token.NUM, Literal: temp_lexeme, LineNum: l.lineNumber}
				}
				if tok, exist := l.symbols[l.lexeme]; exist {
					temp_lexeme := l.lexeme
					l.lexeme = ""
					return token.Token{Type: tok, Literal: temp_lexeme, LineNum: l.lineNumber}
				}
				if l.id_compiled.MatchString(l.lexeme) {
					temp_lexeme := l.lexeme
					// check if it matches with some keywords (e.g. print, var)
					if tok, exist := l.keywords[l.lexeme]; exist {
						l.lexeme = ""
						return token.Token{Type: tok, Literal: temp_lexeme, LineNum: l.lineNumber}
					}
					l.lexeme = ""
					return token.Token{Type: token.ID, Literal: temp_lexeme, LineNum: l.lineNumber}
				}

			} else {
				// unknown error
				log.Fatal(err)
			}
		} else {
			temp_line_number := l.lineNumber
			if l.isComment {
				if r == '\n' || r == '\r' { // newline
					l.isComment = false
					r_next, _, _ := l.reader.ReadRune()
					if r_next != '\n' {
						l.reader.UnreadRune()
					}
					l.lineNumber++
				}
				continue
			}
			curr_lexeme := l.lexeme + string(r)
			_, exist := l.symbols[curr_lexeme]
			if l.number_compiled.MatchString(curr_lexeme) || l.id_compiled.MatchString(curr_lexeme) || exist {
				l.lexeme = curr_lexeme
				continue
			}

			if len(l.lexeme) == 0 && !l.whitespaces.MatchString(string(r)) {
				return token.Token{Type: token.ILLEGAL, Literal: "ILLEGAL", LineNum: temp_line_number}
			}

			// if r != ' ' && r != '\n' && r != '\t' {
			// 	l.reader.UnreadRune() // rollback
			// }
			// rollback if not a whitespace/newline/carriage return/etc
			if !l.whitespaces.MatchString(string(r)) {
				l.reader.UnreadRune()
			} else if r == '\n' || r == '\r' { // newline
				r_next, _, _ := l.reader.ReadRune()
				if r_next != '\n' {
					l.reader.UnreadRune()
				}
				l.lineNumber++
			}

			temp_lexeme := l.lexeme
			l.lexeme = ""
			if l.number_compiled.MatchString(temp_lexeme) {
				return token.Token{Type: token.NUM, Literal: temp_lexeme, LineNum: temp_line_number}
			}
			if l.id_compiled.MatchString(temp_lexeme) {
				if tok, exist := l.keywords[temp_lexeme]; exist {
					return token.Token{Type: tok, Literal: temp_lexeme, LineNum: temp_line_number}
				}
				return token.Token{Type: token.ID, Literal: temp_lexeme, LineNum: temp_line_number}
			}
			if tok, exist := l.symbols[temp_lexeme]; exist {
				if tok == token.COMMENT {
					l.isComment = true
				}
				return token.Token{Type: tok, Literal: temp_lexeme, LineNum: temp_line_number}
			}
		}
	}

}