package adapter

// client 需要我们提供的方法接口
type Request interface {
	Request()
}

// 被适配的对象
type Adaptee struct {
}

// 被适配的对象的方法名不符合client的需求
func(Adaptee) SpecialRequest() {
}

// 适配器封装Adaptee
type Adapter struct {
	Adaptee
}

// 实现client需要的方法，实际调用Adaptee的方法
func (a Adapter) Request() {
	a.SpecialRequest()
}