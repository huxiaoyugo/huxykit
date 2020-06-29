package parser

import (
	"fmt"
	"testing"
)



func TestParser_Start(t *testing.T) {
	//
	//file := "app.proto"
	//s := &scanner.Scanner{}
	//s.Init(src, file, 9, 0)

	p := Parser{
		protoPath:[]string{
			"/usr/local/include",
			"./",
		},
		//scan: []*scanner.Scanner{s},
	}

	err := p.ParseFile("test.md")
	//err := p.Start()
	if err != nil {
		t.Error(err)
	}
}



type I interface {

}
func TestInterface(t *testing.T) {
	var i I
	SetNum(i)
	fmt.Print(i)
}
func SetNum(i I) {
	i = "sss"
}