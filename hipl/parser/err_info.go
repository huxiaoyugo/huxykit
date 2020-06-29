package parser

import "icode.baidu.com/baidu/duer/protocol-gen/hipl/err"

func (p *Parser) errExpect(msg string) {
	if p.curToken != nil {
		panic(err.ErrF("%s Expected '%s'", p.curToken.Pos.String(), msg))
	} else {
		panic(err.ErrF("Expected '%s'", msg))
	}
}

func (p *Parser) errExpectBut(expect, got string) {
	if p.curToken != nil {
		panic(err.ErrF("%s Expected '%s', but got '%s'",p.curToken.Pos.String(), expect, got))
	} else {
		panic(err.ErrF("Expected '%s', but got '%s'", expect, got))
	}
}

func(p *Parser) errorInfo(msg string) {
	if p.curToken != nil {
		panic(err.ErrF("%s%s", p.curToken.Pos.String(), msg))
	} else {
		panic(err.Err(msg))
	}
}
