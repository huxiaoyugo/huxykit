package abstrict_factory

import (
	"github.com/huxiaoyugo/huxykit/design/create/common"
	"testing"
)

func TestDescPartFactory(t *testing.T) {

	DescPartFactory(common.BMWPartFactory{})
	DescPartFactory(common.BYDPartFactory{})
	DescPartFactory(common.BenChiPartFactory{})
}
