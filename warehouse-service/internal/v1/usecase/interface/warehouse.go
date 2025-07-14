package usecase_interface

import (
	"context"

	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type IWarehouseUsecase interface {
	ReserveStock(ctx context.Context, request dto.ReserveStockRequest) error
	ReturnStock(ctx context.Context, request dto.ReserveStockRequest) error
	GetWarehouse(ctx context.Context, request *dto.GetWarehouseRequest) (*dto.GetWarehouseResponse, error)
	TransferProduct(ctx context.Context, request *dto.TransferProductRequest) error
	UpdateStatusWarehouse(ctx context.Context, request *dto.UpdateStatusWarehouseRequest) error
	GetProductStock(ctx context.Context, request *dto.GetProductStockRequest) (*dto.GetProductStockResponse, error)

	StoreToInbox(ctx context.Context, i model.Inbox) error
	ProcessInbox(ctx context.Context, request *dto.ProcessInboxRequest) error
	ProcessOutbox(ctx context.Context, request *dto.ProcessOutboxRequest) error
}
