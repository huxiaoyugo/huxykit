/*
造者模式(Builder Pattern)：将一个复杂对象的构建与它的表示分离，使得同样的构建过程可以创建不同的表示。
建造者模式是一步一步创建一个复杂的对象，它允许用户只通过指定复杂对象的类型和内容就可以构建它们，用户不
需要知道内部的具体构建细节。建造者模式属于对象创建型模式。根据中文翻译的不同，建造者模式又可以称为生成
器模式。

参考资料：
https://design-patterns.readthedocs.io/zh_CN/latest/creational_patterns/builder.html
 */
package builder

import "fmt"

type Builder interface {
	Part1()
	Part2()
	Part3()
}

type Director struct {
	builder Builder
}

func CreateDirector(b Builder) *Director {
	d := &Director{builder: b}
	d.builder.Part1()
	d.builder.Part2()
	d.builder.Part3()
	return d
}

type BuilderA struct {
}

func (BuilderA) Part1() {
	fmt.Println("==part1==")
}

func (BuilderA) Part2() {
	fmt.Println("==part2==")
}

func (BuilderA) Part3() {
	fmt.Println("==part3==")
}

type BuilderB struct {
}

func (BuilderB) Part1() {
	fmt.Println("**part1**")
}

func (BuilderB) Part2() {
	fmt.Println("**part2**")
}

func (BuilderB) Part3() {
	fmt.Println("**part3**")
}
