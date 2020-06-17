package research

import (
	"fmt"
	"testing"
)

//
// 堆栈内存分配性能测试
// ================================================

type Tree struct {
	Val int
	N1 Node
	N2 Node
}

type Node struct {
	Name string
}

// 不返回指针
func NewTarget() Tree{
	return Tree{

		N1:  Node{},
		N2:  Node{},
	}
}

// 返回指针
func NewTargetPointer() *Tree{
	return &Tree{
		N1:  Node{},
		N2:  Node{},
	}
}

func BenchmarkNewTarget(b *testing.B) {
	var t Tree
	for i:=0;i<b.N;i++ {
		t = NewTarget()
	}
	fmt.Print(t.Val)
}

func BenchmarkNewTargetPointer(b *testing.B) {
	var t *Tree
	for i:=0;i<b.N;i++ {
		t = NewTargetPointer()
	}
	fmt.Print(t.Val)
}
