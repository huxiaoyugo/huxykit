package token

import (
	"strconv"
	"strings"
	"unicode"
)

type Token int

/*
+ STRING,
	+ INT: 被当做INT32处理，
	+ INT32,
	+ INT64,
	+ UINT32,
	+ UINT64,
	+ FLOAT,
	+ DOUBLE,
	+ BOOL(或者BOOLEAN)

[ {{ARRAY_TYPE}} ]：数组，ARRAY_TYPE必须用{{}}包裹起来（兼容老版本）
	+ MAP<K_TYPE, V_TYPE>： MAP，K_TYPE,和V_TYPE不用{{}}包裹，
	+ ENUM<ENUM_TYPE>:
*/
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT
	EMPTY

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)

	INTEGER // 100
	STR     // "abc" or 'abc'
	IDENT   // 标识符 变量名
	literal_end

	operator_beg
	COLON     //  :
	LBRACK    // [
	LBRACE    // {
	RBRACK    // ]
	RBRACE    // }
	COMMA     // ,
	SEMICOLON // ;
	LSS       // <
	GTR       // >
	operator_end

	// Keywords
	keyword_beg
	INT
	INT32
	UINT32
	INT64
	UINT64
	FLOAT
	DOUBLE
	BOOL
	BOOLEAN
	STRING
	ENUM
	MAP
	PACKAGE
	IMPORT
	STRUCT
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "eof",
	COMMENT: "comment",
	EMPTY:   "empty",
	IDENT:   "identifier",
	INTEGER: "integer",
	STR:     "string",


	LSS:       "<",
	GTR:       ">",
	COLON:     ":",
	SEMICOLON: ";",
	LBRACK:    "[",
	LBRACE:    "{",
	RBRACK:    "]",
	RBRACE:    "}",
	COMMA:     ",",

	// key_world
	INT:     "INT",
	STRING:  "STRING",
	INT32:   "INT32",
	UINT32:  "UINT32",
	INT64:   "INT64",
	UINT64:  "UINT64",
	FLOAT:   "FLOAT",
	DOUBLE:  "DOUBLE",
	BOOL:    "BOOL",
	BOOLEAN: "BOOLEAN",
	ENUM:    "ENUM",
	STRUCT:  "STRUCT",
	MAP:     "MAP",
	IMPORT:  "IMPORT",
	PACKAGE: "PACKAGE",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[strings.ToUpper(ident)]; is_keyword {
		return tok
	}
	return IDENT
}

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// IsKeyword reports whether name is a  keyword, such as "INT" or "ENUM".
func IsKeyword(name string) bool {
	_, ok := keywords[strings.ToUpper(name)]
	return ok
}

// IsIdentifier reports whether name is a Go identifier, that is, a non-empty
// string made up of letters, digits, and underscores, where the first character
// is not a digit. Keywords are not identifiers.
//
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}
