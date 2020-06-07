package validator

import "testing"

type Student struct {
	Name string `json:"name" validator:"required"`
	Age  int    `json:"age" validator:">=0, gt1, , , gt 00020"`
}

func TestValidator_Valid(t *testing.T) {

	var target interface{}
	target = struct {
		Good string `validator:"required"`
	}{
		Good:"adfd",
	}
	valid := NewValidator()
	err := valid.IsValid(target)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("target is valid")
}
