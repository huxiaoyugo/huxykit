package rule

import (
	"testing"
)

func TestError(t *testing.T) {

	msg := Error(GreaterRule{},">=10", 10).Error()
	t.Log(msg)
}
