package decorator

import (
	"fmt"
	"testing"
)

func TestDecorator(t *testing.T) {

	add := NewAddOperator(DefaultOperator{}, 9)
	mul := NewMulOperator(add, 100)

	fmt.Println(mul.Cale())

}
