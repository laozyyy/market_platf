package model

import (
	"time"
)

type Tree struct {
	TreeID           string                `json:"treeId"`           // 规则树ID
	TreeName         string                `json:"treeName"`         // 规则树名称
	TreeDesc         string                `json:"treeDesc"`         // 规则树描述
	TreeRootRuleNode string                `json:"treeRootRuleNode"` // 规则根节点
	TreeNodeMap      map[string]TreeNodeVO `json:"treeNodeMap"`      // 规则节点
}

type TreeNodeVO struct {
	TreeID           string
	RuleKey          string
	RuleDesc         string
	RuleValue        string
	TreeNodeLineList []TreeNodeLineVO
}

type TreeNodeLineVO struct {
	TreeId         string
	RuleNodeFrom   string
	RuleNodeTo     string
	RuleLimitType  string
	RuleLimitValue string
}

// RuleTree 表的结构体定义
type RuleTree struct {
	ID              uint64    `json:"id"`                  // 自增ID
	TreeID          string    `json:"tree_id"`             // 规则树ID
	TreeName        string    `json:"tree_name"`           // 规则树名称
	TreeDesc        string    `json:"tree_desc,omitempty"` // 规则树描述
	TreeNodeRuleKey string    `json:"tree_node_rule_key"`  // 规则树根入口规则
	CreateTime      time.Time `json:"create_time"`         // 创建时间
	UpdateTime      time.Time `json:"update_time"`         // 更新时间
}

// RuleTreeNode 表的结构体定义
type RuleTreeNode struct {
	ID         uint64    `json:"id"`                   // 自增ID
	TreeID     string    `json:"tree_id"`              // 规则树ID
	RuleKey    string    `json:"rule_key"`             // 规则Key
	RuleDesc   string    `json:"rule_desc"`            // 规则描述
	RuleValue  string    `json:"rule_value,omitempty"` // 规则比值
	CreateTime time.Time `json:"create_time"`          // 创建时间
	UpdateTime time.Time `json:"update_time"`          // 更新时间
}

// RuleTreeNodeLine 表的结构体定义
type RuleTreeNodeLine struct {
	ID             uint64    `json:"id"`               // 自增ID
	TreeID         string    `json:"tree_id"`          // 规则树ID
	RuleNodeFrom   string    `json:"rule_node_from"`   // 规则Key节点 From
	RuleNodeTo     string    `json:"rule_node_to"`     // 规则Key节点 To
	RuleLimitType  string    `json:"rule_limit_type"`  // 限定类型；1:=;2:>;3:<;4:>=;5<=;6:enum[枚举范围];
	RuleLimitValue string    `json:"rule_limit_value"` // 限定值（到下个节点）
	CreateTime     time.Time `json:"create_time"`      // 创建时间
	UpdateTime     time.Time `json:"update_time"`      // 更新时间
}
