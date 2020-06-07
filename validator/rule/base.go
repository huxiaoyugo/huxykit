package rule

import (
	"fmt"
	"regexp"
)

type BaseRule struct {
}

func (n BaseRule) isMatchRule(pattern string, ruleStr string) bool {
	ok, err := regexp.Match(pattern, []byte(ruleStr))
	if err != nil {
		fmt.Println(err)
	}
	return ok
}
