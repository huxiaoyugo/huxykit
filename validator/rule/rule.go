package rule

type Rule interface {
	// 规则名称，仅用于表示规则名称，无实际意义
	RuleName() string

	/*
	* 检查是否合法
	*
	* PARAM:
	* val 待检测的值
	* ruleExp 规则表达式
	*
	* RETURN：
	* bool 是否匹配该规则，如果不匹配，则忽略
	* error 在匹配该规则的前提下，检查参数是否符合改规则，如果符合则返回nil,否则返回对应的错误信息
	 */
	Check(val interface{}, ruleExp string) (bool, error)
}

var (
	DefaultRules []Rule
)

func registerDefaultRule(r Rule) {
	DefaultRules = append(DefaultRules, r)
}
