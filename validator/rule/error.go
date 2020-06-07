package rule

import (
	"fmt"
)

type RuleError struct {
	rule Rule
	// 目标规则
	ruleExp string
	// 实际值
	actualVal interface{}
	// 错误信息
	otherErrMsg []interface{}
}

func (e *RuleError) Error() string {
	errMsg := fmt.Sprintf("rulename is %s and ruleExp is '%s' but actual value is %v", e.rule.RuleName(), e.ruleExp, e.actualVal)
	if len(e.otherErrMsg) > 0 {
		return errMsg + ", "+ fmt.Sprint(e.otherErrMsg...)
	}
	return errMsg
}

func Error(rule Rule, ruleExp string, actualVal interface{}, errMsg ...interface{}) *RuleError {
	return &RuleError{
		rule:        rule,
		ruleExp:     ruleExp,
		actualVal:   actualVal,
		otherErrMsg: errMsg,
	}
}
