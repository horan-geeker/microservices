package response

import "microservices/entity/model"

type GoodsListResponse struct {
	List []*model.Goods `json:"list"`
}
