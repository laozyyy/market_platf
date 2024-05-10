package request

type CreateOrderRequest struct {
	UserID        string `json:"user_id"`
	Sku           int64  `json:"sku"`
	OutBusinessNo string `json:"out_business_no"`
}
