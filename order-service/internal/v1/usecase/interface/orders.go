package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
)

type IOrdersUsecase interface {
	Checkout(ctx context.Context, request *dto.CheckoutRequest) error
	UpdateOrderStatusExpired(ctx context.Context) error
	UpdateOrderStatusCompleted(ctx context.Context, orderUUID uuid.UUID) error

	ProcessOutbox(ctx context.Context, request *dto.ProcessOutboxRequest) error
	ProcessInbox(ctx context.Context, request *dto.ProcessInboxRequest) error
	StoreToInbox(ctx context.Context, i model.Inbox) error
}
