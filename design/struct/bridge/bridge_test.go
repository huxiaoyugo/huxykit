package bridge

import "testing"

// 将Color 与Shape完全抽离，可以进行任意的组合
func TestBridge(t *testing.T) {

	red := Red{}
	white := White{}
	c := Circle{}
	square := Square{}
	c.SetColor(red)
	square.SetColor(white)

	c.Draw()
	square.Draw()

	c.SetColor(white)
	square.SetColor(red)

	c.Draw()
	square.Draw()

}
