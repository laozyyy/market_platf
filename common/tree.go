package common

type RuleLimitType int
type RuleLimitValue int

const (
	EQUAL RuleLimitType = iota + 1
	GT
	LT
	GE
	LE
	ENUM
	TreeAllowRule RuleLimitValue = iota
	TreeTakeOver
)

var RuleLimitTypeStrings = map[RuleLimitType]string{
	EQUAL: "等于",
	GT:    "大于",
	LT:    "小于",
	GE:    "大于&等于",
	LE:    "小于&等于",
	ENUM:  "枚举",
}

//var RuleLimitTypeStrings = map[RuleLimitValue]string{
//	TreeAllowRule: "放行；执行后续的流程，不受规则引擎影响",
//	TreeTakeOver:  "接管；后续的流程，受规则引擎执行结果影响",
//}
