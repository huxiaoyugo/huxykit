package rule

type NumberCompare struct {
	BaseRule
}

func (n NumberCompare) RuleName(name string) string {
	return name
}

func (n NumberCompare) getRuleNum(ruleExp string) int64 {
	num := int64(0)
	isPositive := true
	for i := 0; i < len(ruleExp); i++ {
		if ruleExp[i] == '-' {
			isPositive = false
		} else if ruleExp[i] >= '0' && ruleExp[i] <= '9' {
			num *= 10
			num = int64(ruleExp[i] - '0')
		}
	}
	if !isPositive {
		return -num
	}
	return num
}
