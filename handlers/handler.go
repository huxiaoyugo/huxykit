package handlers

import (
	"context"
	"fmt"
)

type EchoContext interface {
	context.Context
	Set(interface{}, interface{})
	Get(interface{}) (interface{}, bool)
	PrintAll()
}

type HandlerFunc func(context EchoContext) error

func OneHandler(handlerFunc HandlerFunc) HandlerFunc {
	return func(c EchoContext) error {
		c.Set("one", 1)
		fmt.Println("==1==")
		return handlerFunc(c)
	}
}

func TwoHandler(handlerFunc HandlerFunc) HandlerFunc {
	return func(c EchoContext) error {
		c.Set("two", 2)
		fmt.Println("==2==")
		return handlerFunc(c)
	}
}

func ThreeHandler(handlerFunc HandlerFunc) HandlerFunc {
	return func(c EchoContext) error {
		c.Set("three", 3)
		fmt.Println("==3==")
		return handlerFunc(c)
	}
}

func DoneHandler(handlerFunc HandlerFunc) HandlerFunc {
	return func(c EchoContext) error {
		err := handlerFunc(c)
		if err != nil {
			return err
		}
		fmt.Println("==done==")
		c.PrintAll()
		return nil
	}
}

func Default(context EchoContext) error {
	context.Set("default", 0)
	return nil
}
