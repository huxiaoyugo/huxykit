/*
简单工厂模式(Simple Factory Pattern)：
又称为静态工厂方法(Static Factory Method)模式，它属于类创建型模式。在简单工厂模式中，
可以根据参数的不同返回不同类的实例。简单工厂模式专门定义一个类来负责创建其他类的实例，
被创建的实例通常都具有共同的父类。

参考资料：
https://design-patterns.readthedocs.io/zh_CN/latest/creational_patterns/simple_factory.html
 */
package simple_factory

import (
	"fmt"
	"github.com/huxiaoyugo/huxykit/design/create/common"
	. "github.com/huxiaoyugo/huxykit/design/create/common/common"
)


type CarSimpleFactory struct {
}

func (cf CarSimpleFactory) CreateCar(brand common.Brand, name string) (common.ICar, error) {
	var car common.ICar
	switch brand {
	case common.BrandBYD:
		car = common.NewBYD(name, common.BYDPartFactory{})
	case common.BrandBMW:
		car = common.NewBMW(name, common.BMWPartFactory{})
	case common.BrandBenchi:
		car = common.NewBenChi(name, common.BenChiPartFactory{})
	}
	if car == nil {
		return nil, fmt.Errorf("brand not found")
	}
	return car, nil
}
