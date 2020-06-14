package decorator

type Operator interface {
	Cale() int
}

type DefaultOperator struct {
}

func (DefaultOperator) Cale() int {
	return 0
}

type MulOperator struct {
	op  Operator
	num int
}

func (m MulOperator) Cale() int {
	return m.op.Cale() * m.num
}

func NewMulOperator(op Operator, num int) Operator {
	return &MulOperator{
		op:  op,
		num: num,
	}
}

type AddOperator struct {
	op Operator
	num int
}

func (a AddOperator) Cale() int {
	return a.op.Cale() + a.num
}

func NewAddOperator(op Operator, num int) Operator {
	return &AddOperator{
		op:  op,
		num: num,
	}
}

