package simple_factory

import (
	"github.com/huxiaoyugo/huxykit/design/create/common"
	"testing"
)

func TestCarSimpleFactory_CreateCar(t *testing.T) {
	factory := CarSimpleFactory{}
	bmw, _ := factory.CreateCar(common.BrandBMW, "320")
	byd, _ := factory.CreateCar(common.BrandBYD, "宋")
	ben, _ := factory.CreateCar(common.BrandBenchi, "C级")
	bmw.Run()
	byd.Run()
	ben.Run()
}
