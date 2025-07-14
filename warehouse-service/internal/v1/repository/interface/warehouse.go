package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
)

type IWarehouseRepository interface {
	GetProductStock(ctx context.Context, request dto.GetProductStockRequest) (*model.Warehouse, error)
	UpdateReserveStock(ctx context.Context, request dto.UpdateStockRequest) error
	UpdateReturnStock(ctx context.Context, request dto.UpdateStockRequest) error
	GetWarehouse(ctx context.Context, warehouseUUID uuid.UUID) (*model.Warehouse, error)
	UpdateStatusWarehouse(ctx context.Context, request dto.UpdateStatusWarehouseRequest) error
	TransferProduct(ctx context.Context, request dto.TransferProductRequest) error

	CreateInbox(ctx context.Context, i model.Inbox) error
	GetInboxList(ctx context.Context, p dto.ProcessInboxRequest) ([]model.Inbox, error)
	GetOutboxList(ctx context.Context, p dto.ProcessOutboxRequest) ([]model.Outbox, error)
	UpdateOutboxStatus(ctx context.Context, o *model.Outbox, newStatus model.OutboxStatusType) error
}
