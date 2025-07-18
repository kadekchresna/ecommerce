package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	standart_error "github.com/kadekchresna/ecommerce/warehouse-service/helper/error"
	"github.com/kadekchresna/ecommerce/warehouse-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/warehouse-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/warehouse-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/warehouse-service/helper/uuid"
	"github.com/kadekchresna/ecommerce/warehouse-service/infrastructure/messaging"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/helper/dto"
	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/model"

	"github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/warehouse-service/internal/v1/repository/interface"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type warehouseRepository struct {
	db         *gorm.DB
	time       *helper_time.TimeHelper
	helperUUID *helper_uuid.UUIDHelper
	Messager   messaging.Producer
}

func NewWarehouseRepository(
	db *gorm.DB,
	time *helper_time.TimeHelper,
	helperUUID *helper_uuid.UUIDHelper,
	Messager messaging.Producer,
) repository_interface.IWarehouseRepository {
	return &warehouseRepository{
		db:         db,
		time:       time,
		helperUUID: helperUUID,
		Messager:   Messager,
	}
}

func (r *warehouseRepository) GetProductStock(ctx context.Context, request dto.GetProductStockRequest) (*model.Warehouse, error) {

	w := dao.WarehouseStockDAO{}
	if err := r.db.WithContext(ctx).Model(dao.WarehouseStockDAO{}).Where("product_uuid = ?", request.ProductUUID).Preload("Warehouse").First(&w).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	res := model.Warehouse{
		UUID:     w.Warehouse.UUID,
		Name:     w.Warehouse.Name,
		Code:     w.Warehouse.Code,
		Desc:     w.Warehouse.Desc,
		ShopUUID: w.Warehouse.ShopUUID,
		Status:   w.Warehouse.Status,
		WarehouseStock: model.WarehouseStock{
			UUID:            w.UUID,
			ProductUUID:     w.ProductUUID,
			StartQuantity:   w.StartQuantity,
			ReserveQuantity: w.ReserveQuantity,
		},
	}

	return &res, nil
}

func (r *warehouseRepository) UpdateReserveStock(ctx context.Context, request dto.UpdateStockRequest) error {

	var note error
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("PANIC :: warehouseRepository.UpdateStock().Transaction() :: %v", r)
				logger.LogWithContext(ctx).Error(err.Error())
				return
			}

		}()

		var inbox dao.InboxDAO
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uuid = ?", request.InboxUUID).
			Where("status IN ?", []model.InboxStatusType{model.InboxStatusCreated, model.InboxStatusFailed}).
			First(&inbox).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			err = fmt.Errorf("error query create :: warehouseRepository.UpdateStock().First().Inbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		updateStocks := []dao.WarehouseStockDAO{}
		for _, s := range request.OrderDetails {

			ws := dao.WarehouseStockDAO{}
			if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_uuid = ?", s.ProductUUID).First(&ws).Error; err != nil {
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().SELECT-FOR-UPDATE. %s", err.Error()))
				return err
			}

			w := dao.WarehouseDAO{}
			if err := tx.Model(dao.WarehouseDAO{}).Where("uuid = ?", ws.WarehouseUUID).First(&w).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = fmt.Errorf("warehouse not found for product %s", s.ProductUUID.String())
					logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().Get-Warehouse. %s", err.Error()))
					return err
				}

				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().Get-Warehouse. %s", err.Error()))
				return err
			}

			available := ws.StartQuantity - ws.ReserveQuantity

			if w.Status != model.WarehouseStatusActive {
				note = standart_error.ErrorWarehouseInactive
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().Get-Warehouse. %s", note.Error()))
				break
			}

			if available < s.StockAmountToReduce {
				note = standart_error.ErrorInsufficientStock
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().Stock-Calculation. %s", note.Error()))
				break
			}

			if available >= s.StockAmountToReduce {

				ws.ReserveQuantity += s.StockAmountToReduce
				updateStocks = append(updateStocks, ws)
			}
		}

		if len(updateStocks) == len(request.OrderDetails) {
			if err := tx.Save(&updateStocks).Error; err != nil {
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Update :: warehouseRepository.UpdateStock().Save(). %s", err.Error()))
				return err
			}
		}

		messageID := r.helperUUID.New()
		eventPayload := model.OutboxOrderUpdateStatusMetaRequest{
			Order: model.Order{
				UUID: request.OrderUUID,
			},
			MessageID: messageID.String(),
			Action:    helper_messaging.ACTION_RESERVE_SUCCESS,
		}

		if note != nil && errors.Is(note, standart_error.ErrorInsufficientStock) {
			eventPayload.Action = helper_messaging.ACTION_RESERVE_FAILED_INSUFFICIENT_STOCK
		}

		if note != nil && errors.Is(note, standart_error.ErrorWarehouseInactive) {
			eventPayload.Action = helper_messaging.ACTION_RESERVE_FAILED_WAREHOUSE_INACTIVE
		}

		payload, err := json.Marshal(eventPayload)
		if err != nil {
			err = fmt.Errorf("error Marshal :: warehouseRepository.UpdateStock().Update().Order. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		outbox := dao.OutboxDAO{
			UUID:      messageID,
			Metadata:  string(payload),
			Type:      helper_messaging.TOPIC_WAREHOUSE_EVENTS,
			Status:    model.OutboxStatusType(model.OutboxStatusCreated),
			Reference: request.OrderUUID.String(),
			Action:    eventPayload.Action,
			Response:  "{}",
		}

		if err := tx.Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "type"}, {Name: "reference"}, {Name: "action"}},
				DoNothing: true,
			},
		).Create(&outbox).Error; err != nil {
			err = fmt.Errorf("error Query Create :: warehouseRepository.UpdateStock().Create().Outbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		inbox.Status = model.InboxStatusSuccess
		inbox.UpdatedAt = r.time.Now()
		inbox.RetryCount += 1
		resPayload := model.OutboxGeneralMetaResponse{}

		if err != nil {
			inbox.Status = model.InboxStatusFailed
			resPayload.ErrorReason = err.Error()
		}

		if note != nil {
			resPayload.ErrorReason = note.Error()
		}

		resPayloadString, err := json.Marshal(resPayload)
		if err != nil {
			err = fmt.Errorf("error Marshal :: warehouseRepository.UpdateStock().Update().Inbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())

			return
		}

		inbox.Response = string(resPayloadString)

		if err = tx.Save(&inbox).Error; err != nil {
			err = fmt.Errorf("error query update :: warehouseRepository.UpdateStock().Save().Inbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())

			return
		}

		return nil
	})
}

func (r *warehouseRepository) UpdateReturnStock(ctx context.Context, request dto.UpdateStockRequest) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		var inbox dao.InboxDAO
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("PANIC :: warehouseRepository.UpdateStock().Transaction() :: %v", r)
				logger.LogWithContext(ctx).Error(err.Error())
				return
			}

			inbox.Status = model.InboxStatusSuccess
			inbox.UpdatedAt = r.time.Now()
			inbox.RetryCount += 1
			resPayload := model.OutboxGeneralMetaResponse{}

			if err != nil {
				inbox.Status = model.InboxStatusFailed
				resPayload.ErrorReason = err.Error()
			}

			resPayloadString, err := json.Marshal(resPayload)
			if err != nil {
				err = fmt.Errorf("error Marshal :: warehouseRepository.UpdateStock().Update().Inbox. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())

				return
			}

			inbox.Response = string(resPayloadString)

			if err = tx.Save(&inbox).Error; err != nil {
				err = fmt.Errorf("error query update :: warehouseRepository.UpdateStock().Save().Inbox. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())

				return
			}

		}()

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uuid = ?", request.InboxUUID).
			Where("status IN ?", []model.InboxStatusType{model.InboxStatusCreated, model.InboxStatusFailed}).
			First(&inbox).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			err = fmt.Errorf("error query create :: warehouseRepository.UpdateStock().First().Inbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		for _, s := range request.OrderDetails {

			ws := dao.WarehouseStockDAO{}
			if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_uuid = ?", s.ProductUUID).First(&ws).Error; err != nil {
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().SELECT-FOR-UPDATE. %s", err.Error()))
				return err
			}

			w := dao.WarehouseDAO{}
			if err := tx.Model(dao.WarehouseDAO{}).Where("uuid = ?", ws.WarehouseUUID).First(&w).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = fmt.Errorf("warehouse not found for product %s", s.ProductUUID.String())
					logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().Get-Warehouse. %s", err.Error()))
					return err
				}

				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.UpdateStock().Get-Warehouse. %s", err.Error()))
				return err
			}

			ws.ReserveQuantity += s.StockAmountToReduce

			if err := tx.Save(&ws).Error; err != nil {
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Update :: warehouseRepository.UpdateStock().Save(). %s", err.Error()))
				return err
			}

		}

		return nil
	})
}

func (r *warehouseRepository) TransferProduct(ctx context.Context, request dto.TransferProductRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("PANIC :: warehouseRepository.TransferProduct().Transaction() :: %v", r)
				logger.LogWithContext(ctx).Error(err.Error())
			}
		}()

		w := dao.WarehouseStockDAO{}
		if err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_uuid = ?", request.ProductUUID).First(&w).Error; err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Select :: warehouseRepository.TransferProduct().SELECT-FOR-UPDATE. %s", err.Error()))
			return err
		}

		w.WarehouseUUID = request.TargetWarehouseUUID

		if err := tx.Save(&w).Error; err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Update :: warehouseRepository.TransferProduct().Save(). %s", err.Error()))
			return err
		}

		return nil
	}, &sql.TxOptions{})

}

func (r *warehouseRepository) GetWarehouse(ctx context.Context, warehouseUUID uuid.UUID) (*model.Warehouse, error) {

	w := dao.WarehouseDAO{}
	if err := r.db.WithContext(ctx).Model(dao.WarehouseDAO{}).First(&w).Where("uuid = ?", warehouseUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	res := model.Warehouse{
		UUID:     w.UUID,
		Name:     w.Name,
		Code:     w.Code,
		Desc:     w.Desc,
		ShopUUID: w.ShopUUID,
		Status:   w.Status,
	}

	return &res, nil
}

func (r *warehouseRepository) UpdateStatusWarehouse(ctx context.Context, request dto.UpdateStatusWarehouseRequest) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("PANIC :: warehouseRepository.TransferProduct().Transaction() :: %v", r)
				logger.LogWithContext(ctx).Error(err.Error())
			}
		}()

		if err := tx.Model(dao.WarehouseDAO{}).
			Where("uuid = ?", request.WarehouseUUID).
			Updates(map[string]interface{}{
				"status":     request.Status,
				"updated_by": request.UserUUID,
				"updated_at": r.time.Now(),
			}).Error; err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Update :: warehouseRepository.UpdateStatusWarehouse(). %s", err.Error()))
			return err
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

}

func (r *warehouseRepository) CreateInbox(ctx context.Context, i model.Inbox) error {
	inbox := dao.InboxDAO{
		UUID:       r.helperUUID.New(),
		Metadata:   i.Metadata,
		Response:   i.Response,
		Status:     i.Status,
		Type:       i.Type,
		Reference:  i.Reference,
		Action:     i.Action,
		RetryCount: 0,
	}

	err := r.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "type"}, {Name: "reference"}, {Name: "action"}},
			DoNothing: true,
		},
	).Create(&inbox).Error

	if err != nil {
		err = fmt.Errorf("error Query Create :: warehouseRepository.CreateInbox(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	return nil
}

func (r *warehouseRepository) GetInboxList(ctx context.Context, p dto.ProcessInboxRequest) ([]model.Inbox, error) {
	var ob []dao.InboxDAO

	tx := r.db.WithContext(ctx).Model(&dao.InboxDAO{})

	if p.Status != "" {
		tx = tx.Where("status = ?", p.Status)
	}

	if p.Type != "" {
		tx = tx.Where("type = ?", p.Type)
	}

	if !p.OlderThan.IsZero() {
		tx = tx.Where("created_at <= ?", p.OlderThan)
	}

	if len(p.Statuses) > 0 {
		tx = tx.Where("status IN ?", p.Statuses)
	}

	if p.Limit <= 0 {
		p.Limit = 100
	}

	err := tx.Order("created_at ASC").
		Limit(p.Limit).
		Offset(p.Offset).
		Find(&ob).Error

	if err != nil {
		err = fmt.Errorf("error Query Select :: warehouseRepository.GetInboxList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	result := make([]model.Inbox, 0, len(ob))

	for _, o := range ob {
		result = append(result, model.Inbox{
			UUID:       o.UUID,
			Metadata:   o.Metadata,
			Response:   o.Response,
			Status:     o.Status,
			Type:       o.Type,
			RetryCount: o.RetryCount,
			Reference:  o.Reference,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
			Action:     o.Action,
		})
	}

	return result, err
}

func (r *warehouseRepository) GetOutboxList(ctx context.Context, p dto.ProcessOutboxRequest) ([]model.Outbox, error) {
	var ob []dao.OutboxDAO

	tx := r.db.WithContext(ctx).Model(&dao.OutboxDAO{})

	if p.Status != "" {
		tx = tx.Where("status = ?", p.Status)
	}

	if p.Type != "" {
		tx = tx.Where("type = ?", p.Type)
	}

	if !p.OlderThan.IsZero() {
		tx = tx.Where("created_at <= ?", p.OlderThan)
	}

	if len(p.Statuses) > 0 {
		tx = tx.Where("status IN ?", p.Statuses)
	}

	if p.Limit <= 0 {
		p.Limit = 100
	}

	err := tx.Order("created_at ASC").
		Limit(p.Limit).
		Offset(p.Offset).
		Find(&ob).Error

	if err != nil {
		err = fmt.Errorf("error Query Select :: ordersRepository.GetOutboxList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	result := make([]model.Outbox, 0, len(ob))

	for _, o := range ob {
		result = append(result, model.Outbox{
			UUID:       o.UUID,
			Metadata:   o.Metadata,
			Response:   o.Response,
			Status:     o.Status,
			Type:       o.Type,
			RetryCount: o.RetryCount,
			Reference:  o.Reference,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
			Action:     o.Action,
		})
	}

	return result, err
}

func (r *warehouseRepository) UpdateOutboxStatus(ctx context.Context, o *model.Outbox, newStatus model.OutboxStatusType) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var outbox dao.OutboxDAO

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uuid = ?", o.UUID).
			Where("status IN ?", []model.OutboxStatusType{model.OutboxStatusCreated, model.OutboxStatusFailed}).
			First(&outbox).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			err := fmt.Errorf("error query create :: warehouseRepository.First().Outbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		outbox.Status = newStatus
		outbox.UpdatedAt = r.time.Now()
		outbox.RetryCount += 1

		resPayload := model.OutboxGeneralMetaResponse{
			MessageID: o.UUID.String(),
		}
		errPublish := r.Messager.Publish(ctx, []byte(o.Type), []byte(o.Metadata))
		if errPublish != nil {
			errPublish = fmt.Errorf("error publish message :: Messager.Publish(). %s", errPublish.Error())
			logger.LogWithContext(ctx).Error(errPublish.Error())

			resPayload.ErrorReason = errPublish.Error()
			outbox.Status = model.OutboxStatusFailed
		}

		resPayloadString, err := json.Marshal(resPayload)
		if err != nil {
			err = fmt.Errorf("error Marshal :: warehouseRepository.UpdateOrderStatusWithOutbox().Update().Order. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		outbox.Response = string(resPayloadString)

		if err := tx.Save(&outbox).Error; err != nil {
			err := fmt.Errorf("error query create :: warehouseRepository.Save().Outbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		return nil
	})
}
