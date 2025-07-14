package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	standart_error "github.com/kadekchresna/ecommerce/order-service/helper/error"
	"github.com/kadekchresna/ecommerce/order-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/order-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/order-service/helper/uuid"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/lock"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	repository_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/interface"
	usecase_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/usecase/interface"
)

type ordersUsecase struct {
	timeHelper        helper_time.TimeHelper
	uuidHelper        helper_uuid.UUIDHelper
	distributedLock   lock.DistributedLock
	OrdersRepository  repository_interface.IOrdersRepository
	ProductRepository repository_interface.IProductRepository
}

func NewOrdersUsecase(
	timeHelper helper_time.TimeHelper,
	uuidHelper helper_uuid.UUIDHelper,
	distributedLock lock.DistributedLock,
	OrdersRepository repository_interface.IOrdersRepository,
	ProductRepository repository_interface.IProductRepository,
) usecase_interface.IOrdersUsecase {
	return &ordersUsecase{
		uuidHelper:        uuidHelper,
		timeHelper:        timeHelper,
		distributedLock:   distributedLock,
		OrdersRepository:  OrdersRepository,
		ProductRepository: ProductRepository,
	}
}

func (u *ordersUsecase) Checkout(ctx context.Context, request *dto.CheckoutRequest) error {
	orderUUID := u.uuidHelper.New()
	order := model.Order{
		UUID:        orderUUID,
		Code:        fmt.Sprintf("ORD-%s", orderUUID.String()[:8]),
		UserUUID:    request.UserUUID,
		TotalAmount: 0,
		ExpiredAt:   u.timeHelper.Now().Add(15 * time.Minute),
		Status:      model.OrderStatusCreated,
		CreatedBy:   request.UserUUID,
		UpdatedBy:   request.UserUUID,
		Metadata:    "{}",
	}

	orderDetails := []model.OrderDetail{}

	for _, d := range request.OrderDetails {
		product, err := u.ProductRepository.GetProduct(ctx, d.ProductUUID)
		if err != nil {
			err := fmt.Errorf("error Query Select :: ProductRepository.GetProduct(). %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		if product == nil {
			err := errors.New("error Query Select :: ProductRepository.GetProduct(). product is not found")
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		subTotal := product.Price * float64(d.Quantity)
		order.TotalAmount += subTotal

		orderDetails = append(orderDetails, model.OrderDetail{
			ProductUUID:  product.UUID,
			ProductTitle: product.Title,
			ProductPrice: product.Price,
			Quantity:     d.Quantity,
			SubTotal:     subTotal,
		})
	}

	eventPayload := model.OutboxOrderCreatedMetaRequest{
		Order:       order,
		OrderDetail: orderDetails,
		Action:      helper_messaging.ACTION_ORDER_CREATED,
	}

	if err := u.OrdersRepository.Checkout(ctx, &dto.CreateCheckoutRequest{
		Order:        order,
		OrderDetails: orderDetails,
		EventType:    helper_messaging.TOPIC_ORDER_EVENTS,
		UserUUID:     request.UserUUID,
		EventPayload: eventPayload,
	}); err != nil {
		err = fmt.Errorf("error Query create :: OrdersRepository.Checkout(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	return nil
}

func (u *ordersUsecase) UpdateOrderStatusExpired(ctx context.Context) error {

	orders, err := u.OrdersRepository.GetOrderList(ctx, dto.ProcessExpiredOrderRequest{
		Statuses:  []model.OrderStatusType{model.OrderStatusCreated, model.OrderStatusInPayment},
		OlderThan: u.timeHelper.Now(),
		Limit:     10,
	})
	if err != nil {
		err := fmt.Errorf("error Query Select :: OrdersRepository.GetOrderList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	for _, o := range orders {
		u.OrdersRepository.UpdateOrderStatusWithOutbox(ctx, dto.UpdateOrderStatusRequest{
			OrderUUID: o.UUID,
			NewStatus: model.OrderStatusExpired,
			EventType: helper_messaging.TOPIC_ORDER_EVENTS,
		})
	}

	return nil
}

func (u *ordersUsecase) ProcessOutbox(ctx context.Context, request *dto.ProcessOutboxRequest) error {

	outboxes, err := u.OrdersRepository.GetOutboxList(ctx, *request)
	if err != nil {
		err := fmt.Errorf("error Query Select :: OrdersRepository.GetOutboxList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	for _, o := range outboxes {

		u.ProcessPublishOutbox(ctx, o)
	}

	return nil
}

func (u *ordersUsecase) ProcessPublishOutbox(ctx context.Context, o model.Outbox) {

	key := fmt.Sprintf("outbox-%s", o.UUID.String())
	locked, err := u.distributedLock.Acquire(ctx, key, 10)
	if err != nil {
		err := fmt.Errorf("error distributed lock :: distributedLock.Acquire(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return
	}

	if !locked {
		return
	}

	defer u.distributedLock.Release(ctx, key)

	if err := u.OrdersRepository.UpdateOutboxStatus(ctx, &o, model.OutboxStatusSuccess); err != nil {
		err := fmt.Errorf("error process outbox :: OrdersRepository.UpdateOutboxStatus(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return
	}
}

func (u *ordersUsecase) StoreToInbox(ctx context.Context, i model.Inbox) error {
	return u.OrdersRepository.CreateInbox(ctx, i)
}

func (u *ordersUsecase) ProcessInbox(ctx context.Context, request *dto.ProcessInboxRequest) error {

	inboxes, err := u.OrdersRepository.GetInboxList(ctx, *request)
	if err != nil {
		err := fmt.Errorf("error Query Select :: WarehouseRepository.GetInboxList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	for _, i := range inboxes {
		u.ProcessUpdateOrderStatus(ctx, i)
	}

	return nil
}

func (u *ordersUsecase) ProcessUpdateOrderStatus(ctx context.Context, i model.Inbox) {

	key := fmt.Sprintf("inbox-%s", i.UUID.String())

	locked, err := u.distributedLock.Acquire(ctx, key, 10)
	if err != nil {
		err := fmt.Errorf("error distributed lock :: distributedLock.Acquire(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return
	}

	if !locked {
		return
	}

	defer u.distributedLock.Release(ctx, key)

	payload := model.OutboxOrderCreatedMetaRequest{}

	if err := json.Unmarshal([]byte(i.Metadata), &payload); err != nil {

		logger.LogWithContext(ctx).Error(fmt.Sprintf("error Marshal :: warehouseHandler.OrderStockEvent(). %s", err.Error()))
		return
	}

	status := model.OrderStatusInPayment

	metadataOrder := model.OrderMetadata{}

	if payload.Action == helper_messaging.ACTION_RESERVE_FAILED_INSUFFICIENT_STOCK {
		status = model.OrderStatusFailed
		metadataOrder.Reason = standart_error.ErrorInsufficientStock.Error()
	}

	if payload.Action == helper_messaging.ACTION_RESERVE_FAILED_WAREHOUSE_INACTIVE {
		status = model.OrderStatusFailed
		metadataOrder.Reason = standart_error.ErrorWarehouseInactive.Error()
	}

	u.OrdersRepository.UpdateOrderStatusWithOutbox(ctx, dto.UpdateOrderStatusRequest{
		OrderUUID: payload.Order.UUID,
		NewStatus: status,
		Metadata:  metadataOrder,
		InboxUUID: i.UUID,
	})

}

func (u *ordersUsecase) UpdateOrderStatusCompleted(ctx context.Context, orderUUID uuid.UUID) error {

	return u.OrdersRepository.UpdateOrderStatusWithOutbox(ctx, dto.UpdateOrderStatusRequest{
		OrderUUID: orderUUID,
		NewStatus: model.OrderStatusCompleted,
	})
}
