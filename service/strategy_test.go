package service

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	list, _ := getStrategyAwardList(100002)
	fmt.Println(list)
}
