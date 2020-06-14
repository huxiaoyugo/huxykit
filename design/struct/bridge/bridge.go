package bridge

import "fmt"

type Shape interface {
	SetColor(color Color)
	Draw()
}

type BaseShape struct {
	color Color
}

func (b *BaseShape) SetColor(color Color) {
	b.color = color
}

type Square struct {
	BaseShape
}

func (s Square) Draw() {
	fmt.Printf("square -- %s\n", s.color.GetName())
}

type Circle struct {
	BaseShape
}

func (s Circle) Draw() {
	fmt.Printf("circle -- %s\n", s.color.GetName())
}

type Color interface {
	GetName() string
}

type Red struct {
}

func (r Red) GetName() string {
	return "red"
}

type White struct {
}

func (w White) GetName() string {
	return "white"
}
