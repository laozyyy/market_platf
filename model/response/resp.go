package response

type StrategyAwardListResponse struct {
	Data    StrategyAwardPageInfo `json:"data"`
	Message string                `json:"message"`
}

type StrategyAwardPageInfo struct {
	PageNum       int64 `json:"page_num"`
	PageSize      int64 `json:"page_size"`
	IsEnd         bool  `json:"is_end"`
	Count         int32 `json:"count"`
	StrategyAward []StrategyAwardDTO
}

type StrategyAwardDTO struct {
	AwardId       int    `json:"awardId"`       // 抽奖奖品ID - 内部流转使用
	AwardTitle    string `json:"awardTitle"`    // 抽奖奖品标题
	AwardSubtitle string `json:"awardSubtitle"` // 抽奖奖品副标题
	Sort          int    `json:"sort"`          // 排序
}
