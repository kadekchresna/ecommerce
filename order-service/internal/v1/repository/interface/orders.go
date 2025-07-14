package repository_interface

import (
	"context"

	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type IOrdersRepository interface {
	Checkout(ctx context.Context, request *dto.CreateCheckoutRequest) error
	UpdateOrderStatusWithOutbox(ctx context.Context, p dto.UpdateOrderStatusRequest) error
	GetOrderList(ctx context.Context, p dto.ProcessExpiredOrderRequest) ([]model.Order, error)

	GetOutboxList(ctx context.Context, p dto.ProcessOutboxRequest) ([]model.Outbox, error)
	UpdateOutboxStatus(ctx context.Context, o *model.Outbox, newStatus model.OutboxStatusType) error
	CreateInbox(ctx context.Context, i model.Inbox) error
	GetInboxList(ctx context.Context, p dto.ProcessInboxRequest) ([]model.Inbox, error)
}
