package ast

import (
	"icode.baidu.com/baidu/duer/protocol-gen/hipl/scanner"
	"regexp"
)

/*
"request": {
    "url": "dueros://account/GetUserInfo",
    "params": {
		"name": {{STRING}},
		"extra_info": {{RequestParam}}
	}
}
*/

type Request struct {
	Url RequestUrl

	// struct name
	Params *RequestParams
}

type RequestParams struct {
	Tok        *scanner.TokenInfo
	StructName string
}


type RequestUrl struct {
	Tok *scanner.TokenInfo
	Name string
}


func(r *Request) GeneParamsName() string{

	// todo
	s := regexp.MustCompile("//.+$").FindString(r.Url.Name)

	return s
}