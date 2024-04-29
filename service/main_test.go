package service

import (
	"big_market/common"
	"big_market/model"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestStrategyArmory(t *testing.T) {
	int64s := []int64{100001, 100002, 100003}
	for _, i := range int64s {
		err := AssembleLotteryStrategyWithRules(i)
		if err != nil {
			common.Log.Infof("装配失败")
		} else {
			common.Log.Infof("装配成功, strategyID: %d", i)
		}
	}
}

// 权重 （积分根据userid）
func TestPerformRaffle(t *testing.T) {
	success, err := PerformRaffle(model.RaffleFactor{
		UserID:     "zym",
		StrategyID: 100001,
	})
	if err != nil {
		common.Log.Errorf("error: %v", err)
	} else {
		common.Log.Info("抽奖成功")
		common.Log.Infof("抽奖结果 %+v", success)
	}

}

// 黑名单
func TestPerformRaffleBlackList(t *testing.T) {
	success, err := PerformRaffle(model.RaffleFactor{
		UserID:     "user001",
		StrategyID: 100001,
	})
	if err != nil {
		common.Log.Errorf("error: %v", err)
	} else {
		common.Log.Info("抽奖成功")
		common.Log.Infof("抽奖结果 %+v", success)
	}
}

// 解锁
func TestPerformRaffleLock(t *testing.T) {
	success, err := PerformRaffle(model.RaffleFactor{
		UserID:     "zym",
		StrategyID: 100003,
	})
	if err != nil {
		common.Log.Errorf("error: %v", err)
	} else {
		common.Log.Info("抽奖成功")
		common.Log.Infof("抽奖结果 %+v", success)
	}
}

/*
*
从装配的策略中随机获取奖品ID值
*/
func TestGetAssembleRandomVal(t *testing.T) {
	i := 200
	for i > 0 {
		i--
		result := GetRandomAwardIdByWeight("100001", "4000:102,103,104,105")
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
