package scanner

import (
	"fmt"
	"github.com/huxiaoyugo/huxykit/hipl/token"
	"testing"
	"time"
)

func TestScanner_Scan(t *testing.T) {

	src := []byte(`// 自定义对象
"RequestParam":{
	"name": {{STRING}},
	"age": {{INT32}},
	"score": {{FLOAT}},
	"is_stu": {{BOOL}},
	"list": [ {{STRING}} ], // 数组字段
	"favorit": {{ENUM<Sports>}}, // 枚举类型的字段
	"ext_info": {{MAP<STRING, STRING>}},
	"school": {{School}}  // 自定义对象类型的字段
}

// 自定义对象
"School": {
	"name":{{STRING}},
	"age":{{INT32}}
}

// enum类型定义
"Sports": {
	"basketball":1,
	"football":2
}
`)
	file := "app.proto"
	s := Scanner{}
	s.Init(src, file, 0, 0)

	tok:= s.Scan()
	for tok.Tok != token.EOF {
		//t.Log(tok.String(), lit)
		tok = s.Scan()
		fmt.Println(tok.Tok.String(), tok.Lit)
	}
	time.Sleep(time.Second)
}
