package research

import (
	"reflect"
	"testing"
)

type Target struct {
	Val int
}

func (t *Target) Run() {
	t.Val++
}

func RunDirect(t *Target) {
	t.Run()
}

func RunByReflect(v interface{}) {
	f := reflect.ValueOf(v)
	f.MethodByName("Run").Call(nil)
}

func BenchmarkRunFuncDirect(b *testing.B) {
	var t = Target{}
	for i:=0;i<b.N;i++ {
		RunDirect(&t)
	}
}

func BenchmarkRunFuncReflect(b *testing.B) {
	var t = Target{}
	for i:=0;i<b.N;i++ {
		RunByReflect(&t)
	}
}
