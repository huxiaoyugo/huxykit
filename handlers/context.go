package handlers

import (
	"context"
	"fmt"
	"sync"
)

type eContext struct {
	context.Context
	sync.Map
}

func (e *eContext) Set(k interface{},v interface{}) {
	e.Map.Store(k, v)
}

func (e *eContext) Get(k interface{}) (interface{}, bool) {
	return e.Map.Load(k)
}

func (e *eContext) PrintAll() {
	e.Map.Range(func(key, value interface{}) bool {
		fmt.Printf("%v: %v\n", key, value)
		return true
	})
}

func NewEchoContext() EchoContext {
	return &eContext{
		Context: context.Background(),
		Map:     sync.Map{},
	}
}

var _ EchoContext = &eContext{}

