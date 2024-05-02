package model

import (
	"big_market/common"
	"big_market/common/constant"
	"big_market/common/log"
)

var LogicTreeNodeGroup = map[string]TreeNode{
	constant.LockTreeNode:      LockTreeNode{},
	constant.LuckAwardTreeNode: LuckAwardTreeNode{},
	constant.StockNode:         StockTreeNode{},
}

type TreeNode interface {
	Logic(userID string, strategyID int64, awardID int) (code string, strategyAwardVO *StrategyAwardVO)
}

type TreeNodeVO struct {
	TreeID           int
	RuleKey          string
	RuleDesc         string
	RuleValue        string
	TreeNodeLineList []TreeNodeLine
}

type TreeNodeLine struct {
	TreeId               int
	RuleNodeFrom         string
	RuleNodeTo           string
	RuleLimitTypeVO      common.RuleLimitType
	RuleLogicCheckTypeVO string
}

type LockTreeNode struct {
}

type LuckAwardTreeNode struct {
}

type StockTreeNode struct {
}

func (l LockTreeNode) Logic(userID string, strategyID int64, awardID int) (code string, strategyAwardVO *StrategyAwardVO) {
	log.Infof("决策树-积分解锁")
	return constant.Allow, strategyAwardVO
}

func (l LuckAwardTreeNode) Logic(userID string, strategyID int64, awardID int) (code string, strategyAwardVO *StrategyAwardVO) {
	log.Infof("决策树-幸运奖")
	return constant.TakeOver, &StrategyAwardVO{
		AwardId:   101,
		RuleValue: "1,100",
	}
}

func (s StockTreeNode) Logic(userID string, strategyID int64, awardID int) (code string, strategyAwardVO *StrategyAwardVO) {
	log.Infof("决策树-库存")
	return constant.TakeOver, strategyAwardVO
}
