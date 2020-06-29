package scanner

import (
	"fmt"
	"icode.baidu.com/baidu/duer/protocol-gen/hipl/err"
	"icode.baidu.com/baidu/duer/protocol-gen/hipl/token"
	"unicode"
	"unicode/utf8"
)

const bom = 0xFEFF // byte order mark, only permitted as very first character

type Scanner struct {
	pos     Pos
	tokenBeginPos Pos

	baseRow int

	src []byte
	ch  rune // current character
	//last int
	offset   int
	rdOffset int
}

type Pos struct {
	File     string
	Row, Col int
}

func(p Pos)String() string{
	return fmt.Sprintf("%s:%d:%d: ", p.File, p.Row, p.Col)
}

func(p Pos) GetCol() int {
	return p.Col
}

func(p Pos) GetRow() int {
	return p.Row
}

func(p Pos) GetFile() string {
	return p.File
}

func (s *Scanner) PosString() string {
	return s.tokenBeginPos.String()
}

// Read the next Unicode char into s.ch.
// s.ch < 0 means end-of-file.
//

func (s *Scanner) Init(src []byte, filename string, baseRow int, col int) {
	s.baseRow = baseRow
	s.src = src
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.pos = Pos{
		File: filename,
		Row:  baseRow,
		Col:  col,
	}
	s.tokenBeginPos = s.pos
	s.next()
	if s.ch == bom {
		s.next()
	}
}

func (s *Scanner) FileName()string {
	return s.pos.File
}

func (s *Scanner) Error(msg string) err.Error {
	return err.Err(s.PosString(), msg)
}

func (s *Scanner) ErrorWithPos(pos Pos, msg string) err.Error {
	return err.Err(pos.String(), msg)
}

func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.pos.Row++
			s.pos.Col = 0
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			panic(s.Error("illegal character NUL"))
		case r >= utf8.RuneSelf:
			// todo: 这段不是很懂
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				panic(s.Error("illegal UTF-8 encoding"))
			} else if r == bom && s.offset > 0 {
				panic(s.Error("illegal byte order mark"))
			}
		}
		s.pos.Col+=w
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		s.ch = -1 // eof
	}
}

func (s *Scanner) scanWhitespace() (token.Token, string) {
	offs := s.offset - 1
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
	return token.EMPTY, string(s.src[offs:s.offset])
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_'
}

func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) || s.ch == '.' {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanNumber() string {
	offs := s.offset
	for isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanString() string {
	// '"' opening already consumed
	offs := s.offset - 1
	c := rune(s.src[offs])
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			panic(s.ErrorWithPos(s.tokenBeginPos, "string literal not terminated"))
		}
		s.next()
		if ch == c {
			break
		}
		if ch == '\\' {
			s.scanEscape(c)
		}
	}

	return string(s.src[offs:s.offset])
}

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax err, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.
func (s *Scanner) scanEscape(quote rune) bool {
	var n int
	var base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		panic(s.Error(msg))
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			panic(s.Error(msg))
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		panic(s.Error("escape sequence is invalid Unicode code point"))
		return false
	}

	return true
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}

func (s *Scanner) scanComment() string {
	// initial '/' already consumed; s.ch == '/' || s.ch == '*'
	offs := s.offset - 1 // position of initial '/'

	if s.ch == '/' {
		//-style comment
		// (the final '\n' is not considered part of the comment)
		s.next()
		for s.ch != '\n' && s.ch >= 0 {
			s.next()
		}
		if s.ch == -1 {
			return string(s.src[offs:])
		}
		return string(s.src[offs:s.offset])
	}

	/*-style comment */
	s.next()
	for s.ch >= 0 {
		ch := s.ch
		s.next()
		if ch == '*' && s.ch == '/' {
			s.next()
			return string(s.src[offs:s.offset])
		}
	}
	panic(s.Error("comment not terminated"))
}

func (s *Scanner) Scan() *TokenInfo {
	tokenInfo := &TokenInfo{
		Pos: s.pos,
	}
	s.tokenBeginPos = s.pos
	switch ch := s.ch; {
	case isLetter(ch):
		tokenInfo.Lit = s.scanIdentifier()
		tokenInfo.Tok = token.Lookup(tokenInfo.Lit)

	case isDecimal(ch):
		tokenInfo.Lit = s.scanNumber()
		tokenInfo.Tok = token.INTEGER

	default:
		s.next() // always make progress
		switch ch {
		case -1:
			tokenInfo.Tok = token.EOF
		case ' ', '\n', '\t', '\r':
			tokenInfo.Tok, tokenInfo.Lit = s.scanWhitespace()
		case '"','\'':
			tokenInfo.Tok = token.STR
			tokenInfo.Lit = s.scanString()
		case ',':
			tokenInfo.Tok = token.COMMA
			tokenInfo.Lit = ","
		case '[':
			tokenInfo.Tok = token.LBRACK
			tokenInfo.Lit = "["
		case ']':
			tokenInfo.Tok = token.RBRACK
			tokenInfo.Lit = "]"
		case '{':
			tokenInfo.Tok = token.LBRACE
			tokenInfo.Lit = "{"
		case '}':
			tokenInfo.Tok = token.RBRACE
			tokenInfo.Lit = "}"
		case '<':
			tokenInfo.Tok = token.LSS
			tokenInfo.Lit = "<"
		case '>':
			tokenInfo.Tok = token.GTR
			tokenInfo.Lit = ">"
		case ':':
			tokenInfo.Tok = token.COLON
			tokenInfo.Lit = ":"
		case '/':
			if s.ch == '/' || s.ch == '*' {
				comment := s.scanComment()
				tokenInfo.Tok = token.COMMENT
				tokenInfo.Lit = comment
			} else {
				panic(s.Error("invalid char"))
			}
		default:
			panic(s.Error(fmt.Sprintf("invalid char '%c'", ch)))
		}
	}
	return tokenInfo
}

type TokenInfo struct {
	Tok token.Token
	Pos Pos
	Lit string
}
