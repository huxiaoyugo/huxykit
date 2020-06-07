package rule

import (
	"fmt"
	"reflect"
)

const (
	ruleNameGt = "greater"
	// 正则匹配规则
	patternGt = `(>|(gt)) {0,}(\+|-)?[0-9]+`
)

type GreaterRule struct {
	NumberCompare
}

func (g GreaterRule) RuleName() string {
	return g.NumberCompare.RuleName(ruleNameGt)
}

func (g GreaterRule) Check(val interface{}, ruleExp string) (bool, error) {

	if !g.isMatchRule(patternGt, ruleExp) {
		return false, nil
	}
	v := reflect.ValueOf(val)
	num := g.getRuleNum(ruleExp)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() <= num {
			return true, Error(g, ruleExp, val)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		// 如果为负数，必然满足
		if num < 0 {
			break
		}
		// 如果为正数, 则可以直接转成uint64进行比较
		if v.Uint() <= uint64(num) {
			return true, Error(g, ruleExp, val)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() <= float64(num) {
			return true, Error(g, ruleExp, val)
		}
	default:
		return true, Error(g, ruleExp, val, fmt.Sprintf("%s can not match this rule", v.Kind()))
	}
	return true, nil
}

func init() {
	registerDefaultRule(&GreaterRule{})
}
