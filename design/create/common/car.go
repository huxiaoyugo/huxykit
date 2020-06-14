package common

import (
	"fmt"

)

type ICar interface {
	Run()
	GetName() string
	GetBrand() Brand

	SetChar(char Char)
	SetWheel(wheel Wheel)
}

type car struct {
	name  string
	brand Brand
	char  Char
	wheel Wheel
}

func (c *car) SetChar(char Char) {
	c.char = char
}

func (c *car) SetWheel(wheel Wheel) {
	c.wheel = wheel
}

func (c *car) GetBrand() Brand {
	return c.brand
}

func (c *car) Run() {
	fmt.Printf("[brand:%s][char:%s][wheel:%s] %s is running\n", c.GetBrand().ToString(), GetName(), GetName(), c.GetName())
}

func (c *car) GetName() string {
	return c.name
}

func newCar(name string, brand Brand) *car {
	return &car{
		name:  name,
		brand: brand,
	}
}

type BMW struct {
	car
}

func (bmw *BMW) Run() {
	fmt.Printf("***[%s]*** %s is running\n", bmw.GetBrand().ToString(), bmw.GetName())
}

type BYD struct {
	*car
}

type BenChi struct {
	*car
}

func NewBMW(name string, factory PartFactory) ICar {
	c := newCar(name, BrandBMW)
	SetPart(c, factory)
	return c
}

func NewBYD(name string, factory PartFactory) ICar {
	c := newCar(name, BrandBYD)
	SetPart(c, factory)
	return c
}

func NewBenChi(name string, factory PartFactory) ICar {
	c := newCar(name, BrandBenchi)
	SetPart(c, factory)
	return c
}

func SetPart(iCar ICar, factory PartFactory) {
	iCar.SetChar(CreateChar())
	iCar.SetWheel(CreateWheel())
}
