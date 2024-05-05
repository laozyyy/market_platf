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
	nodeVO, ok := t.TreeNodeMap[nextNodeStr]
	if !ok {
		log.Errorf("err: %v", "cannot get tree node")
		return nil, errors.New("cannot get tree node")
	}
	code := ""
	treeNode, ok := LogicTreeNodeGroup[nodeVO.RuleKey]
	for ok {
		code, result, err = treeNode.Logic(userID, strategyID, awardID, nodeVO.RuleValue)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		log.Infof("决策树引擎【%v】treeId:%s node:%s code:%s", t.TreeName, t.TreeID, nextNodeStr, code)
		nodeVO = t.nextNode(code, nodeVO.TreeNodeLineList)
		treeNode, ok = LogicTreeNodeGroup[nodeVO.RuleKey]
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
