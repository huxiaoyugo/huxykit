package ast

import "icode.baidu.com/baidu/duer/protocol-gen/hipl/scanner"

type Enum struct {
	Tok    *scanner.TokenInfo
	Name   string
	Fields []*EnumField
}

type EnumField struct {
	Tok    *scanner.TokenInfo
	Name string
	Tag  EnumTag
}

type EnumTag struct {
	Tag int
	Tok *scanner.TokenInfo
}
