package model

import (
	"big_market/common"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StrategyRule struct {
	ID         int       `json:"id"`         // 自增ID
	StrategyID int       `json:"strategyId"` // 抽奖策略ID
	AwardID    int       `json:"awardId"`    // 抽奖奖品ID【规则类型为策略，则不需要奖品ID】
	RuleType   int       `json:"ruleType"`   // 抽象规则类型；1-策略规则、2-奖品规则
	RuleModel  string    `json:"ruleModel"`  // 抽奖规则类型【rule_random - 随机值计算、rule_lock - 抽奖几次后解锁、rule_luck_award - 幸运奖(兜底奖品)】
	RuleValue  string    `json:"ruleValue"`  // 抽奖规则比值
	RuleDesc   string    `json:"ruleDesc"`   // 抽奖规则描述
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateTime time.Time `json:"updateTime"` // 更新时间
}

// 数据案例；4000:102,103,104,105 5000:102,103,104,105,106,107 6000:102,103,104,105,106,107,108,109
func (s StrategyRule) GetWeightValues() (map[string][]int, error) {
	result := make(map[string][]int)
	valueGroups := strings.Split(s.RuleValue, common.Space)
	for _, group := range valueGroups {
		if group == "" {
			continue
		}
		parts := strings.Split(group, common.COLON)
		if len(parts) != 2 {
			return nil, errors.New(fmt.Sprintf("%s%s", "rule_weight rule_rule invalid input format", group))
		}
		values := strings.Split(parts[1], common.Split)
		for _, value := range values {
			valueStr, _ := strconv.Atoi(value)
			result[group] = append(result[group], valueStr)
		}
	}
	return result, nil
}
