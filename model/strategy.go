package model

import (
	"big_market/common"
	"sort"
	"strings"
	"time"
)

// todo 要改json字段名
type Strategy struct {
	ID           int       `json:"id"`
	StrategyID   int       `json:"strategyId"`
	StrategyDesc string    `json:"strategyDesc"`
	RuleModels   string    `json:"ruleModels"` //规则模型，rule配置的模型同步到此表，便于使用
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

func (s Strategy) GetRuleModels() []string {
	if s.RuleModels == "" {
		return nil
	} else {
		rules := strings.Split(s.RuleModels, common.Split)
		sort.Slice(rules, func(i, j int) bool {
			return common.WeightOfRules[rules[i]] < common.WeightOfRules[rules[j]]
		})
		return rules
	}
}
