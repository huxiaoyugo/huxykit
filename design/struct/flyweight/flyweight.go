package flyweight

import (
	"sync"
)

type Share struct {
	Name string
}

type Pool struct {
	sync.Pool
}

var (
	sharePool Pool
)

func init() {
	sharePool = Pool{sync.Pool{
		New: func() interface{} {
			return &Share{}
		},
	}}
}

func(p *Pool) GetOneShare() *Share {
	s := sharePool.Get().(*Share)
	s.Name = ""
	return s
}


type Obj1 struct {
	s *Share
}

type Obj2 struct {
	s *Share
}

func NewObj1() *Obj1 {
	return &Obj1{
		s: sharePool.GetOneShare(),
	}
}

func NewObj2() *Obj2 {
	return &Obj2{
		s: sharePool.GetOneShare(),
	}
}



