package factory_method

import "testing"

func TestMethodFactory(t *testing.T) {
	Run(BMWFactory{}, "320")
	Run(BYDFactory{}, "宋Pro")
	Run(BenChiFactory{}, "320")
}

func Run(factory CarFactory, name string) {
	factory.Create(name).Run()
}
