package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kadekchresna/ecommerce/order-service/helper/logger"
	helper_messaging "github.com/kadekchresna/ecommerce/order-service/helper/messaging"
	helper_time "github.com/kadekchresna/ecommerce/order-service/helper/time"
	helper_uuid "github.com/kadekchresna/ecommerce/order-service/helper/uuid"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/messaging"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/dto"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/model"
	"github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/dao"
	repository_interface "github.com/kadekchresna/ecommerce/order-service/internal/v1/repository/interface"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ordersRepository struct {
	db         *gorm.DB
	Messager   messaging.Producer
	timeHelper helper_time.TimeHelper
	uuidHelper helper_uuid.UUIDHelper
}

func NewOrdersRepository(
	db *gorm.DB,
	timeHelper helper_time.TimeHelper,
	uuidHelper helper_uuid.UUIDHelper,
	Messager messaging.Producer,
) repository_interface.IOrdersRepository {
	return &ordersRepository{
		db:         db,
		timeHelper: timeHelper,
		uuidHelper: uuidHelper,
		Messager:   Messager,
	}
}

func (r *ordersRepository) Checkout(ctx context.Context, request *dto.CreateCheckoutRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("PANIC :: ordersRepository.Checkout().defer() :: %v", r)
				logger.LogWithContext(ctx).Error(err.Error())
			}
		}()

		order := dao.OrderDAO{
			UUID:        request.Order.UUID,
			Code:        request.Order.Code,
			UserUUID:    request.Order.UserUUID,
			TotalAmount: request.Order.TotalAmount,
			ExpiredAt:   request.Order.ExpiredAt,
			Status:      request.Order.Status,
			CreatedBy:   request.Order.CreatedBy,
			UpdatedBy:   request.Order.UpdatedBy,
			Metadata:    request.Order.Metadata,
		}

		if err := tx.Create(&order).Error; err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Create :: ordersRepository.Checkout().Create().Order. %s", err.Error()))
			return err
		}

		orderDetails := make([]dao.OrderDetailDAO, 0, len(request.OrderDetails))

		for i := range request.OrderDetails {

			orderDetails = append(orderDetails, dao.OrderDetailDAO{

				ProductUUID:  request.OrderDetails[i].ProductUUID,
				ProductTitle: request.OrderDetails[i].ProductTitle,
				ProductPrice: request.OrderDetails[i].ProductPrice,
				Quantity:     request.OrderDetails[i].Quantity,
				SubTotal:     request.OrderDetails[i].SubTotal,
				OrderUUID:    request.Order.UUID,
			})

		}

		if err := tx.Create(&orderDetails).Error; err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Create :: ordersRepository.Checkout().Create().OrderDetail. %s", err.Error()))
			return err
		}

		messageID := r.uuidHelper.New()
		request.EventPayload.MessageID = messageID.String()

		payload, err := json.Marshal(request.EventPayload)
		if err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Marshall payload :: ordersRepository.Checkout().Marshal. %s", err.Error()))
			return err
		}

		outbox := dao.OutboxDAO{
			UUID:      messageID,
			Metadata:  string(payload),
			Type:      request.EventType,
			Status:    model.OutboxStatusType(model.OrderStatusCreated),
			Reference: request.Order.UUID.String(),
			Action:    helper_messaging.ACTION_ORDER_CREATED,
			Response:  "{}",
		}

		if err := tx.Create(&outbox).Error; err != nil {
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error Query Create :: ordersRepository.Checkout().Create().Outbox. %s", err.Error()))
			return err
		}

		return nil
	})
}

func (r *ordersRepository) UpdateOrderStatusWithOutbox(ctx context.Context, p dto.UpdateOrderStatusRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("PANIC :: ordersRepository.UpdateOrderStatusWithOutbox().defer() :: %v", r)
				logger.LogWithContext(ctx).Error(err.Error())
			}
		}()

		var o dao.OrderDAO

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uuid = ?", p.OrderUUID).
			First(&o).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.New("error Query Select :: ordersRepository.().UpdateOrderStatusWithOutbox().SELECT-FOR-UPDATE.Order. Order not found")
				logger.LogWithContext(ctx).Error(err.Error())
				return err
			}
			return err
		}

		od := []dao.OrderDetailDAO{}
		if err := tx.
			Where("order_uuid = ?", p.OrderUUID).
			Find(&od).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.New("error Query Select :: ordersRepository.().UpdateOrderStatusWithOutbox().OrderDetail. order detail not found")
				logger.LogWithContext(ctx).Error(err.Error())
				return err
			}
			return err
		}

		metadataString, err := json.Marshal(p.Metadata)
		if err != nil {
			err = fmt.Errorf("error marshal :: ordersRepository.().UpdateOrderStatusWithOutbox().Marshal.Metadata. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		if err := tx.Model(&o).Updates(map[string]interface{}{
			"status":     p.NewStatus,
			"metadata":   string(metadataString),
			"updated_at": r.timeHelper.Now(),
		}).Error; err != nil {
			err = fmt.Errorf("error Query Update :: ordersRepository.UpdateOrderStatusWithOutbox().Update().Order. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		if p.NewStatus == model.OrderStatusExpired {

			order := model.Order{
				UUID: o.UUID,
			}

			orderDetail := []model.OrderDetail{}
			for _, d := range od {
				orderDetail = append(orderDetail, model.OrderDetail{
					ProductUUID: d.ProductUUID,
					Quantity:    d.Quantity,
				})
			}

			messageID := r.uuidHelper.New()
			eventPayload := model.OutboxOrderUpdateStatusMetaRequest{
				Order:       order,
				OrderDetail: orderDetail,
				MessageID:   messageID.String(),
				Action:      helper_messaging.ACTION_ORDER_EXPIRED,
			}

			payload, err := json.Marshal(eventPayload)
			if err != nil {
				err = fmt.Errorf("error Marshal :: ordersRepository.UpdateOrderStatusWithOutbox().Update().Order. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())
				return err
			}

			outbox := dao.OutboxDAO{
				UUID:      messageID,
				Metadata:  string(payload),
				Type:      p.EventType,
				Status:    model.OutboxStatusType(model.OrderStatusCreated),
				Reference: order.UUID.String(),
				Action:    helper_messaging.ACTION_ORDER_EXPIRED,
				Response:  "{}",
			}

			if err := tx.Clauses(
				clause.OnConflict{
					Columns: []clause.Column{{Name: "type"}, {Name: "reference"}, {Name: "action"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"metadata": string(payload),
					}),
				},
			).Create(&outbox).Error; err != nil {
				err = fmt.Errorf("error Query Create :: ordersRepository.UpdateOrderStatusWithOutbox().Create().Outbox. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())
				return err
			}

		}

		if p.NewStatus == model.OrderStatusInPayment ||
			p.NewStatus == model.OrderStatusFailed {

			var inbox dao.InboxDAO
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("uuid = ?", p.InboxUUID).
				Where("status IN ?", []model.InboxStatusType{model.InboxStatusCreated, model.InboxStatusFailed}).
				First(&inbox).Error; err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil
				}

				err = fmt.Errorf("error query create :: warehouseRepository.UpdateStock().First().Inbox. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())
				return err
			}

			inbox.Status = model.InboxStatusSuccess
			inbox.UpdatedAt = r.timeHelper.Now()
			inbox.RetryCount += 1
			resPayload := model.OutboxGeneralMetaResponse{}

			resPayloadString, err := json.Marshal(resPayload)
			if err != nil {
				err = fmt.Errorf("error Marshal :: warehouseRepository.UpdateStock().Update().Inbox. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())

				return err
			}

			inbox.Response = string(resPayloadString)

			if err = tx.Save(&inbox).Error; err != nil {
				err = fmt.Errorf("error query update :: warehouseRepository.UpdateStock().Save().Inbox. %s", err.Error())
				logger.LogWithContext(ctx).Error(err.Error())

				return err
			}

		}
		return nil
	})
}

func (r *ordersRepository) GetOrderList(ctx context.Context, p dto.ProcessExpiredOrderRequest) ([]model.Order, error) {
	var result []dao.OrderDAO

	tx := r.db.WithContext(ctx).Model(&dao.OrderDAO{})

	if p.Status != "" {
		tx = tx.Where("status = ?", p.Status)
	}

	if !p.OlderThan.IsZero() {
		tx = tx.Where("expired_at <= ?", p.OlderThan)
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
		Find(&result).Error

	if err != nil {
		err = fmt.Errorf("error Query Select :: ordersRepository.GetOrderList(). %s", err.Error())
		logger.LogWithContext(ctx).Error(err.Error())
		return nil, err
	}

	orders := make([]model.Order, 0, len(result))
	for _, o := range result {
		orders = append(orders, model.Order{
			UUID:        o.UUID,
			Code:        o.Code,
			UserUUID:    o.UserUUID,
			TotalAmount: o.TotalAmount,
			ExpiredAt:   o.ExpiredAt,
			Status:      o.Status,
		})
	}

	return orders, err
}

func (r *ordersRepository) GetOutboxList(ctx context.Context, p dto.ProcessOutboxRequest) ([]model.Outbox, error) {
	var ob []dao.OutboxDAO

	tx := r.db.WithContext(ctx).Model(&dao.OutboxDAO{})

	if p.Status != "" {
		tx = tx.Where("status = ?", p.Status)
	}

	if len(p.Statuses) > 0 {
		tx = tx.Where("status IN ?", p.Statuses)
	}

	if p.Type != "" {
		tx = tx.Where("type = ?", p.Type)
	}

	if !p.OlderThan.IsZero() {
		tx = tx.Where("created_at <= ?", p.OlderThan)
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

func (r *ordersRepository) UpdateOutboxStatus(ctx context.Context, o *model.Outbox, newStatus model.OutboxStatusType) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var outbox dao.OutboxDAO

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uuid = ?", o.UUID).
			Where("status IN ?", []model.OutboxStatusType{model.OutboxStatusCreated, model.OutboxStatusFailed}).
			First(&outbox).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			err := fmt.Errorf("error query create :: ordersRepository.First().Outbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		outbox.Status = newStatus
		outbox.UpdatedAt = r.timeHelper.Now()
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
			err = fmt.Errorf("error Marshal :: ordersRepository.UpdateOrderStatusWithOutbox().Update().Order. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		outbox.Response = string(resPayloadString)

		if err := tx.Save(&outbox).Error; err != nil {
			err := fmt.Errorf("error query create :: ordersRepository.Save().Outbox. %s", err.Error())
			logger.LogWithContext(ctx).Error(err.Error())
			return err
		}

		return nil
	})
}

func (r *ordersRepository) CreateInbox(ctx context.Context, i model.Inbox) error {
	inbox := dao.InboxDAO{
		UUID:       r.uuidHelper.New(),
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

func (r *ordersRepository) GetInboxList(ctx context.Context, p dto.ProcessInboxRequest) ([]model.Inbox, error) {
	var ob []dao.InboxDAO

	tx := r.db.WithContext(ctx).Model(&dao.InboxDAO{})

	if p.Status != "" {
		tx = tx.Where("status = ?", p.Status)
	}

	if len(p.Statuses) > 0 {
		tx = tx.Where("status IN ?", p.Statuses)
	}

	if p.Type != "" {
		tx = tx.Where("type = ?", p.Type)
	}

	if !p.OlderThan.IsZero() {
		tx = tx.Where("created_at <= ?", p.OlderThan)
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
