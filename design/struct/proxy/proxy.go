package proxy



type Subject interface {
	Request() string
}

type RealSub struct {
}

func(RealSub) Request() string{
	return "real"
}

type Proxy struct {
	real RealSub
}

func(p Proxy) Request() string {
	// before request
	var res string
	res += "before "

	res += p.real.Request()

	// after request
	res += " after"
	return res
}