package dto

import "github.com/google/uuid"

type UpdateStockRequest struct {
	OrderDetails []OrderDetails
	OrderUUID    uuid.UUID
	InboxUUID    uuid.UUID
}

type OrderDetails struct {
	ProductUUID         uuid.UUID `json:"product_uuid"`
	StockAmountToReduce int       `json:"stock_amount_to_reduce"`
}

type UpdateStockResponse struct {
	ProductUUID uuid.UUID `json:"product_uuid"`
}

type ReserveStockRequest struct {
	ReserveStockDetails []ReserveStockDetailsRequest
	OrderUUID           uuid.UUID
	InboxUUID           uuid.UUID
}

type ReserveStockDetailsRequest struct {
	ProductUUID uuid.UUID `json:"product_uuid"`
	StockAmount int       `json:"stock_amount"`
}
