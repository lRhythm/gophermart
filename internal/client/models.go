package client

import (
	"github.com/lRhythm/gophermart/internal/models"
	"net/http"
)

type Client struct {
	client *http.Client
}

// In comments for local testing.

type orderSendRequest struct {
	ID string `json:"order"`
	//Goods []orderSendGoodsRequest `json:"goods"`
}

//type orderSendGoodsRequest struct {
//	Description string `json:"description"`
//	Price       uint   `json:"price"`
//}

func newOrderSendRequest(orderID string) *orderSendRequest {
	return &orderSendRequest{
		ID: orderID,
		//Goods: []orderSendGoodsRequest{
		//	{
		//		Description: "apple 1",
		//		Price:       1,
		//	},
		//	{
		//		Description: "apple 2",
		//		Price:       10,
		//	},
		//	{
		//		Description: "apple 3",
		//		Price:       100,
		//	},
		//	{
		//		Description: "apple 4",
		//		Price:       1000,
		//	},
		//	{
		//		Description: "apple 5",
		//		Price:       10000,
		//	},
		//},
	}
}

type orderGetResponse struct {
	ID      string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func (o *orderGetResponse) toDTO() models.AccrualOrderDTO {
	return models.AccrualOrderDTO{
		ID:      o.ID,
		Status:  o.Status,
		Accrual: o.Accrual,
	}
}
