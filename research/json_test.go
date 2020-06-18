package research

import (
	"encoding/json"
	"reflect"
	"testing"
)


type tree struct {
	Val int
	Name string
}

func setTreeNameDirect(t *tree) {
	t.Name = "huxiaoyu"
}

func setTreeNameInterface(t interface{})  {
	t.(*tree).Name = "huxiaoyu"
}


func setTreeNameReflect(t interface{}) {
	reflect.ValueOf(t).Elem().FieldByName("Name").SetString("huxiaoyu")
}

func BenchmarkDirect(b *testing.B) {
	t := &tree{}
	for i:=0;i<b.N;i++ {
		setTreeNameDirect(t)
	}
}


func BenchmarkInterface(b *testing.B) {
	t := &tree{}
	for i:=0;i<b.N;i++ {
		setTreeNameInterface(t)
	}
}


func BenchmarkReflect(b *testing.B) {
	t := &tree{}
	for i:=0;i<b.N;i++ {
		setTreeNameReflect(t)
	}
}



type Car struct {
	Name string
	Age int
	Brand string
	Address string
	Level  int
	*Wheel
}

type Wheel struct {
	Name string
	Brand string
	*Material
}

type Material struct {
	Name string
	Type string
}


var c = &Car{
		Name:    "BMW320",
		Age:     100,
		Brand:   "BMW",
		Address: "北京",
		Level:   1,
		Wheel:   &Wheel{
			Name:     "",
			Brand:    "",
			Material: &Material{
				Name: "橡胶",
				Type: "1",
			},
		},
	}

func BenchmarkJson(b *testing.B) {
	for i:=0;i<b.N;i++ {
		json.Marshal(c)
	}
}

func TestJson(t *testing.T) {
	r , err := json.Marshal(c)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(r))
}