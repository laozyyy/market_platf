package main

import (
	"big_market/cache"
	"big_market/common"
	"big_market/common/constant"
	log2 "big_market/common/log"
	"big_market/crons"
	"big_market/database"
	"big_market/model"
	"big_market/mq"
	"big_market/service"
	"big_market/service/reposity"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
)

func init() {
	cache.Init()
	database.Init()
	mq.Init()
	crons.AddCron()
}

func TestCreateRaffleActivityOrder(t *testing.T) {
	cart := model.ActivityShopCartEntity{
		UserID: GenerateRandomString(10),
		Sku:    9011,
	}
	order, err := service.CreateRaffleActivityOrder(cart)
	if err != nil {
		log2.Log.Infof("装配失败")
	} else {
		log2.Log.Infof("order:%v", order)
	}
}

func TestStrategyArmory(t *testing.T) {
	int64s := []int64{100001, 100002, 100003, 100004, 100005, 100006}
	for _, i := range int64s {
		err := reposity.AssembleLotteryStrategyWithRules(i)
		if err != nil {
			log2.Log.Infof("装配失败")
		} else {
			log2.Log.Infof("装配成功, strategyID: %d", i)
		}
	}
}

// 权重 （积分根据userid）
func TestPerformRaffle(t *testing.T) {
	success, err := service.PerformRaffle(model.RaffleFactor{
		UserID:     "zym",
		StrategyID: 100001,
	})
	if err != nil {
		log2.Log.Errorf("error: %v", err)
	} else {
		log2.Log.Info("抽奖成功")
		log2.Log.Infof("抽奖结果 %+v", success)
	}

}

// 黑名单
func TestPerformRaffleBlackList(t *testing.T) {
	success, err := service.PerformRaffle(model.RaffleFactor{
		UserID:     "user001",
		StrategyID: 100001,
	})
	if err != nil {
		log2.Log.Errorf("error: %v", err)
	} else {
		log2.Log.Info("抽奖成功")
		log2.Log.Infof("抽奖结果 %+v", success)
	}
}

func TestRedis(t *testing.T) {
	var wg sync.WaitGroup
	num := 10
	wg.Add(num)
	cache.Init()
	//cache.Client.Set(context.Background(), "test", 0, 0)
	for i := num; i > 0; i-- {
		go func() {
			defer func() {
				wg.Done()
			}()
			cache.Client.Decr(context.Background(), "test")
		}()

	}
	wg.Wait()
	log2.Errorf("并发完成")
}

// 解锁(决策树)
// todo err: Exception (501) Reason: \"read tcp [::1]:4862->[::1]:5672: i/o timeout\"
func TestPerformRaffleTree(t *testing.T) {
	crons.AddCron()
	TestStrategyArmory(t)
	var wg sync.WaitGroup
	num := 10
	wg.Add(num)
	success, fail := 0, 0
	for i := num; i > 0; i-- {
		go func() {
			defer func() {
				wg.Done()
			}()
			result, err := service.PerformRaffle(model.RaffleFactor{
				UserID:     "zym",
				StrategyID: 100006,
			})
			if err != nil {
				log2.Log.Errorf("error: %v", err)
				fail++
			} else {
				success++
				log2.Log.Info("抽奖成功")
				log2.Log.Infof("抽奖结果 %+v", result)
			}
		}()

	}
	wg.Wait()
	log2.Errorf("并发抽奖完成, 成功：%d，失败%d，共%d", success, fail, success+fail)
}

func TestUpdate(t *testing.T) {
	err := database.UpdateStrategyAwardAwardCountSurplus(nil, 100001, 101, "80000")
	if err != nil {
		log2.Log.Errorf("error: %v", err)
	}
}

func TestTree(t *testing.T) {
	ruleLock := model.TreeNodeVO{
		TreeID:    "100000001",
		RuleKey:   "rule_lock",
		RuleDesc:  "限定用户已完成N次抽奖后解锁",
		RuleValue: "1",
		TreeNodeLineList: []model.TreeNodeLineVO{
			{
				TreeId:         "100000001",
				RuleNodeFrom:   "rule_lock",
				RuleNodeTo:     "rule_luck_award",
				RuleLimitType:  common.EQUAL,
				RuleLimitValue: constant.TakeOver,
			},
			{
				TreeId:         "100000001",
				RuleNodeFrom:   "rule_lock",
				RuleNodeTo:     "rule_stock",
				RuleLimitType:  common.EQUAL,
				RuleLimitValue: constant.Allow,
			},
		},
	}

	ruleLuckAward := model.TreeNodeVO{
		TreeID:           "100000001",
		RuleKey:          "rule_luck_award",
		RuleDesc:         "限定用户已完成N次抽奖后解锁",
		RuleValue:        "1",
		TreeNodeLineList: nil,
	}

	ruleStock := model.TreeNodeVO{
		TreeID:    "100000001",
		RuleKey:   "rule_stock",
		RuleDesc:  "库存处理规则",
		RuleValue: "",
		TreeNodeLineList: []model.TreeNodeLineVO{
			{
				TreeId:         "100000001",
				RuleNodeFrom:   "rule_lock",
				RuleNodeTo:     "rule_luck_award",
				RuleLimitType:  common.EQUAL,
				RuleLimitValue: constant.TakeOver,
			},
		},
	}

	ruleTreeVO := model.Tree{
		TreeID:           "100000001",
		TreeName:         "决策树规则；增加dall-e-3画图模型",
		TreeDesc:         "决策树规则；增加dall-e-3画图模型",
		TreeRootRuleNode: "rule_lock",
		TreeNodeMap: map[string]model.TreeNodeVO{
			"rule_lock":       ruleLock,
			"rule_stock":      ruleStock,
			"rule_luck_award": ruleLuckAward,
		},
	}

	engine := service.TreeEngine{ruleTreeVO}
	result, err := engine.Process("zym", 100001, 100)

	if err != nil {
		log2.Log.Errorf("err: %v", err)
		return
	}

	// 将对象转为JSON字符串
	jsonStr, err := json.Marshal(result)
	if err != nil {
		log2.Log.Errorf("err: %v", err)
		return
	}

	log2.Log.Infof("测试结果：%v", string(jsonStr))
}

/*
*
从装配的策略中随机获取奖品ID值
*/
func TestGetAssembleRandomVal(t *testing.T) {
	i := 50
	for i > 0 {
		i--
		result, _ := reposity.GetRandomAwardIdByWeight("100001", "4000:102,103,104,105")
		log.Printf("测试结果：%v - 奖品ID值\n", result)
	}

}

//还有with rules的情况

func TestShuffle(t *testing.T) {
	strategyAwardSearchRateTable := make(map[int]int)
	// 添加内容到Map中
	strategyAwardSearchRateTable[1] = 10
	strategyAwardSearchRateTable[2] = 20
	strategyAwardSearchRateTable[3] = 30
	strategyAwardSearchRateTable[4] = 40

	// 将Map中的值转换为List
	valueList := make([]int, 0, len(strategyAwardSearchRateTable))
	for _, value := range strategyAwardSearchRateTable {
		valueList = append(valueList, value)
	}

	// 使用rand.Shuffle()方法对值的List进行乱序
	rand.Shuffle(len(valueList), func(i, j int) {
		valueList[i], valueList[j] = valueList[j], valueList[i]
	})

	// 将乱序后的值重新放回Map中
	randomizedMap := make(map[int]int)
	index := 0
	for key := range strategyAwardSearchRateTable {
		randomizedMap[key] = valueList[index]
		index++
	}

	// 打印乱序后的Map内容
	for key, value := range randomizedMap {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
	}
}

func GenerateRandomString(length int) string {
	// 计算生成字节数
	byteLength := length * 3 / 4

	// 生成随机字节
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	// 将随机字节编码为字符串
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	// 截取指定长度的随机字符串
	randomString = randomString[:length]

	return randomString
}
