package ast

import (
	"icode.baidu.com/baidu/duer/protocol-gen/hipl/scanner"
	"icode.baidu.com/baidu/duer/protocol-gen/hipl/token"
	"strings"
)

type IType int

const (
	basic_type_beg IType = iota
	STRING
	INT32
	INT64
	UINT32
	UINT64
	FLOAT
	DOUBLE
	BOOL
	basic_type_end
	MAP
	ARRAY
	STRUCT
	ENUM
)

var iTypes = [...]string{
	STRING: "string",
	INT32:  "int32",
	INT64:  "int64",
	UINT32: "uint32",
	UINT64: "uint64",
	FLOAT:  "float",
	DOUBLE: "double",
	BOOL:   "bool",
	MAP:    "map",
	ARRAY:  "[]",
	STRUCT: "struct",
	ENUM:   "enum",
}

type Language int

const (
	Hipl Language = iota
	Proto3
	Go
)

type Type interface {
	GetType() IType
	ToString(language Language) string
}

func (i IType) String() string {
	return iTypes[i]
}

type Basic struct {
	Tok *scanner.TokenInfo
	Typ IType
}

func (b Basic) ToString(language Language) string {
	switch language {
	case Hipl:
		return strings.ToUpper(b.Typ.String())
	default:
		return strings.ToUpper(b.Typ.String())
	}
}

func (b Basic) GetType() IType {
	return b.Typ
}

var _ Type = Basic{}

type Map struct {
	Basic
	Key, Val Type
}

type Array struct {
	Basic
	ItemType Type
}

type StructType struct {
	Basic
	Name string
}

type EnumType struct {
	Basic
	ValType Type
}

var _ Type = StructType{}


func NewBasicType(tokenInfo *scanner.TokenInfo)Type {
	b := &Basic{
		Tok: tokenInfo,
	}
	switch tokenInfo.Tok {
	case token.STRING:
		b.Typ = STRING
	case token.INT, token.INT32:
		b.Typ = INT32
	case token.UINT32:
		b.Typ = UINT32
	case token.INT64:
		b.Typ = INT64
	case token.UINT64:
		b.Typ = UINT64
	case token.BOOLEAN, token.BOOL:
		b.Typ = BOOL
	case token.DOUBLE:
		b.Typ = DOUBLE
	case token.FLOAT:
		b.Typ = FLOAT
	default:
		panic("invalid type")
	}
	return b
}

func NewMapType(tokenInfo *scanner.TokenInfo) *Map{
	return &Map{
		Basic: Basic{
			Tok: tokenInfo,
			Typ: MAP,
		},
	}
}


func NewEnumType(tokenInfo *scanner.TokenInfo) *EnumType{
	return &EnumType{
		Basic: Basic{
			Tok: tokenInfo,
			Typ: ENUM,
		},
	}
}


func NewStructType(tokenInfo *scanner.TokenInfo) *StructType{
	return &StructType{
		Basic: Basic{
			Tok: tokenInfo,
			Typ: STRUCT,
		},
		Name: tokenInfo.Lit,
	}
}

func NewArrayType(tokenInfo *scanner.TokenInfo) *Array{
	return &Array{
		Basic: Basic{
			Tok: tokenInfo,
			Typ: ARRAY,
		},
	}
}