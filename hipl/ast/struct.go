package ast

import "github.com/huxiaoyugo/huxykit/hipl/scanner"

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
