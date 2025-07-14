package dto

import (
	"github.com/kadekchresna/ecommerce/product-service/internal/v1/model"
)

type GetProductsPaginateRequest struct {
	Limit  int    `query:"limit"`
	Page   int    `query:"page"`
	Search string `query:"search"`
}

type GetProductsPaginateResponse struct {
	Products   []ProductWithStock `json:"products"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	Total      int64              `json:"total"`
	TotalPages int                `json:"total_pages"`
}

type ProductWithStock struct {
	model.Products
	ShopName        string `json:"shop_name"`
	AvailableStock  int    `json:"available_stock"`
	WarehouseStatus string `json:"warehouse_status"`
	WarehouseName   string `json:"warehouse_name"`
}
