package service

import (
	"big_market/common"
	"big_market/common/log"
	"big_market/model"
	"errors"
)

type Tree struct {
	TreeID           int                         `json:"treeId"`           // 规则树ID
	TreeName         string                      `json:"treeName"`         // 规则树名称
	TreeDesc         string                      `json:"treeDesc"`         // 规则树描述
	TreeRootRuleNode string                      `json:"treeRootRuleNode"` // 规则根节点
	TreeNodeMap      map[string]model.TreeNodeVO `json:"treeNodeMap"`      // 规则节点
}

func (t Tree) Process(userID string, strategyID int64, awardID int) (result *model.StrategyAwardVO, err error) {
	nextNodeStr := t.TreeRootRuleNode
	nextNodeVO, ok := t.TreeNodeMap[nextNodeStr]
	if !ok {
		log.Errorf("err: %v", "cannot get tree node")
		return nil, errors.New("cannot get tree node")
	}
	code := ""
	treeNode, ok := model.LogicTreeNodeGroup[nextNodeVO.RuleKey]
	for ok {
		code, result = treeNode.Logic(userID, strategyID, awardID)
		log.Infof("决策树引擎【%v】treeId:%d node:%s code:%s", t.TreeName, t.TreeID, nextNodeStr, code)
		nextNodeVO = t.nextNode(code, nextNodeVO.TreeNodeLineList)
		treeNode, ok = model.LogicTreeNodeGroup[nextNodeVO.RuleKey]
	}
	return result, err
}

func (t Tree) nextNode(code string, nodeLines []model.TreeNodeLine) (treeNode model.TreeNodeVO) {
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

func judge(value string, line model.TreeNodeLine) bool {
	switch line.RuleLimitTypeVO {
	case common.EQUAL:
		return line.RuleLogicCheckTypeVO == value
	case common.GT:
	case common.LT:
	case common.GE:
	case common.LE:
	default:
		return false
	}
	return false
}
