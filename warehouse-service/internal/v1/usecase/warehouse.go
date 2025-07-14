package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kadekchresna/ecommerce/warehouse-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/warehouse-service/helper/messaging"
	"github.com/kadekchresna/ecommerce/warehouse-service/infrastructure/lock"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"
	repository_interface "github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/repository/interface"
	usecase_interface "github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/usecase/interface"
)

type warehouseUsecase struct {
	distributedLock     lock.DistributedLock
	WarehouseRepository repository_interface.IWarehouseRepository
}

func NewWarehouseUsecase(
	distributedLock lock.DistributedLock,
	WarehouseRepository repository_interface.IWarehouseRepository,
) usecase_interface.IWarehouseUsecase {
	return &warehouseUsecase{
		distributedLock:     distributedLock,
		WarehouseRepository: WarehouseRepository,
	}
}

func (u *warehouseUsecase) ReserveStock(ctx context.Context, request dto.ReserveStockRequest) error {

	payload := dto.UpdateStockRequest{
		OrderUUID: request.OrderUUID,
		InboxUUID: request.InboxUUID,
	}
	for _, r := range request.ReserveStockDetails {

		if r.StockAmount <= 0 {
			err := fmt.Errorf("error validation :: warehouseUsecase.ReserveStock().StockAmount-Validation <= 0")
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		payload.OrderDetails = append(payload.OrderDetails, dto.OrderDetails{
			ProductUUID:         r.ProductUUID,
			StockAmountToReduce: r.StockAmount,
		})
	}

	return u.WarehouseRepository.UpdateReserveStock(ctx, payload)
}

func (u *warehouseUsecase) ReturnStock(ctx context.Context, request dto.ReserveStockRequest) error {

	payload := dto.UpdateStockRequest{
		OrderUUID: request.OrderUUID,
		InboxUUID: request.InboxUUID,
	}

	for _, r := range request.ReserveStockDetails {

		if r.StockAmount <= 0 {
			err := fmt.Errorf("error validation :: warehouseUsecase.ReturnStock().StockAmount-Validation <= 0")
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		payload.OrderDetails = append(payload.OrderDetails, dto.OrderDetails{
			ProductUUID:         r.ProductUUID,
			StockAmountToReduce: -1 * r.StockAmount,
		})
	}

	return u.WarehouseRepository.UpdateReturnStock(ctx, payload)
}

func (u *warehouseUsecase) GetWarehouse(ctx context.Context, request *dto.GetWarehouseRequest) (*dto.GetWarehouseResponse, error) {

	w, err := u.WarehouseRepository.GetWarehouse(ctx, request.UUID)
	if err != nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	if w == nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. warehouse is not found")
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	return &dto.GetWarehouseResponse{
		Warehouse: w,
	}, nil
}

func (u *warehouseUsecase) GetProductStock(ctx context.Context, request *dto.GetProductStockRequest) (*dto.GetProductStockResponse, error) {

	w, err := u.WarehouseRepository.GetProductStock(ctx, *request)
	if err != nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	if w == nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. product is not found")
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	return &dto.GetProductStockResponse{
		ProductUUID:     w.WarehouseStock.ProductUUID,
		WarehouseUUID:   w.UUID,
		ReserveQuantity: w.WarehouseStock.ReserveQuantity,
		StartQuantity:   w.WarehouseStock.StartQuantity,
		WarehouseName:   w.Name,
		ShopUUID:        w.ShopUUID,
		Status:          w.Status,
	}, nil
}

func (u *warehouseUsecase) TransferProduct(ctx context.Context, request *dto.TransferProductRequest) error {

	w, err := u.WarehouseRepository.GetWarehouse(ctx, request.TargetWarehouseUUID)
	if err != nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	if w == nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. warehouse is not found")
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	return u.WarehouseRepository.TransferProduct(ctx, *request)
}

func (u *warehouseUsecase) UpdateStatusWarehouse(ctx context.Context, request *dto.UpdateStatusWarehouseRequest) error {
	w, err := u.WarehouseRepository.GetWarehouse(ctx, request.WarehouseUUID)
	if err != nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	if w == nil {
		err := fmt.Errorf("error query select :: WarehouseRepository.GetWarehouse. warehouse is not found")
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	return u.WarehouseRepository.UpdateStatusWarehouse(ctx, *request)
}

func (u *warehouseUsecase) StoreToInbox(ctx context.Context, i model.Inbox) error {
	return u.WarehouseRepository.CreateInbox(ctx, i)
}

func (u *warehouseUsecase) ProcessInbox(ctx context.Context, request *dto.ProcessInboxRequest) error {

	inboxes, err := u.WarehouseRepository.GetInboxList(ctx, *request)
	if err != nil {
		err := fmt.Errorf("error Query Select :: WarehouseRepository.GetInboxList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	for _, i := range inboxes {
		u.ProcessReserveStock(ctx, i)
	}

	return nil
}

func (u *warehouseUsecase) ProcessReserveStock(ctx context.Context, i model.Inbox) {

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

	switch payload.Action {
	case helper_messaging.ACTION_ORDER_EXPIRED:

		reserveStockReq := dto.ReserveStockRequest{
			OrderUUID: payload.Order.UUID,
			InboxUUID: i.UUID,
		}
		for _, od := range payload.OrderDetail {
			reserveStockReq.ReserveStockDetails = append(reserveStockReq.ReserveStockDetails, dto.ReserveStockDetailsRequest{
				ProductUUID: od.ProductUUID,
				StockAmount: od.Quantity,
			})
		}

		if err := u.ReturnStock(ctx, reserveStockReq); err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Marshal :: warehouseHandler.ReserveStock(). %s", err.Error()))
			return
		}
	default:
		reserveStockReq := dto.ReserveStockRequest{
			OrderUUID: payload.Order.UUID,
			InboxUUID: i.UUID,
		}
		for _, od := range payload.OrderDetail {
			reserveStockReq.ReserveStockDetails = append(reserveStockReq.ReserveStockDetails, dto.ReserveStockDetailsRequest{
				ProductUUID: od.ProductUUID,
				StockAmount: od.Quantity,
			})
		}

		if err := u.ReserveStock(ctx, reserveStockReq); err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Marshal :: warehouseHandler.ReserveStock(). %s", err.Error()))
			return
		}

	}

}

func (u *warehouseUsecase) ProcessOutbox(ctx context.Context, request *dto.ProcessOutboxRequest) error {

	outboxes, err := u.WarehouseRepository.GetOutboxList(ctx, *request)
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

func (u *warehouseUsecase) ProcessPublishOutbox(ctx context.Context, o model.Outbox) {

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

	if err := u.WarehouseRepository.UpdateOutboxStatus(ctx, &o, model.OutboxStatusSuccess); err != nil {
		err := fmt.Errorf("error process outbox :: OrdersRepository.UpdateOutboxStatus(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return
	}
}
