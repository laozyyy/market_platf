package reposity

import (
	"big_market/cache"
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

func GetTreeByTreeID(treeID string) (*model.Tree, error) {
	treeStr, err := cache.Client.Get(context.Background(), constant.RuleTreeVOKey+treeID).Result()
	// 意外错误
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("err %v", err)
		return nil, err
	}
	// redis存在
	if !errors.Is(err, redis.Nil) {
		var tree model.Tree
		_ = json.Unmarshal([]byte(treeStr), &tree)
		return &tree, nil
	}
	tree, err := assembleTreeByTreeID(treeID)
	if err != nil {
		log.Errorf("err %v", err)
		return nil, err
	}
	marshal, err := json.Marshal(tree)
	err = cache.Client.Set(context.Background(), constant.RuleTreeVOKey+treeID, marshal, 0).Err()
	if err != nil {
		log.Errorf("err %v", err)
		return nil, err
	}
	return tree, nil
}

func assembleTreeByTreeID(treeID string) (*model.Tree, error) {
	rawTree, err := database.QueryRuleTreeByTreeID(nil, treeID)
	if err != nil {
		log.Errorf("err %v", err)
		return nil, err
	}
	treeNodes, err := database.QueryRuleTreeNodesByTreeID(nil, treeID)
	treeNodeLines, err := database.QueryRuleTreeNodeLinesByTreeID(nil, treeID)
	ruleKeyToTreeNodeLine := make(map[string][]model.TreeNodeLineVO)
	treeNodeMap := make(map[string]model.TreeNodeVO)
	for _, line := range treeNodeLines {
		treeNodeLineVO := model.TreeNodeLineVO{
			TreeId:         treeID,
			RuleNodeFrom:   line.RuleNodeFrom,
			RuleNodeTo:     line.RuleNodeTo,
			RuleLimitType:  line.RuleLimitType,
			RuleLimitValue: line.RuleLimitValue,
		}
		ruleKeyToTreeNodeLine[line.RuleNodeFrom] = append(ruleKeyToTreeNodeLine[line.RuleNodeFrom], treeNodeLineVO)
	}
	for _, treeNode := range treeNodes {
		treeNodeVO := model.TreeNodeVO{
			TreeID:           treeID,
			RuleKey:          treeNode.RuleKey,
			RuleDesc:         treeNode.RuleDesc,
			RuleValue:        treeNode.RuleValue,
			TreeNodeLineList: ruleKeyToTreeNodeLine[treeNode.RuleKey],
		}
		treeNodeMap[treeNodeVO.RuleKey] = treeNodeVO
	}
	tree := model.Tree{
		TreeID:           treeID,
		TreeName:         rawTree.TreeName,
		TreeDesc:         rawTree.TreeName,
		TreeRootRuleNode: rawTree.TreeNodeRuleKey,
		TreeNodeMap:      treeNodeMap,
	}
	return &tree, nil
}
