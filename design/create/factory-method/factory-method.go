/*
工厂方法模式(Factory Method Pattern)又称为工厂模式，也叫虚拟构造器(Virtual Constructor)模式
或者多态工厂(Polymorphic Factory)模式，它属于类创建型模式。在工厂方法模式中，工厂父类负责定义创
建产品对象的公共接口，而工厂子类则负责生成具体的产品对象，这样做的目的是将产品类的实例化操作延迟到
工厂子类中完成，即通过工厂子类来确定究竟应该实例化哪一个具体产品类。
参考资料：
https://design-patterns.readthedocs.io/zh_CN/latest/creational_patterns/factory_method.html#
 */
package factory_method

import . "github.com/huxiaoyugo/huxykit/design/create/common"

type CarFactory interface {
	Create(name string) ICar
}

type BMWFactory struct {
}

func (BMWFactory) Create(name string) ICar {
	return NewBMW(name, BMWPartFactory{})
}

type BYDFactory struct {
}

func (BYDFactory) Create(name string) ICar {
	return NewBYD(name, BYDPartFactory{})
}

type BenChiFactory struct {
}

func (BenChiFactory) Create(name string) ICar {
	return NewBenChi(name, BenChiPartFactory{})
}
