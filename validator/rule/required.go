package rule

import (
	"reflect"
)

type RequiredRule struct {
	BaseRule
}

const (
	requiredRuleName = "required"
	patternRequired  = "required"
)

func (r RequiredRule) RuleName() string {
	return requiredRuleName
}

func (r RequiredRule) Check(val interface{}, ruleExp string) (bool, error) {
	if !r.isMatchRule(patternRequired, ruleExp) {
		return false, nil
	}
	if reflect.ValueOf(val).IsZero() {
		return true, Error(r, ruleExp, val)
	}
	return true, nil
}

func init() {
	registerDefaultRule(&RequiredRule{})
}
