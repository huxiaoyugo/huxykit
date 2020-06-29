package ast

import (
	"icode.baidu.com/baidu/duer/protocol-gen/hipl/scanner"
	"strings"
)

/*
// resource声明
"resource": {
    "header": {
        "namespace": "ai.dueros.resource.account",
        "name": "PhoneInfo",
        "version": "v1.0"
    },
    "payload": {
        "phone": {{STRING}}
    }
}
*/
type Resource struct {
	Header ResourceHeader
	Payload *ResourcePayload
}

type ResourceHeader struct {
	Namespace string
	Name      string
	Version   string
	Tok       *scanner.TokenInfo
}

type ResourcePayload struct {
	Tok *scanner.TokenInfo
	// payload struct name
	StructName string
}

func (r *Resource) Key() string {
	return r.Header.Namespace + "_" + r.Header.Name
}

func (r *Resource) GenePayloadStructName() string {
	return strings.Replace(r.Header.Namespace, ".", "_ ", -1) + "_" + r.Header.Name
}

