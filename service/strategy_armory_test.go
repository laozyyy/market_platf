package service

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestStrategyArmory(t *testing.T) {
	success := AssembleLotteryStrategy(100002)
	log.Printf("测试结果：%v\n", success)
}

/*
*

从装配的策略中随机获取奖品ID值
*/
func TestGetAssembleRandomVal(t *testing.T) {
	i := 200
	for i > 0 {
		i--
		result := GetRandomAwardId(100002)
		log.Printf("测试结果：%v - 奖品ID值\n", result)
	}

}

func testMap() {
	strategyMap := make(map[int]int)
	// 添加内容到Map中
	strategyMap[1] = 101
	strategyMap[2] = 101
	strategyMap[3] = 101
	strategyMap[4] = 102
	strategyMap[5] = 102
	strategyMap[6] = 102
	strategyMap[7] = 103
	strategyMap[8] = 103
	strategyMap[9] = 104
	strategyMap[10] = 105

	log.Printf("测试结果：%v\n", strategyMap[1])
}

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
