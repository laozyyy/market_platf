package service

import (
	"big_market/common"
	"big_market/common/log"
	"big_market/model"
	"errors"
)

type TreeEngine struct {
	model.Tree
}

func (t TreeEngine) Process(userID string, strategyID int64, awardID int) (result *model.StrategyAwardVO, err error) {
	nextNodeStr := t.TreeRootRuleNode
	nextNodeVO, ok := t.TreeNodeMap[nextNodeStr]
	if !ok {
		log.Errorf("err: %v", "cannot get tree node")
		return nil, errors.New("cannot get tree node")
	}
	code := ""
	treeNode, ok := LogicTreeNodeGroup[nextNodeVO.RuleKey]
	for ok {
		code, result = treeNode.Logic(userID, strategyID, awardID)
		log.Infof("决策树引擎【%v】treeId:%s node:%s code:%s", t.TreeName, t.TreeID, nextNodeStr, code)
		nextNodeVO = t.nextNode(code, nextNodeVO.TreeNodeLineList)
		treeNode, ok = LogicTreeNodeGroup[nextNodeVO.RuleKey]
	}
	return result, err
}

func (t TreeEngine) nextNode(code string, nodeLines []model.TreeNodeLineVO) (treeNode model.TreeNodeVO) {
	if nodeLines == nil {
		return
	}
	for _, nodeLine := range nodeLines {
		if ok := judge(code, nodeLine); ok {
			return t.TreeNodeMap[nodeLine.RuleNodeTo]
		}
	}
	return
}

func judge(value string, line model.TreeNodeLineVO) bool {
	switch line.RuleLimitType {
	case common.EQUAL:
		return line.RuleLimitValue == value
	case common.GT:
	case common.LT:
	case common.GE:
	case common.LE:
	default:
		return false
	}
	return false
}
