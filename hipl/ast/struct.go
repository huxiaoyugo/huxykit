package ast

import "icode.baidu.com/baidu/duer/protocol-gen/hipl/scanner"

type Struct struct {
	Tok    *scanner.TokenInfo
	Name   string
	Fields []*StructField
}

type StructField struct {
	Tok  *scanner.TokenInfo
	Name string
	Type Type
}
