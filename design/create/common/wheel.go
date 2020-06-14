package common

import "fmt"

type Wheel interface {
	GetName() string
}

type BaseWheel struct {
	Name string
}

type AWheel struct {
	BaseWheel
}

func (a AWheel) GetName() string {
	return fmt.Sprintf("[wheel_brand: A]-%s", a.Name)
}

type BWheel struct {
	BaseWheel
}

func (b BWheel) GetName() string {
	return fmt.Sprintf("[wheel_brand: B]-%s", b.Name)
}
